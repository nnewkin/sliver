package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type SliverClient struct {
	conn        *grpc.ClientConn
	address     string
	certsDir    string
	mu          sync.RWMutex
	connected   bool
	reconnectCh chan struct{}
}

func NewSliverClient(address, certsDir string) (*SliverClient, error) {
	client := &SliverClient{
		address:     address,
		certsDir:    certsDir,
		reconnectCh: make(chan struct{}, 1),
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	go client.watchConnection()

	return client, nil
}

func (c *SliverClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	certFile := filepath.Join(c.certsDir, "client.crt")
	keyFile := filepath.Join(c.certsDir, "client.key")
	caFile := filepath.Join(c.certsDir, "ca.crt")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load client certificate: %w", err)
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to add CA certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	kap := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.Dial(
		c.address,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithKeepaliveParams(kap),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to sliver server: %w", err)
	}

	c.conn = conn
	c.connected = true

	return nil
}

func (c *SliverClient) watchConnection() {
	for range c.reconnectCh {
		for {
			time.Sleep(5 * time.Second)
			if err := c.connect(); err == nil {
				break
			}
		}
	}
}

func (c *SliverClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.reconnectCh)

	if c.conn != nil {
		return c.conn.Close()
	}
	c.connected = false
	return nil
}

func (c *SliverClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.conn != nil
}

func (c *SliverClient) GetConn() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

func (c *SliverClient) NotifyReconnect() {
	select {
	case c.reconnectCh <- struct{}{}:
	default:
	}
}

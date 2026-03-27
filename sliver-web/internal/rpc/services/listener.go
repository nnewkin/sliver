package services

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	rpcpb "github.com/BishopFox/sliver/protobuf/rpcpb"
)

type ListenerService struct {
	client *grpc.ClientConn
}

func NewListenerService(conn *grpc.ClientConn) *ListenerService {
	return &ListenerService{client: conn}
}

func (s *ListenerService) GetClient() rpcpb.SliverRPCClient {
	return rpcpb.NewSliverRPCClient(s.client)
}

func (s *ListenerService) GetAllListeners(ctx context.Context) ([]*Listener, error) {
	client := s.GetClient()
	resp, err := client.GetJobs(ctx, &rpcpb.JobsReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}

	listeners := make([]*Listener, 0, len(resp.GetJobs()))
	for _, j := range resp.GetJobs() {
		listeners = append(listeners, sliverJobToModel(j))
	}

	return listeners, nil
}

func (s *ListenerService) StartMTLSListener(ctx context.Context, name, bindAddr string) (*Listener, error) {
	client := s.GetClient()
	resp, err := client.StartMTLSListener(ctx, &rpcpb.MTLSListenerReq{
		Name:     name,
		BindAddr: bindAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start MTLS listener: %w", err)
	}

	return sliverJobToModel(resp.GetJob()), nil
}

func (s *ListenerService) StartDNSListener(ctx context.Context, name, bindAddr string, domains []string) (*Listener, error) {
	client := s.GetClient()
	resp, err := client.StartDNSListener(ctx, &rpcpb.DNSListenerReq{
		Name:     name,
		BindAddr: bindAddr,
		Domains:  domains,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start DNS listener: %w", err)
	}

	return sliverJobToModel(resp.GetJob()), nil
}

func (s *ListenerService) StartHTTPListener(ctx context.Context, name, bindAddr string, domains []string) (*Listener, error) {
	client := s.GetClient()
	resp, err := client.StartHTTPListener(ctx, &rpcpb.HTTPListenerReq{
		Name:     name,
		BindAddr: bindAddr,
		Domains:  domains,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start HTTP listener: %w", err)
	}

	return sliverJobToModel(resp.GetJob()), nil
}

func (s *ListenerService) StartHTTPSListener(ctx context.Context, name, bindAddr string, domains []string) (*Listener, error) {
	client := s.GetClient()
	resp, err := client.StartHTTPSListener(ctx, &rpcpb.HTTPSListenerReq{
		Name:     name,
		BindAddr: bindAddr,
		Domains:  domains,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start HTTPS listener: %w", err)
	}

	return sliverJobToModel(resp.GetJob()), nil
}

func (s *ListenerService) StartWGListener(ctx context.Context, name, bindAddr string, port int) (*Listener, error) {
	client := s.GetClient()
	resp, err := client.StartWGListener(ctx, &rpcpb.WGListenerReq{
		Name:     name,
		BindAddr: bindAddr,
		Port:     int32(port),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start WireGuard listener: %w", err)
	}

	return sliverJobToModel(resp.GetJob()), nil
}

func (s *ListenerService) StopListener(ctx context.Context, listenerID int64) error {
	client := s.GetClient()
	_, err := client.KillJob(ctx, &rpcpb.KillReq{
		JobID: listenerID,
	})
	if err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}
	return nil
}

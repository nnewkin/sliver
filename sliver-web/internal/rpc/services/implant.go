package services

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	rpcpb "github.com/BishopFox/sliver/protobuf/rpcpb"
)

type ImplantService struct {
	client *grpc.ClientConn
}

func NewImplantService(conn *grpc.ClientConn) *ImplantService {
	return &ImplantService{client: conn}
}

func (s *ImplantService) GetClient() rpcpb.SliverRPCClient {
	return rpcpb.NewSliverRPCClient(s.client)
}

type ImplantProfile struct {
	Name   string         `json:"name"`
	Config *ImplantConfig `json:"config"`
}

type ImplantConfig struct {
	Format         string   `json:"format"`
	GOOS           string   `json:"goos"`
	GOARCH         string   `json:"goarch"`
	MTLSAddresses  []string `json:"mtls_addresses"`
	DNSAddresses   []string `json:"dns_addresses"`
	HTTPAddresses  []string `json:"http_addresses"`
	HTTPSAddresses []string `json:"https_addresses"`
	WGAddresses    []string `json:"wg_addresses"`
	Interval       int32    `json:"interval"`
	Jitter         int32    `json:"jitter"`
	Reconnect      int32    `json:"reconnect"`
	MaxConns       int32    `json:"max_conns"`
	Name           string   `json:"name"`
	Limit          int32    `json:"limit"`
	LimitWait      int32    `json:"limit_wait"`
	Evasion        bool     `json:"evasion"`
	Obfuscate      bool     `json:"obfuscate"`
	OS             string   `json:"os"`
	Arch           string   `json:"arch"`
}

func (s *ImplantService) ListProfiles(ctx context.Context) ([]*ImplantProfile, error) {
	client := s.GetClient()
	resp, err := client.ImplantProfiles(ctx, &rpcpb.ImplantProfilesReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	profiles := make([]*ImplantProfile, 0, len(resp.GetProfiles()))
	for _, p := range resp.GetProfiles() {
		profiles = append(profiles, &ImplantProfile{
			Name: p.GetName(),
			Config: &ImplantConfig{
				Format:         p.GetConfig().GetFormat().String(),
				GOOS:           p.GetConfig().GetGOOS(),
				GOARCH:         p.GetConfig().GetGOARCH(),
				MTLSAddresses:  p.GetConfig().GetMTLSAddresses(),
				DNSAddresses:   p.GetConfig().GetDNSAddresses(),
				HTTPAddresses:  p.GetConfig().GetHTTPAddresses(),
				HTTPSAddresses: p.GetConfig().GetHTTPSAddresses(),
				WGAddresses:    p.GetConfig().GetWGAddresses(),
				Interval:       p.GetConfig().GetInterval(),
				Jitter:         p.GetConfig().GetJitter(),
				Reconnect:      p.GetConfig().GetReconnect(),
				MaxConns:       p.GetConfig().GetMaxConns(),
				Name:           p.GetConfig().GetName(),
				Limit:          p.GetConfig().GetLimit(),
				LimitWait:      p.GetConfig().GetLimitWait(),
				Evasion:        p.GetConfig().GetEvasion(),
				Obfuscate:      p.GetConfig().GetObfuscate(),
				OS:             p.GetConfig().GetOs(),
				Arch:           p.GetConfig().GetArch(),
			},
		})
	}

	return profiles, nil
}

func (s *ImplantService) GenerateImplant(ctx context.Context, cfg *ImplantConfig) ([]byte, error) {
	client := s.GetClient()

	config := &rpcpb.ImplantConfig{
		Format:         rpcpb.ImplantFormat(rpcpb.ImplantFormat_value[cfg.Format]),
		GOOS:           cfg.GOOS,
		GOARCH:         cfg.GOARCH,
		MTLSAddresses:  cfg.MTLSAddresses,
		DNSAddresses:   cfg.DNSAddresses,
		HTTPAddresses:  cfg.HTTPAddresses,
		HTTPSAddresses: cfg.HTTPSAddresses,
		WGAddresses:    cfg.WGAddresses,
		Interval:       cfg.Interval,
		Jitter:         cfg.Jitter,
		Reconnect:      cfg.Reconnect,
		MaxConns:       cfg.MaxConns,
		Name:           cfg.Name,
		Limit:          cfg.Limit,
		LimitWait:      cfg.LimitWait,
		Evasion:        cfg.Evasion,
		Obfuscate:      cfg.Obfuscate,
	}

	resp, err := client.Generate(ctx, &rpcpb.GenerateReq{
		Config: config,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate implant: %w", err)
	}

	return resp.GetFile().GetData(), nil
}

func (s *ImplantService) GenerateFromProfile(ctx context.Context, profileName string) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.GenerateFromProfile(ctx, &rpcpb.GenerateFromProfileReq{
		ProfileName: profileName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate from profile: %w", err)
	}

	return resp.GetFile().GetData(), nil
}

func (s *ImplantService) SaveProfile(ctx context.Context, name string, cfg *ImplantConfig) error {
	client := s.GetClient()

	config := &rpcpb.ImplantConfig{
		Format:         rpcpb.ImplantFormat(rpcpb.ImplantFormat_value[cfg.Format]),
		GOOS:           cfg.GOOS,
		GOARCH:         cfg.GOARCH,
		MTLSAddresses:  cfg.MTLSAddresses,
		DNSAddresses:   cfg.DNSAddresses,
		HTTPAddresses:  cfg.HTTPAddresses,
		HTTPSAddresses: cfg.HTTPSAddresses,
		WGAddresses:    cfg.WGAddresses,
		Interval:       cfg.Interval,
		Jitter:         cfg.Jitter,
		Reconnect:      cfg.Reconnect,
		MaxConns:       cfg.MaxConns,
		Name:           cfg.Name,
		Limit:          cfg.Limit,
		LimitWait:      cfg.LimitWait,
		Evasion:        cfg.Evasion,
		Obfuscate:      cfg.Obfuscate,
	}

	_, err := client.SaveImplantProfile(ctx, &rpcpb.SaveImplantProfileReq{
		Name:   name,
		Config: config,
	})
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	return nil
}

func (s *ImplantService) DeleteProfile(ctx context.Context, name string) error {
	client := s.GetClient()
	_, err := client.DeleteImplantProfile(ctx, &rpcpb.DeleteImplantProfileReq{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}

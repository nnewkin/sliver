package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sliverpb "github.com/BishopFox/sliver/protobuf/client"
	rpcpb "github.com/BishopFox/sliver/protobuf/rpcpb"
)

type SessionService struct {
	client *grpc.ClientConn
}

func NewSessionService(conn *grpc.ClientConn) *SessionService {
	return &SessionService{client: conn}
}

func (s *SessionService) GetClient() rpcpb.SliverRPCClient {
	return rpcpb.NewSliverRPCClient(s.client)
}

func (s *SessionService) GetAllSessions(ctx context.Context) ([]*Session, error) {
	client := s.GetClient()
	resp, err := client.GetSessions(ctx, &rpcpb.SessionsReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	sessions := make([]*Session, 0, len(resp.GetSessions()))
	for _, sliverSession := range resp.GetSessions() {
		sessions = append(sessions, sliverSessionToModel(sliverSession))
	}

	return sessions, nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	client := s.GetClient()
	resp, err := client.GetSession(ctx, &rpcpb.SessionReq{
		ID: sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return sliverSessionToModel(resp.GetSession()), nil
}

func (s *SessionService) KillSession(ctx context.Context, sessionID string) error {
	client := s.GetClient()
	_, err := client.KillSession(ctx, &rpcpb.KillReq{
		SessionID: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to kill session: %w", err)
	}
	return nil
}

func (s *SessionService) RenameSession(ctx context.Context, sessionID, newName string) (*Session, error) {
	client := s.GetClient()
	resp, err := client.RenameSession(ctx, &rpcpb.RenameReq{
		SessionID: sessionID,
		NewName:   newName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to rename session: %w", err)
	}

	return sliverSessionToModel(resp.GetSession()), nil
}

type Session struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Hostname          string    `json:"hostname"`
	Username          string    `json:"username"`
	UUID              string    `json:"uuid"`
	OS                string    `json:"os"`
	Arch              string    `json:"arch"`
	Transport         string    `json:"transport"`
	RemoteAddress     string    `json:"remote_address"`
	LocalAddress      string    `json:"local_address"`
	PID               int32     `json:"pid"`
	ProxyUUID         string    `json:"proxy_uuid"`
	ReconnectInterval int32     `json:"reconnect_interval"`
	LastCheckin       time.Time `json:"last_checkin"`
	NextCheckin       time.Time `json:"next_checkin"`
	IsActive          bool      `json:"is_active"`
	Version           string    `json:"version"`
	_revision         string    `json:"_revision,omitempty"`
}

func sliverSessionToModel(s *sliverpb.Session) *Session {
	if s == nil {
		return nil
	}
	return &Session{
		ID:                s.ID,
		Name:              s.Name,
		Hostname:          s.Hostname,
		Username:          s.Username,
		UUID:              s.UUID,
		OS:                s.OS,
		Arch:              s.Arch,
		Transport:         s.Transport,
		RemoteAddress:     s.RemoteAddress,
		LocalAddress:      s.LocalAddress,
		PID:               s.PID,
		ProxyUUID:         s.ProxyUUID,
		ReconnectInterval: s.ReconnectInterval,
		LastCheckin:       time.Unix(s.LastCheckin, 0),
		NextCheckin:       time.Unix(s.NextCheckin, 0),
		IsActive:          s.Active,
		Version:           s.Version,
		_revision:         s.Revision,
	}
}

type Beacon struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Hostname          string    `json:"hostname"`
	Username          string    `json:"username"`
	UUID              string    `json:"uuid"`
	OS                string    `json:"os"`
	Arch              string    `json:"arch"`
	RemoteAddress     string    `json:"remote_address"`
	PID               int32     `json:"pid"`
	LastCheckin       time.Time `json:"last_checkin"`
	NextCheckin       time.Time `json:"next_checkin"`
	Interval          int32     `json:"interval"`
	Jitter            int32     `json:"jitter"`
	ReconnectInterval int32     `json:"reconnect_interval"`
	IsActive          bool      `json:"is_active"`
	Version           string    `json:"version"`
}

func sliverBeaconToModel(b *sliverpb.Beacon) *Beacon {
	if b == nil {
		return nil
	}
	return &Beacon{
		ID:                b.ID,
		Name:              b.Name,
		Hostname:          b.Hostname,
		Username:          b.Username,
		UUID:              b.UUID,
		OS:                b.OS,
		Arch:              b.Arch,
		RemoteAddress:     b.RemoteAddress,
		PID:               b.PID,
		LastCheckin:       time.Unix(b.LastCheckin, 0),
		NextCheckin:       time.Unix(b.NextCheckin, 0),
		Interval:          b.Interval,
		Jitter:            b.Jitter,
		ReconnectInterval: b.ReconnectInterval,
		IsActive:          b.Active,
		Version:           b.Version,
	}
}

type BeaconTask struct {
	ID        string    `json:"id"`
	BeaconID  string    `json:"beacon_id"`
	Type      string    `json:"type"`
	Command   []string  `json:"command"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Completed bool      `json:"completed"`
	Success   bool      `json:"success"`
}

func sliverTaskToModel(t *sliverpb.Task) *BeaconTask {
	if t == nil {
		return nil
	}
	return &BeaconTask{
		ID:        t.ID,
		BeaconID:  t.BeaconID,
		Type:      t.Type,
		Command:   t.Command,
		CreatedAt: time.Unix(t.CreatedAt, 0),
		UpdatedAt: time.Unix(t.UpdatedAt, 0),
		Completed: t.Completed,
		Success:   t.Success,
	}
}

type Listener struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	BindAddress   string `json:"bind_address"`
	Port          int32  `json:"port"`
	Domain        string `json:"domain"`
	Options       string `json:"options"`
	Connected     int32  `json:"connected"`
	MaxConnection int32  `json:"max_connection"`
}

func sliverJobToModel(j *sliverpb.Job) *Listener {
	if j == nil {
		return nil
	}
	return &Listener{
		ID:            fmt.Sprintf("%d", j.ID),
		Name:          j.Name,
		Type:          j.Type,
		Status:        getJobStatus(j),
		BindAddress:   j.BindAddress,
		Port:          j.Port,
		Domain:        j.Domain,
		Options:       j.Options,
		Connected:     j.Connected,
		MaxConnection: j.MaxConnection,
	}
}

func getJobStatus(j *sliverpb.Job) string {
	if j.Active {
		return "active"
	}
	return "stopped"
}

func MapGrpcError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return fmt.Errorf("authentication failed")
	case codes.PermissionDenied:
		return fmt.Errorf("permission denied")
	case codes.NotFound:
		return fmt.Errorf("resource not found")
	case codes.AlreadyExists:
		return fmt.Errorf("resource already exists")
	case codes.InvalidArgument:
		return fmt.Errorf("invalid argument: %s", st.Message())
	case codes.DeadlineExceeded:
		return fmt.Errorf("operation timeout")
	case codes.Unavailable:
		return fmt.Errorf("service unavailable")
	default:
		return fmt.Errorf("rpc error: %s", st.Message())
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

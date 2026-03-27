package services

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	rpcpb "github.com/BishopFox/sliver/protobuf/rpcpb"
)

type InteractiveService struct {
	client *grpc.ClientConn
}

func NewInteractiveService(conn *grpc.ClientConn) *InteractiveService {
	return &InteractiveService{client: conn}
}

func (s *InteractiveService) GetClient() rpcpb.SliverRPCClient {
	return rpcpb.NewSliverRPCClient(s.client)
}

type Process struct {
	PID          int32  `json:"pid"`
	PPID         int32  `json:"ppid"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Architecture string `json:"architecture"`
	SessionID    int64  `json:"session_id"`
	CommandLine  string `json:"command_line"`
	Path         string `json:"path"`
	User         string `json:"user"`
	Started      string `json:"started"`
}

type PrivilegeInfo struct {
	User          string `json:"user"`
	Authenticated bool   `json:"authenticated"`
	Admin         bool   `json:"admin"`
	Impersonating bool   `json:"impersonating"`
	Username      string `json:"username"`
}

func (s *InteractiveService) GetProcesses(ctx context.Context, sessionID string) ([]*Process, error) {
	client := s.GetClient()
	resp, err := client.Ps(ctx, &rpcpb.PsReq{
		SessionID: sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	processes := make([]*Process, 0, len(resp.GetProcs()))
	for _, p := range resp.GetProcs() {
		processes = append(processes, &Process{
			PID:          p.GetPid(),
			PPID:         p.GetPpid(),
			Name:         p.GetName(),
			Owner:        p.GetOwner(),
			Architecture: p.GetArchitecture(),
			SessionID:    p.GetSessionID(),
			CommandLine:  p.GetCommandLine(),
			Path:         p.GetPath(),
			User:         p.GetUser(),
			Started:      p.GetStarted(),
		})
	}

	return processes, nil
}

func (s *InteractiveService) GetPrivileges(ctx context.Context, sessionID string) (*PrivilegeInfo, error) {
	client := s.GetClient()
	resp, err := client.GetPrivs(ctx, &rpcpb.GetPrivsReq{
		SessionID: sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get privileges: %w", err)
	}

	privs := resp.GetPrivs()
	if len(privs) == 0 {
		return &PrivilegeInfo{}, nil
	}

	priv := privs[0]
	return &PrivilegeInfo{
		User:          priv.GetUser(),
		Authenticated: priv.GetAuthenticated(),
		Admin:         priv.GetAdmin(),
		Impersonating: priv.GetImpersonating(),
		Username:      priv.GetUsername(),
	}, nil
}

func (s *InteractiveService) Execute(ctx context.Context, sessionID string, command string) (string, error) {
	client := s.GetClient()
	resp, err := client.Execute(ctx, &rpcpb.ExecuteReq{
		SessionID: sessionID,
		Command:   command,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) ExecuteAssembly(ctx context.Context, sessionID, assembly string, processID int32, args string) (string, error) {
	client := s.GetClient()
	resp, err := client.ExecuteAssembly(ctx, &rpcpb.ExecuteAssemblyReq{
		SessionID:   sessionID,
		Assembly:    []byte(assembly),
		ProcessID:   processID,
		AssemblyArg: args,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute assembly: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) ExecuteToken(ctx context.Context, sessionID string, processID int32, command string) (string, error) {
	client := s.GetClient()
	resp, err := client.ExecuteToken(ctx, &rpcpb.ExecuteTokenReq{
		SessionID: sessionID,
		TokenID:   processID,
		Command:   command,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute with token: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) OpenShell(ctx context.Context, sessionID string, cols, rows uint32) (io.Reader, error) {
	client := s.GetClient()

	req := &rpcpb.ShellReq{
		SessionID: sessionID,
		Cols:      cols,
		Rows:      rows,
	}

	stream, err := client.Shell(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to open shell: %w", err)
	}

	return &shellReader{stream: stream}, nil
}

type shellReader struct {
	stream rpcpb.SliverRPC_ShellClient
}

func (r *shellReader) Read(p []byte) (n int, err error) {
	resp, err := r.stream.Recv()
	if err != nil {
		if err == io.EOF {
			return 0, io.EOF
		}
		return 0, err
	}

	data := resp.GetResponse()
	if len(data) > len(p) {
		copy(p, data)
		return len(data), io.ErrShortBuffer
	}

	copy(p, data)
	return len(data), nil
}

func (r *shellReader) Close() error {
	return r.stream.CloseSend()
}

func (s *InteractiveService) Shell(ctx context.Context, sessionID string, data string, cols, rows uint32) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.Shell(ctx, &rpcpb.ShellReq{
		SessionID: sessionID,
		Data:      data,
		Cols:      cols,
		Rows:      rows,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send shell data: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) Screenshot(ctx context.Context, sessionID string) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.Screenshot(ctx, &rpcpb.ScreenshotReq{
		SessionID: sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to take screenshot: %w", err)
	}

	return resp.GetImage(), nil
}

func (s *InteractiveService) Upload(ctx context.Context, sessionID, remotePath string, data []byte) error {
	client := s.GetClient()
	_, err := client.Upload(ctx, &rpcpb.UploadReq{
		SessionID: sessionID,
		Path:      remotePath,
		Data:      data,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (s *InteractiveService) Download(ctx context.Context, sessionID, remotePath string) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.Download(ctx, &rpcpb.DownloadReq{
		SessionID: sessionID,
		Path:      remotePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return resp.GetData(), nil
}

func (s *InteractiveService) ListDirectory(ctx context.Context, sessionID, path string) ([]*FileInfo, error) {
	client := s.GetClient()
	resp, err := client.Ls(ctx, &rpcpb.LsReq{
		SessionID: sessionID,
		Path:      path,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	files := make([]*FileInfo, 0, len(resp.GetPath()))
	for _, f := range resp.GetPath() {
		files = append(files, &FileInfo{
			Name:     f.GetName(),
			Size:     f.GetSize(),
			IsDir:    f.GetIsDir(),
			Mode:     f.GetMode(),
			Modified: f.GetModified(),
			Perms:    f.GetPerms(),
		})
	}

	return files, nil
}

type FileInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Mode     uint32 `json:"mode"`
	Modified string `json:"modified"`
	Perms    string `json:"perms"`
}

func (s *InteractiveService) Mkdir(ctx context.Context, sessionID, path string) error {
	client := s.GetClient()
	_, err := client.Mkdir(ctx, &rpcpb.MkdirReq{
		SessionID: sessionID,
		Path:      path,
	})
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func (s *InteractiveService) Rm(ctx context.Context, sessionID, path string) error {
	client := s.GetClient()
	_, err := client.Rm(ctx, &rpcpb.RmReq{
		SessionID: sessionID,
		Path:      path,
	})
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}

func (s *InteractiveService) Cp(ctx context.Context, sessionID, src, dst string) error {
	client := s.GetClient()
	_, err := client.Cp(ctx, &rpcpb.CpReq{
		SessionID: sessionID,
		Src:       src,
		Dst:       dst,
	})
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

func (s *InteractiveService) Mv(ctx context.Context, sessionID, src, dst string) error {
	client := s.GetClient()
	_, err := client.Mv(ctx, &rpcpb.MvReq{
		SessionID: sessionID,
		Src:       src,
		Dst:       dst,
	})
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}

func (s *InteractiveService) ProcessDump(ctx context.Context, sessionID string, pid int32) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.ProcessDump(ctx, &rpcpb.ProcessDumpReq{
		SessionID: sessionID,
		PID:       pid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dump process: %w", err)
	}

	return resp.GetFile().GetData(), nil
}

func (s *InteractiveService) RunAs(ctx context.Context, sessionID, username, program string) (string, error) {
	client := s.GetClient()
	resp, err := client.RunAs(ctx, &rpcpb.RunAsReq{
		SessionID: sessionID,
		Username:  username,
		Program:   program,
	})
	if err != nil {
		return "", fmt.Errorf("failed to run as: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) Impersonate(ctx context.Context, sessionID, username string) error {
	client := s.GetClient()
	_, err := client.Impersonate(ctx, &rpcpb.ImpersonateReq{
		SessionID: sessionID,
		Username:  username,
	})
	if err != nil {
		return fmt.Errorf("failed to impersonate: %w", err)
	}
	return nil
}

func (s *InteractiveService) RevToSelf(ctx context.Context, sessionID string) error {
	client := s.GetClient()
	_, err := client.RevToSelf(ctx, &rpcpb.RevToSelfReq{
		SessionID: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to revert to self: %w", err)
	}
	return nil
}

func (s *InteractiveService) GetSystem(ctx context.Context, sessionID string) (string, error) {
	client := s.GetClient()
	resp, err := client.GetSystem(ctx, &rpcpb.GetSystemReq{
		SessionID: sessionID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get system: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *InteractiveService) MakeToken(ctx context.Context, sessionID, username, password string) error {
	client := s.GetClient()
	_, err := client.MakeToken(ctx, &rpcpb.MakeTokenReq{
		SessionID: sessionID,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return fmt.Errorf("failed to make token: %w", err)
	}
	return nil
}

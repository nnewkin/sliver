package services

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	rpcpb "github.com/BishopFox/sliver/protobuf/rpcpb"
)

type BeaconService struct {
	client *grpc.ClientConn
}

func NewBeaconService(conn *grpc.ClientConn) *BeaconService {
	return &BeaconService{client: conn}
}

func (s *BeaconService) GetClient() rpcpb.SliverRPCClient {
	return rpcpb.NewSliverRPCClient(s.client)
}

func (s *BeaconService) GetAllBeacons(ctx context.Context) ([]*Beacon, error) {
	client := s.GetClient()
	resp, err := client.GetBeacons(ctx, &rpcpb.BeaconsReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to get beacons: %w", err)
	}

	beacons := make([]*Beacon, 0, len(resp.GetBeacons()))
	for _, b := range resp.GetBeacons() {
		beacons = append(beacons, sliverBeaconToModel(b))
	}

	return beacons, nil
}

func (s *BeaconService) GetBeacon(ctx context.Context, beaconID string) (*Beacon, error) {
	client := s.GetClient()
	resp, err := client.GetBeacon(ctx, &rpcpb.BeaconReq{
		ID: beaconID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon: %w", err)
	}

	return sliverBeaconToModel(resp.GetBeacon()), nil
}

func (s *BeaconService) GetBeaconTasks(ctx context.Context, beaconID string) ([]*BeaconTask, error) {
	client := s.GetClient()
	resp, err := client.GetBeaconTasks(ctx, &rpcpb.BeaconTasksReq{
		BeaconID: beaconID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get beacon tasks: %w", err)
	}

	tasks := make([]*BeaconTask, 0, len(resp.GetTasks()))
	for _, t := range resp.GetTasks() {
		tasks = append(tasks, sliverTaskToModel(t))
	}

	return tasks, nil
}

func (s *BeaconService) GetTaskOutput(ctx context.Context, beaconID, taskID string) ([]byte, error) {
	client := s.GetClient()
	resp, err := client.GetTaskOutput(ctx, &rpcpb.TaskReq{
		BeaconID: beaconID,
		TaskID:   taskID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get task output: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *BeaconService) CancelTask(ctx context.Context, beaconID, taskID string) error {
	client := s.GetClient()
	_, err := client.CancelTask(ctx, &rpcpb.TaskReq{
		BeaconID: beaconID,
		TaskID:   taskID,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel task: %w", err)
	}
	return nil
}

func (s *BeaconService) ExecuteCommand(ctx context.Context, beaconID string, command []string) (string, error) {
	client := s.GetClient()
	resp, err := client.Execute(ctx, &rpcpb.ExecuteReq{
		BeaconID: beaconID,
		Command:  command,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return resp.GetResponse(), nil
}

func (s *BeaconService) SetBeaconOptions(ctx context.Context, beaconID string, interval, jitter int32) error {
	client := s.GetClient()
	_, err := client.BeaconOptions(ctx, &rpcpb.BeaconOptions{
		BeaconID: beaconID,
		Interval: interval,
		Jitter:   jitter,
	})
	if err != nil {
		return fmt.Errorf("failed to set beacon options: %w", err)
	}
	return nil
}

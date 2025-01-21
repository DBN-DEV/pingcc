package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/DBN-DEV/pingpb/gopb"

	"github.com/DBN-DEV/pingcc/domain"
	"github.com/DBN-DEV/pingcc/log"
)

var _ gopb.ControllerServer = &Controller{}

type Controller struct {
	gopb.UnimplementedControllerServer

	agentRepo       AgentRepo
	pingTaskRepo    PingTaskRepo
	tcpPingTaskRepo TcpPingTaskRepo

	logger *zap.Logger
}

func (c *Controller) Heartbeat(ctx context.Context, req *gopb.GrpcHeartbeatReq) (*gopb.GrpcHeartbeatResp, error) {
	c.logger.Debug("Heartbeat", log.AgentUID(req.AgentUID), zap.Uint64("ping_task_version", req.PingTaskVersion), zap.Uint64("tcp_ping_task_version", req.TcpPingTaskVersion))
	agent, err := c.agentRepo.Find(ctx, req.AgentUID)
	if err != nil {
		c.logger.Warn("Failed to find agent", zap.Error(err), log.AgentUID(req.AgentUID))
		return nil, err
	}

	var resp gopb.GrpcHeartbeatResp

	resp.NeedUpdatePingTask = agent.PingTaskVersion != req.PingTaskVersion
	resp.NeedUpdateTcpPingTask = agent.TcpPingTaskVersion != req.TcpPingTaskVersion
	c.logger.Debug("Heartbeat response", log.AgentUID(req.AgentUID), zap.Bool("need_update_ping_task", resp.NeedUpdatePingTask), zap.Bool("need_update_tcp_ping_task", resp.NeedUpdateTcpPingTask))
	if err := c.agentRepo.Heartbeat(ctx, req.AgentUID); err != nil {
		c.logger.Warn("Failed to update agent last active time", zap.Error(err), log.AgentUID(req.AgentUID))
	}

	return &resp, nil
}

func convertTcpPingTaskToGrpc(task domain.TcpPingTask) *gopb.GrpcTcpPingTask {
	return &gopb.GrpcTcpPingTask{
		TcpPingTaskUID: task.UID,
		Dest:           task.Dest,
		Src:            task.Src,
		TimeoutMS:      task.TimeoutMS,
		IntervalMS:     task.IntervalMS,
	}
}

func (c *Controller) GetTcpPingTask(ctx context.Context, req *gopb.GrpcTaskReq) (*gopb.GrpcTcpPingTaskResp, error) {
	tasks, err := c.tcpPingTaskRepo.FindByAgentID(ctx, req.AgentUID)
	if err != nil {
		c.logger.Warn("Failed to find tcp ping task", zap.Error(err), log.AgentUID(req.AgentUID))
		return nil, err
	}

	agent, err := c.agentRepo.Find(ctx, req.AgentUID)
	if err != nil {
		c.logger.Warn("Failed to find agent", zap.Error(err), log.AgentUID(req.AgentUID))
		return nil, err
	}

	resp := &gopb.GrpcTcpPingTaskResp{Version: agent.TcpPingTaskVersion}
	grpcTasks := make([]*gopb.GrpcTcpPingTask, 0, len(tasks))
	for _, task := range tasks {
		grpcTasks = append(grpcTasks, convertTcpPingTaskToGrpc(task))
	}
	resp.TcpPingTasks = grpcTasks
	return resp, nil
}

func convertPingTaskToGrpc(task domain.PingTask) *gopb.GrpcPingTask {
	return &gopb.GrpcPingTask{
		PingTaskUID: task.UID,
		Dest:        task.Dest,
		Src:         task.Src,
		TimeoutMS:   task.TimeoutMS,
		IntervalMS:  task.IntervalMS,
	}
}

func (c *Controller) GetPingTask(ctx context.Context, req *gopb.GrpcTaskReq) (*gopb.GrpcPingTaskResp, error) {
	tasks, err := c.pingTaskRepo.FindByAgentID(ctx, req.AgentUID)
	if err != nil {
		c.logger.Warn("Failed to find tcp ping task", zap.Error(err), log.AgentUID(req.AgentUID))
		return nil, err
	}
	agent, err := c.agentRepo.Find(ctx, req.AgentUID)
	if err != nil {
		c.logger.Warn("Failed to find agent", zap.Error(err), log.AgentUID(req.AgentUID))
		return nil, err
	}

	resp := &gopb.GrpcPingTaskResp{Version: agent.PingTaskVersion}
	grpcTasks := make([]*gopb.GrpcPingTask, 0, len(tasks))
	for _, task := range tasks {
		grpcTasks = append(grpcTasks, convertPingTaskToGrpc(task))
	}
	resp.PingTasks = grpcTasks
	return resp, nil
}

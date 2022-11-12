package app

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"pingcc/domain"
	"pingcc/log"
	"pingcc/pb"
)

type Controller struct {
	pb.UnimplementedControllerServer

	agentRepo domain.AgentRepo
	chL       sync.RWMutex
	agentIDCh map[uint64]chan *pb.UpdateCommandResp

	logger *zap.Logger
}

func NewController(repo domain.AgentRepo) *Controller {
	i := &Controller{
		agentRepo: repo,
		agentIDCh: make(map[uint64]chan *pb.UpdateCommandResp),
		logger:    log.LWithSvcName("controller"),
	}

	return i
}

func (i *Controller) checkUpdate() {
	m := make(map[uint64]chan *pb.UpdateCommandResp)
	i.chL.RLock()
	for k, v := range m {
		m[k] = v
	}
	i.chL.RUnlock()

	for id, ch := range i.agentIDCh {
		agent, err := i.agentRepo.Find(context.Background(), id)
		if err != nil {
			i.logger.Warn("Can not get agent, on watch process", zap.Error(err))
			continue
		}
		if agent.AgentPingCommandVersion != agent.ControllerPingCommandVersion {
			ch <- &pb.UpdateCommandResp{
				CommandType: pb.CommandType_Ping,
				Version:     agent.ControllerPingCommandVersion,
			}
		}
		if agent.AgentTcpPingCommandVersion != agent.ControllerTcpPingCommandVersion {
			ch <- &pb.UpdateCommandResp{
				CommandType: pb.CommandType_TcpPing,
				Version:     agent.ControllerTcpPingCommandVersion,
			}
		}
	}
}

func (i *Controller) checkUpdateProc() {
	ticker := time.NewTimer(30 * time.Second)
	for range ticker.C {
		i.checkUpdate()
	}
}

func (i *Controller) Register(req *pb.RegisterReq, server pb.Controller_RegisterServer) error {
	agent, err := i.agentRepo.FindWithPingTargets(context.Background(), req.AgentID)
	if err != nil {
		i.logger.Info("Fail to find agent with ping targets", zap.Uint64("agent_id", req.AgentID))
		return err
	}

	if err := i.sendInitCommand(agent, server); err != nil {
		i.logger.Info("Fail to send init command", zap.Uint64("agent_id", req.AgentID), zap.Error(err))
		return err
	}

	ch := i.initCH(agent)

	for resp := range ch {
		if err := server.Send(resp); err != nil {
			i.logger.Info("Fail to send command", zap.Uint64("agent_id", req.AgentID), zap.Error(err))
			return err
		}
	}

	return nil
}

// 　sendInitCommand 给 agent 发送初始化指令
func (i *Controller) sendInitCommand(agent *domain.Agent, server pb.Controller_RegisterServer) error {
	resps := []*pb.UpdateCommandResp{{
		CommandType: pb.CommandType_Ping,
		Version:     agent.ControllerPingCommandVersion,
	}, {
		CommandType: pb.CommandType_TcpPing,
		Version:     agent.ControllerTcpPingCommandVersion,
	}}

	for _, resp := range resps {
		if err := server.Send(resp); err != nil {
			i.logger.Info("Fail to send command", zap.Uint64("agent_id", agent.ID), zap.Error(err))
			return err
		}
	}

	return nil
}

// initCh 给注册的 agent 初始化 ch
func (i *Controller) initCH(agent *domain.Agent) chan *pb.UpdateCommandResp {
	i.chL.Lock()
	defer i.chL.Unlock()

	if ch, ok := i.agentIDCh[agent.ID]; ok {
		close(ch)
		delete(i.agentIDCh, agent.ID)
	}

	ch := make(chan *pb.UpdateCommandResp, 2)
	i.agentIDCh[agent.ID] = ch
	return ch
}

func (i *Controller) GetTcpPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.TcpPingCommandResp, error) {
	agent, err := i.agentRepo.FindWithTcpPingTargets(ctx, req.AgentID)
	if err != nil {
		i.logger.Info("Fail to find agent with tcp ping targets", zap.Uint64("agent_id", agent.ID), zap.Error(err))
		return nil, err
	}
	agent.ActivateByGetTcpPingComm(req.Version)
	if err := i.agentRepo.Save(ctx, agent); err != nil {
		i.logger.Info("Fail to save agent", zap.Uint64("agent_id", agent.ID), zap.Error(err))
	}

	comms := make([]*pb.GrpcTcpPingCommand, 0, len(agent.TcpPingTargets))
	for _, target := range agent.TcpPingTargets {
		comm := &pb.GrpcTcpPingCommand{
			ID:         target.ID,
			Target:     target.Address,
			TimeoutMS:  target.TimeoutMS,
			IntervalMS: target.IntervalMS,
		}
		comms = append(comms, comm)
	}

	return &pb.TcpPingCommandResp{
		Version:         agent.ControllerTcpPingCommandVersion,
		TcpPingCommands: comms,
	}, nil
}

func (i *Controller) GetPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.PingCommandsResp, error) {
	agent, err := i.agentRepo.FindWithPingTargets(ctx, req.AgentID)
	if err != nil {
		return nil, err
	}
	agent.ActivateByGetPingComm(req.Version)
	if err := i.agentRepo.Save(ctx, agent); err != nil {
		i.logger.Info("Fail to save agent", zap.Uint64("agent_id", agent.ID), zap.Error(err))
	}

	comms := make([]*pb.GrpcPingCommand, 0, len(agent.PingTargets))
	for _, target := range agent.PingTargets {
		comm := &pb.GrpcPingCommand{
			ID:         target.ID,
			IP:         target.IP,
			TimeoutMS:  target.TimeoutMS,
			IntervalMS: target.IntervalMS,
		}
		comms = append(comms, comm)
	}

	return &pb.PingCommandsResp{
		Version:      agent.ControllerPingCommandVersion,
		PingCommands: comms,
	}, nil
}

func (i *Controller) GetFpingCommand(ctx context.Context, req *pb.CommandReq) (*pb.FpingCommandResp, error) {
	return nil, nil
}

func (i *Controller) GetMtrCommand(ctx context.Context, req *pb.CommandReq) (*pb.MtrCommandResp, error) {
	return nil, nil
}

package controller

import (
	"context"
	"go.uber.org/zap"
	"pingcc/log"
	"sync"

	"pingcc/domain"

	"pingcc/pb"
)

type impl struct {
	pb.UnimplementedControllerServer

	agentRepo domain.AgentRepo
	chL       sync.RWMutex
	agentIDCh map[uint]chan *pb.UpdateCommandResp
}

func New(repo domain.AgentRepo) *impl {
	i := impl{
		agentRepo: repo,
		agentIDCh: make(map[uint]chan *pb.UpdateCommandResp),
	}

	return &i
}

func (i *impl) Register(req *pb.RegisterReq, server pb.Controller_RegisterServer) error {
	agent, err := i.agentRepo.FindWithPingTargets(context.Background(), uint(req.AgentID))
	if err != nil {
		return err
	}

	if err := i.sendInitCommand(agent, server); err != nil {
		return err
	}

	ch := i.initCH(agent)

	for resp := range ch {
		if err := server.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//　sendInitCommand 给 agent 发送初始化指令
func (i *impl) sendInitCommand(agent *domain.Agent, server pb.Controller_RegisterServer) error {
	resps := []*pb.UpdateCommandResp{{
		CommandType: pb.CommandType_Ping,
		Version:     agent.ControllerPingCommandVersion,
	}, {
		CommandType: pb.CommandType_TcpPing,
		Version:     agent.ControllerTcpPingCommandVersion,
	}}

	for _, resp := range resps {
		if err := server.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// initCh 给注册的 agent 初始化 ch
func (i *impl) initCH(agent *domain.Agent) chan *pb.UpdateCommandResp {
	i.chL.Lock()
	defer i.chL.Unlock()

	if ch, ok := i.agentIDCh[agent.ID]; ok {
		close(ch)
		delete(i.agentIDCh, agent.ID)
	}

	ch := make(chan *pb.UpdateCommandResp)
	i.agentIDCh[agent.ID] = ch
	return ch
}

func (i *impl) GetTcpPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.TcpPingCommandResp, error) {
	agent, err := i.agentRepo.FindWithTcpPingTargets(ctx, uint(req.AgentID))
	if err != nil {
		return nil, err
	}
	agent.ActivateByGetTcpPingComm(req.Version)
	if err := i.agentRepo.Save(ctx, agent); err != nil {
		log.L().Info("Fail to save agent", zap.Error(err))
	}

	comms := make([]*pb.GrpcTcpPingCommand, 0, len(agent.TcpPingTargets))
	for _, target := range agent.TcpPingTargets {
		comm := &pb.GrpcTcpPingCommand{
			ID:         uint64(target.ID),
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

func (i *impl) GetPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.PingCommandsResp, error) {
	agent, err := i.agentRepo.FindWithPingTargets(ctx, uint(req.AgentID))
	if err != nil {
		return nil, err
	}
	agent.ActivateByGetPingComm(req.Version)
	if err := i.agentRepo.Save(ctx, agent); err != nil {
		log.L().Info("Fail to save agent", zap.Error(err))
	}

	comms := make([]*pb.GrpcPingCommand, 0, len(agent.PingTargets))
	for _, target := range agent.PingTargets {
		comm := &pb.GrpcPingCommand{
			ID:         uint64(target.ID),
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

func (i *impl) GetFpingCommand(ctx context.Context, req *pb.CommandReq) (*pb.FpingCommandResp, error) {
	return nil, nil
}

func (i *impl) GetMtrCommand(ctx context.Context, req *pb.CommandReq) (*pb.MtrCommandResp, error) {
	return nil, nil
}

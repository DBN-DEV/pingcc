package controller

import (
	"context"
	"pingcc/entry"
	"sync"

	"pingcc/pb"
)

type impl struct {
	pb.UnimplementedControllerServer

	repo      entry.AgentRepo
	chL       sync.RWMutex
	agentIDCh map[uint]chan *pb.UpdateCommandResp
}

func New(agentRepo entry.AgentRepo) pb.ControllerServer {
	i := impl{
		repo:      agentRepo,
		agentIDCh: make(map[uint]chan *pb.UpdateCommandResp),
	}

	return &i
}

func (i *impl) Register(req *pb.RegisterReq, server pb.Controller_RegisterServer) error {
	agent, err := i.repo.FindPreloadPingCommand(int(req.AgentID))
	if err != nil {
		return err
	}

	if err := i.sendInitCommand(agent, server); err != nil {
		return err
	}

	i.initCH(agent)

	return nil
}

//　sendInitCommand 给 agent 发送初始化指令
func (i *impl) sendInitCommand(agent *entry.Agent, server pb.Controller_RegisterServer) error {
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
func (i *impl) initCH(agent *entry.Agent) {
	i.chL.Lock()
	defer i.chL.Unlock()

	if ch, ok := i.agentIDCh[agent.ID]; ok {
		close(ch)
		delete(i.agentIDCh, agent.ID)
	}
	i.agentIDCh[agent.ID] = make(chan *pb.UpdateCommandResp)
}

func (i *impl) GetTcpPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.TcpPingCommandResp, error) {
	return nil, nil
}

func (i *impl) GetPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.PingCommandsResp, error) {
	agent, err := i.repo.FindPreloadPingCommand(int(req.AgentID))
	if err != nil {
		return nil, err
	}

	comms := make([]*pb.GrpcPingCommand, 0, len(agent.PingTargets))
	for _, target := range agent.PingTargets {
		comm := &pb.GrpcPingCommand{
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

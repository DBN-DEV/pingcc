package controller

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"ping-cc/pb"
)

type Impl struct {
	pb.UnimplementedControllerServer

	repo      agentRepo
	chL       sync.RWMutex
	agentIDCh map[uint]chan *pb.UpdateCommandResp
}

func New(db *gorm.DB) pb.ControllerServer {
	repo := &agentRepoImpl{db: db}

	i := Impl{
		repo:      repo,
		agentIDCh: make(map[uint]chan *pb.UpdateCommandResp),
	}

	return &i
}

func (c *Impl) Register(req *pb.RegisterReq, server pb.Controller_RegisterServer) error {
	agent, err := c.repo.findPreloadPingCommand(int(req.AgentID))
	if err != nil {
		return err
	}

	if err := c.sendInitCommand(agent, server); err != nil {
		return err
	}

	c.initCH(agent)

	return nil
}

//　sendInitCommand 给 agent 发送初始化指令
func (c *Impl) sendInitCommand(agent *Agent, server pb.Controller_RegisterServer) error {
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
func (c *Impl) initCH(agent *Agent) {
	c.chL.Lock()
	defer c.chL.Unlock()

	if ch, ok := c.agentIDCh[agent.ID]; ok {
		close(ch)
		delete(c.agentIDCh, agent.ID)
	}
	c.agentIDCh[agent.ID] = make(chan *pb.UpdateCommandResp)
}

func (c *Impl) GetTcpPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.TcpPingCommandResp, error) {
	return nil, nil
}

func (c *Impl) GetPingCommand(ctx context.Context, req *pb.CommandReq) (*pb.PingCommandsResp, error) {
	agent, err := c.repo.findPreloadPingCommand(int(req.AgentID))
	if err != nil {
		return nil, err
	}

	var commands []*pb.GrpcPingCommand
	for _, target := range agent.PingTargets {
		command := &pb.GrpcPingCommand{
			IP:         target.IP,
			TimeoutMS:  target.TimeoutMS,
			IntervalMS: target.IntervalMS,
		}
		commands = append(commands, command)
	}

	return &pb.PingCommandsResp{
		Version:      agent.ControllerPingCommandVersion,
		PingCommands: commands,
	}, nil
}

func (c *Impl) GetFpingCommand(ctx context.Context, req *pb.CommandReq) (*pb.FpingCommandResp, error) {
	return nil, nil
}

func (c *Impl) GetMtrCommand(ctx context.Context, req *pb.CommandReq) (*pb.MtrCommandResp, error) {
	return nil, nil
}

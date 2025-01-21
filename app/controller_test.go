package app

import (
	"context"
	"github.com/DBN-DEV/pingcc/mocks"
	"testing"

	"github.com/DBN-DEV/pingpb/gopb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/DBN-DEV/pingcc/domain"
)

func newTestController() *Controller {
	return &Controller{logger: zap.NewNop()}
}

func TestController_Heartbeat(t *testing.T) {
	grpcReq := &gopb.GrpcHeartbeatReq{
		AgentUID:           1,
		PingTaskVersion:    1,
		TcpPingTaskVersion: 1,
	}

	agent := &domain.Agent{
		AgentUID:           1,
		PingTaskVersion:    2,
		TcpPingTaskVersion: 2,
	}

	agentRepo := mocks.NewAgentRepo(t)
	agentRepo.EXPECT().Find(mock.Anything, uint64(1)).Return(agent, nil).Once()
	agentRepo.EXPECT().Heartbeat(mock.Anything, uint64(1)).Return(nil).Once()

	c := newTestController()
	c.agentRepo = agentRepo

	resp, err := c.Heartbeat(context.Background(), grpcReq)
	assert.Nil(t, err)
	assert.True(t, resp.NeedUpdatePingTask)
	assert.True(t, resp.NeedUpdateTcpPingTask)
}

func TestController_GetTcpPingTask(t *testing.T) {
	grpcReq := &gopb.GrpcTaskReq{AgentUID: 1}

	tasks := []domain.TcpPingTask{{UID: 1, Dest: "dest", Src: "src", TimeoutMS: 1000, IntervalMS: 1000}}

	tcpPingTaskRepo := mocks.NewTcpPingTaskRepo(t)
	tcpPingTaskRepo.EXPECT().FindByAgentID(mock.Anything, uint64(1)).Return(tasks, nil).Once()
	agentRepo := mocks.NewAgentRepo(t)
	agentRepo.EXPECT().Find(mock.Anything, uint64(1)).Return(&domain.Agent{}, nil).Once()

	c := newTestController()
	c.tcpPingTaskRepo = tcpPingTaskRepo
	c.agentRepo = agentRepo

	resp, err := c.GetTcpPingTask(context.Background(), grpcReq)
	assert.Nil(t, err)
	assert.Len(t, resp.TcpPingTasks, 1)
	assert.Equal(t, uint64(1), resp.TcpPingTasks[0].TcpPingTaskUID)
	assert.Equal(t, "dest", resp.TcpPingTasks[0].Dest)
	assert.Equal(t, "src", resp.TcpPingTasks[0].Src)
	assert.Equal(t, uint32(1000), resp.TcpPingTasks[0].TimeoutMS)
	assert.Equal(t, uint32(1000), resp.TcpPingTasks[0].IntervalMS)
}

func TestController_GetPingTask(t *testing.T) {
	grpcReq := &gopb.GrpcTaskReq{AgentUID: 1}

	tasks := []domain.PingTask{{UID: 1, Dest: "dest", Src: "src", TimeoutMS: 1000, IntervalMS: 1000}}

	pingTaskRepo := mocks.NewPingTaskRepo(t)
	pingTaskRepo.EXPECT().FindByAgentID(mock.Anything, uint64(1)).Return(tasks, nil).Once()
	agentRepo := mocks.NewAgentRepo(t)
	agentRepo.EXPECT().Find(mock.Anything, uint64(1)).Return(&domain.Agent{}, nil).Once()

	c := newTestController()
	c.pingTaskRepo = pingTaskRepo
	c.agentRepo = agentRepo

	resp, err := c.GetPingTask(context.Background(), grpcReq)
	assert.Nil(t, err)
	assert.Len(t, resp.PingTasks, 1)
	assert.Equal(t, uint64(1), resp.PingTasks[0].PingTaskUID)
	assert.Equal(t, "dest", resp.PingTasks[0].Dest)
	assert.Equal(t, "src", resp.PingTasks[0].Src)
	assert.Equal(t, uint32(1000), resp.PingTasks[0].TimeoutMS)
	assert.Equal(t, uint32(1000), resp.PingTasks[0].IntervalMS)
}

package app

import (
	"context"
	"testing"

	"github.com/DBN-DEV/pingpb/gopb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/datatypes"

	"github.com/DBN-DEV/pingcc/domain"
	"github.com/DBN-DEV/pingcc/mocks"
)

func newTestCollector() *Collector {
	return &Collector{logger: zap.NewNop()}
}

func TestCollector_PingReport(t *testing.T) {
	grpcReq := &gopb.GrpcPingReportReq{
		Results:  []*gopb.GrpcPingResult{{RttMicros: 100, IsTimeout: false, PingTaskUID: 1}},
		AgentUID: 1,
	}
	tag := make(map[string]string)
	tag["key"] = "value"
	task := &domain.PingTask{Tag: datatypes.NewJSONType(tag)}
	pingTaskRepo := mocks.NewPingTaskRepo(t)
	pingTaskRepo.EXPECT().Find(mock.Anything, uint64(1)).Return(task, nil).Once()

	tsdb := mocks.NewTSDB(t)
	tsdb.EXPECT().InsertPingResult(context.Background(), domain.PingResult{RttMicros: 100, IsTimeout: false, Tag: tag}).
		Return(nil).Once()

	c := newTestCollector()
	c.tsdb = tsdb
	c.pingTaskRepo = pingTaskRepo

	_, err := c.PingReport(context.Background(), grpcReq)
	assert.Nil(t, err)
}

func TestCollector_TcpPingReport(t *testing.T) {
	grpcReq := &gopb.GrpcTcpPingReportReq{
		Results:  []*gopb.GrpcTcpPingResult{{RttMicros: 100, IsTimeout: false, TcpPingTaskUID: 1}},
		AgentUID: 1,
	}
	tag := make(map[string]string)
	tag["key"] = "value"
	task := &domain.TcpPingTask{Tag: datatypes.NewJSONType(tag)}
	tcpPingTaskRepo := mocks.NewTcpPingTaskRepo(t)
	tcpPingTaskRepo.EXPECT().Find(mock.Anything, uint64(1)).Return(task, nil).Once()

	tsdb := mocks.NewTSDB(t)
	tsdb.EXPECT().InsertTcpPingResult(context.Background(), domain.TcpPingResult{RttMicros: 100, IsTimeout: false, Tag: tag}).
		Return(nil).
		Once()

	c := newTestCollector()
	c.tsdb = tsdb
	c.tcpPingTaskRepo = tcpPingTaskRepo

	_, err := c.TcpPingReport(context.Background(), grpcReq)
	assert.Nil(t, err)
}

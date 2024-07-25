package app

import (
	"context"

	"github.com/DBN-DEV/pingpb/gopb"
	"go.uber.org/zap"

	"github.com/DBN-DEV/pingcc/domain"
	"github.com/DBN-DEV/pingcc/log"
)

type Collector struct {
	gopb.UnimplementedCollectorServer

	logger *zap.Logger

	tsdb TSDB

	pingTaskRepo    PingTaskRepo
	tcpPingTaskRepo TcpPingTaskRepo
}

func (i *Collector) PingReport(ctx context.Context, req *gopb.GrpcPingReportReq) (*gopb.Empty, error) {
	for _, result := range req.Results {
		task, err := i.pingTaskRepo.Find(ctx, result.PingTaskUID)
		if err != nil {
			i.logger.Warn("Failed to find ping task", zap.Error(err), zap.Uint64("uid", result.PingTaskUID), log.AgentUID(req.AgentUID))
		}
		r := convertToPingResult(result, task)
		if err := i.tsdb.InsertPingResult(ctx, r); err != nil {
			i.logger.Warn("Failed to insert ping result", zap.Error(err), log.AgentUID(req.AgentUID))
		}
	}

	return &gopb.Empty{}, nil
}

func convertToPingResult(result *gopb.GrpcPingResult, task *domain.PingTask) domain.PingResult {
	return domain.PingResult{
		RttMicros: result.RttMicros,
		IsTimeout: result.IsTimeout,
		Tag:       task.Tag.Data(),
	}
}

func convertToTcpPingResult(result *gopb.GrpcTcpPingResult, task *domain.TcpPingTask) domain.TcpPingResult {
	return domain.TcpPingResult{
		RttMicros: result.RttMicros,
		IsTimeout: result.IsTimeout,
		Tag:       task.Tag.Data(),
	}
}

func (i *Collector) TcpPingReport(ctx context.Context, req *gopb.GrpcTcpPingReportReq) (*gopb.Empty, error) {
	for _, result := range req.Results {
		task, err := i.tcpPingTaskRepo.Find(ctx, result.TcpPingTaskUID)
		if err != nil {
			i.logger.Warn("Failed to find ping task", zap.Error(err), zap.Uint64("uid", result.TcpPingTaskUID), log.AgentUID(req.AgentUID))
		}
		r := convertToTcpPingResult(result, task)
		if err := i.tsdb.InsertTcpPingResult(ctx, r); err != nil {
			i.logger.Warn("Failed to insert ping result", zap.Error(err), log.AgentUID(req.AgentUID))
		}
	}

	return &gopb.Empty{}, nil
}

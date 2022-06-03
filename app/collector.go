package app

import (
	"context"
	"time"

	"github.com/hhyhhy/tsdb"
	"go.uber.org/zap"

	"pingcc/domain"
	"pingcc/log"
	"pingcc/pb"
)

type PingResult struct {
	RttMicros uint32
	IsTimeout bool
}

type Collector struct {
	pb.UnimplementedCollectorServer

	tsdb     *tsdb.TSDB[PingResult]
	pingRepo domain.PingTargetRepo
	tcpRepo  domain.TcpPingTargetRepo

	logger *zap.Logger
}

func NewCollector(tsdb *tsdb.TSDB[PingResult], pingRepo domain.PingTargetRepo, tcpRepo domain.TcpPingTargetRepo) *Collector {
	logger := log.LWithSvcName("collector")
	return &Collector{tsdb: tsdb, pingRepo: pingRepo, tcpRepo: tcpRepo, logger: logger}
}

func (i *Collector) PingReport(ctx context.Context, req *pb.PingReportReq) (*pb.Empty, error) {
	points := make([]tsdb.Point[PingResult], 0, len(req.Results))
	for _, r := range req.Results {
		t, err := i.pingRepo.Find(ctx, r.ID)
		if err != nil {
			i.logger.Info("Fail to find icmp target", zap.Error(err), zap.Uint64("id", r.ID))
			continue
		}

		tags := append(t.Tags(), tsdb.Tag{Key: "measurement", Value: "icmp"})
		data := PingResult{
			RttMicros: r.RttMicros,
			IsTimeout: r.IsTimeout,
		}
		p := tsdb.NewPoint[PingResult](tags, time.Unix(r.SendAt, 0), data)
		points = append(points, p)
	}

	if err := i.tsdb.WritePoints(points); err != nil {
		i.logger.Info("Fail to write ping point", zap.Error(err))
	}

	return &pb.Empty{}, nil
}

func (i *Collector) TcpPingReport(ctx context.Context, req *pb.TcpPingReportReq) (*pb.Empty, error) {
	points := make([]tsdb.Point[PingResult], 0, len(req.Results))
	for _, r := range req.Results {
		t, err := i.tcpRepo.Find(ctx, r.ID)
		if err != nil {
			i.logger.Info("Fail to find tcp target", zap.Error(err), zap.Uint64("id", r.ID))
			continue
		}

		tags := append(t.Tags(), tsdb.Tag{Key: "measurement", Value: "tcp"})
		data := PingResult{
			RttMicros: r.RttMicros,
			IsTimeout: r.IsTimeout,
		}
		p := tsdb.NewPoint[PingResult](tags, time.Unix(r.SendAt, 0), data)
		points = append(points, p)
	}

	if err := i.tsdb.WritePoints(points); err != nil {
		i.logger.Info("Fail to write tcp ping points", zap.Error(err))
	}

	return &pb.Empty{}, nil
}

func (i *Collector) FpingReport(ctx context.Context, req *pb.FPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

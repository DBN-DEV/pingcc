package collector

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

type impl struct {
	pb.UnimplementedCollectorServer

	tsdb     *tsdb.TSDB[PingResult]
	pingRepo domain.PingTargetRepo
	tcpRepo  domain.TcpPingTargetRepo
}

func New(tsdb *tsdb.TSDB[PingResult], pingRepo domain.PingTargetRepo, tcpRepo domain.TcpPingTargetRepo) *impl {
	return &impl{tsdb: tsdb, pingRepo: pingRepo, tcpRepo: tcpRepo}
}

func (i *impl) PingReport(ctx context.Context, req *pb.PingReportReq) (*pb.Empty, error) {
	points := make([]tsdb.Point[PingResult], 0, len(req.Results))
	for _, r := range req.Results {
		t, err := i.pingRepo.Find(ctx, r.ID)
		if err != nil {
			log.L().Info("Fail to find icmp target", zap.Error(err), zap.Uint64("id", r.ID))
			continue
		}

		tags := t.Tags()
		data := PingResult{
			RttMicros: r.RttMicros,
			IsTimeout: r.IsTimeout,
		}
		p := tsdb.NewPoint[PingResult](tags, time.Unix(r.SendAt, 0), data)
		points = append(points, p)
	}

	if err := i.tsdb.WritePoints(points); err != nil {
		log.L().Info("Fail to write point", zap.Error(err))
	}

	return &pb.Empty{}, nil
}

func (i *impl) TcpPingReport(ctx context.Context, req *pb.TcpPingReportReq) (*pb.Empty, error) {
	points := make([]tsdb.Point[PingResult], 0, len(req.Results))
	for _, r := range req.Results {
		t, err := i.tcpRepo.Find(ctx, r.ID)
		if err != nil {
			log.L().Info("Fail to find tcp target", zap.Error(err), zap.Uint64("id", r.ID))
			continue
		}

		tags := t.Tags()
		data := PingResult{
			RttMicros: r.RttMicros,
			IsTimeout: r.IsTimeout,
		}
		p := tsdb.NewPoint[PingResult](tags, time.Unix(r.SendAt, 0), data)
		points = append(points, p)
	}

	if err := i.tsdb.WritePoints(points); err != nil {
		log.L().Info("Fail to write point", zap.Error(err))
	}

	return &pb.Empty{}, nil
}

func (i *impl) FpingReport(ctx context.Context, req *pb.FPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

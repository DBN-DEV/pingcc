package collector

import (
	"context"
	"strconv"
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
}

func New(tsdb *tsdb.TSDB[PingResult], pingRepo domain.PingTargetRepo) *impl {
	return &impl{tsdb: tsdb, pingRepo: pingRepo}
}

func (i *impl) PingReport(ctx context.Context, req *pb.PingReportReq) (*pb.Empty, error) {
	points := make([]tsdb.Point[PingResult], 0, len(req.Results))
	for _, r := range req.Results {
		t, err := i.pingRepo.Find(ctx, r.ID)
		if err != nil {
			log.L().Info("Fail to find target", zap.Error(err), zap.Uint64("id", r.ID))
			continue
		}

		tag := tsdb.Tag{Key: "agent_id", Value: strconv.Itoa(int(req.AgentID))}
		tags := append(t.Tags(), tag)
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
	return &pb.Empty{}, nil
}

func (i *impl) FpingReport(ctx context.Context, req *pb.FPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

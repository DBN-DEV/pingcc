package collector

import (
	"context"

	"ping-cc/pb"
)

type Impl struct {
	pb.UnimplementedCollectorServer
}

func (i *Impl) HeartbeatReport(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (i *Impl) PingReport(ctx context.Context, req *pb.PingReportReq) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}

func (i *Impl) TcpPingReport(ctx context.Context, req *pb.TcpPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (i *Impl) FpingReport(ctx context.Context, req *pb.FPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

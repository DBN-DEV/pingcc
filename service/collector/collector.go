package collector

import (
	"context"

	"pingcc/pb"
)

type impl struct {
	pb.UnimplementedCollectorServer
}

func New() pb.CollectorServer {
	return &impl{}
}

func (i *impl) PingReport(ctx context.Context, req *pb.PingReportReq) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}

func (i *impl) TcpPingReport(ctx context.Context, req *pb.TcpPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (i *impl) FpingReport(ctx context.Context, req *pb.FPingReportReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

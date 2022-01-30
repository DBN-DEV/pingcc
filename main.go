package main

import (
	"net"
	"ping-cc/service/collector"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"google.golang.org/grpc/encoding"

	"ping-cc/pb"
	"ping-cc/service/controller"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/proto"
)

func run() error {
	encoding.RegisterCodec(vtpb.Codec{})

	collS := grpc.NewServer()
	ctrlS := grpc.NewServer()

	dsn := "user=nom password=dashi dbname=nom host=198.18.5.137 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	ctrl := controller.New(dbConn)
	coll := collector.Impl{}

	pb.RegisterControllerServer(ctrlS, ctrl)
	pb.RegisterCollectorServer(collS, &coll)

	lis, err := net.Listen("tcp", "127.0.0.1:5001")
	if err != nil {
		return err
	}
	lisCtrl, err := net.Listen("tcp", "127.0.0.1:5000")

	go func() {
		if err := collS.Serve(lis); err != nil {
			return
		}
	}()

	if err := ctrlS.Serve(lisCtrl); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/spf13/viper"
	"net"
	"ping-cc/service/collector"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ping-cc/pb"
	"ping-cc/service/controller"
)

func run() error {
	encoding.RegisterCodec(vtpb.Codec{})

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	dsn := viper.GetString("database.dsn")
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	ctrl := controller.New(dbConn)
	coll := collector.Impl{}

	collS := grpc.NewServer()
	ctrlS := grpc.NewServer()
	pb.RegisterControllerServer(ctrlS, ctrl)
	pb.RegisterCollectorServer(collS, &coll)

	addr := viper.GetString("server.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := collS.Serve(lis); err != nil {
			return
		}
	}()

	if err := ctrlS.Serve(lis); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

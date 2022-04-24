package main

import (
	"net"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pingcc/entry"
	"pingcc/pb"
	"pingcc/service/collector"
	"pingcc/service/controller"
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
	dbConn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}
	if err := entry.InitTables(dbConn); err != nil {
		return err
	}

	agentRepo := entry.NewAgentRepo(dbConn)

	ctrl := controller.New(agentRepo)
	coll := collector.New()
	collS := grpc.NewServer()
	ctrlS := grpc.NewServer()
	pb.RegisterControllerServer(ctrlS, ctrl)
	pb.RegisterCollectorServer(collS, coll)

	addr := viper.GetString("server.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		err := collS.Serve(lis)
		errCh <- err
	}()
	go func() {
		err := ctrlS.Serve(lis)
		errCh <- err
	}()
	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

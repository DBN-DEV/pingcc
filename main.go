package main

import (
	"net"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"

	"pingcc/app/collector"
	"pingcc/app/controller"
	"pingcc/domain"
	"pingcc/infra"
	"pingcc/log"
	"pingcc/pb"
)

func run() error {
	encoding.RegisterCodec(vtpb.Codec{})

	// init config
	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	log.L().Info("Read config file success", zap.String("path", "./config.toml"))

	// init db connect
	dsn := viper.GetString("database.dsn")
	if err := infra.InitDB(dsn, log.NewGorm(log.L())); err != nil {
		return err
	}
	log.L().Info("Init db connect success")

	// init db table
	if err := domain.InitTables(infra.DB()); err != nil {
		return err
	}
	log.L().Info("Init db tables success")

	agentRepo := &domain.AgentRepoImpl{DB: infra.DB()}

	gsrv := grpc.NewServer()
	pb.RegisterControllerServer(gsrv, controller.New(agentRepo))
	pb.RegisterCollectorServer(gsrv, collector.New())

	addr := viper.GetString("server.addr")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.L().Info("Create listener success", zap.String("addr", addr))
	log.L().Info("Start app")
	if err := gsrv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

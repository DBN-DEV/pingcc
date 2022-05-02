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
	log.L().Info("Create listener success", zap.String("addr", addr))

	errCh := make(chan error)
	go func() {
		log.L().Info("Start collector app")
		err := collS.Serve(lis)
		errCh <- err
	}()
	go func() {
		log.L().Info("Start controller app")
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

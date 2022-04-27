package main

import (
	"net"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pingcc/entry"
	"pingcc/log"
	"pingcc/pb"
	"pingcc/service/collector"
	"pingcc/service/controller"
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
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: log.NewGorm(log.L())})
	if err != nil {
		return err
	}
	log.L().Info("Init db connect success")

	// init db table
	if err := entry.InitTables(dbConn); err != nil {
		return err
	}
	log.L().Info("Init db tables success")

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
	log.L().Info("Create listener success", zap.String("addr", addr))

	errCh := make(chan error)
	go func() {
		log.L().Info("Start collector service")
		err := collS.Serve(lis)
		errCh <- err
	}()
	go func() {
		log.L().Info("Start controller service")
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

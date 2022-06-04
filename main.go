package main

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/hhyhhy/tsdb"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"

	"pingcc/app"
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

	agentRepo := domain.NewAgentRepo(infra.DB())
	pingRepo := domain.NewPingTargetRepo(infra.DB())
	tcpRepo := domain.NewTcpPingTargetRepo(infra.DB())

	memTSDB := tsdb.New[app.PingResult](2 * time.Minute)

	log.L().Info("Init influxdb client")
	cli := influxdb2.NewClient(viper.GetString("influxdb.url"), viper.GetString("influxdb.token"))
	ok, err := cli.Ping(context.Background())
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("ping influxdb not ok")
	}

	writeAPI := cli.WriteAPI(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))

	log.L().Info("Init Aggregator")
	aggregator := app.NewAggregator(agentRepo, memTSDB, writeAPI)
	go aggregator.AggProc()
	log.L().Info("Init aggregator success")

	gsrv := grpc.NewServer()
	pb.RegisterControllerServer(gsrv, app.NewController(agentRepo))
	pb.RegisterCollectorServer(gsrv, app.NewCollector(memTSDB, pingRepo, tcpRepo))

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

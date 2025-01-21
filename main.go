package main

import (
	"fmt"
	"os"

	vtpb "github.com/planetscale/vtprotobuf/codec/grpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"

	"github.com/DBN-DEV/pingcc/config"
	"github.com/DBN-DEV/pingcc/domain"
	"github.com/DBN-DEV/pingcc/log"
)

var root = &cobra.Command{
	Use:   "pingcc",
	Short: "pingcc means pingmesh controller and collector",
	Long:  `pingcc means pingmesh controller and collector`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Flag("conf").Value.String())
	},
}

func run(path string) error {
	encoding.RegisterCodec(vtpb.Codec{})

	// init config
	if err := config.Init(path); err != nil {
		return fmt.Errorf("init config: %w", err)
	}
	log.L().Info("Read config file success", zap.String("path", path))

	// init db connect
	if err := domain.InitDB(config.C().DB.DSN, log.NewGorm(log.L())); err != nil {
		return err
	}
	log.L().Info("Init db connect success")

	// init grpc server

	return nil
}

func main() {
	root.Flags().StringP("conf", "c", "./config.toml", "The config file path")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

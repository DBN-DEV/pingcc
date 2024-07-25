package domain

import (
	"context"
	"fmt"

	greptime "github.com/GreptimeTeam/greptimedb-ingester-go"

	"github.com/DBN-DEV/pingcc/config"
)

var _cli *GreptimeCli

func GrepTimeCli() *GreptimeCli {
	return _cli
}

func InitGreptime() error {
	cfg := greptime.NewConfig(config.C().TSDB.Greptime.Host).
		WithPort(config.C().TSDB.Greptime.Port).
		WithDatabase(config.C().TSDB.Greptime.Database).
		WithAuth(config.C().TSDB.Greptime.Username, config.C().TSDB.Greptime.Password)

	cli, err := greptime.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("init greptime client: %w", err)
	}

	_cli = &GreptimeCli{cli: cli}

	return nil
}

type GreptimeCli struct {
	cli *greptime.Client
}

func (g *GreptimeCli) InsertPingResult(ctx context.Context, result PingResult) error {
	return nil
}

func (g *GreptimeCli) InsertTcpPingResult(ctx context.Context, result TcpPingResult) error {
	return nil
}

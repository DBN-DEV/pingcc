package app

import (
	"context"

	"github.com/DBN-DEV/pingcc/domain"
)

type PingTaskRepo interface {
	Find(ctx context.Context, uid uint64) (*domain.PingTask, error)
	FindByAgentID(ctx context.Context, agentUID uint64) ([]domain.PingTask, error)
}

type TcpPingTaskRepo interface {
	Find(ctx context.Context, uid uint64) (*domain.TcpPingTask, error)
	FindByAgentID(ctx context.Context, agentUID uint64) ([]domain.TcpPingTask, error)
}

type AgentRepo interface {
	Find(ctx context.Context, uid uint64) (*domain.Agent, error)
	Heartbeat(ctx context.Context, uid uint64) error
}

type TSDB interface {
	InsertPingResult(ctx context.Context, result domain.PingResult) error
	InsertTcpPingResult(ctx context.Context, result domain.TcpPingResult) error
}

package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID uint64 `gorm:"primarykey"`

	AgentUID uint64 `gorm:"index"`
	Name     string

	PingTaskVersion    uint64
	TcpPingTaskVersion uint64

	LastActiveTime time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AgentRepoImpl struct {
	db *gorm.DB
}

func (r *AgentRepoImpl) Find(ctx context.Context, uid uint64) (*Agent, error) {
	var agent Agent
	if err := r.db.Where("agent_uid = ?", uid).WithContext(ctx).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepoImpl) Heartbeat(ctx context.Context, uid uint64) error {
	return r.db.Model(&Agent{}).WithContext(ctx).Where("agent_uid = ?", uid).Update("last_active_time", time.Now()).Error
}

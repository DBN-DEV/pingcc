package domain

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gorm.io/datatypes"
)

type TcpPingTask struct {
	ID  uint64 `gorm:"primarykey"`
	UID uint64 `gorm:"index"`

	AgentID uint64

	Dest string
	Src  string

	TimeoutMS  uint32
	IntervalMS uint32

	Tag datatypes.JSONType[map[string]string]

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type TcpPingTaskImpl struct {
	db *gorm.DB
}

func (r *TcpPingTaskImpl) Find(ctx context.Context, uid uint64) (*TcpPingTask, error) {
	var task TcpPingTask
	if err := r.db.Where("uid = ?", uid).WithContext(ctx).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TcpPingTaskImpl) FindByAgentID(ctx context.Context, agentUID uint64) ([]TcpPingTask, error) {
	var tasks []TcpPingTask
	if err := r.db.Where("agent_id = ?", agentUID).WithContext(ctx).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

package domain

import (
	"context"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PingTask struct {
	ID  uint64 `gorm:"primarykey"`
	UID uint64 `gorm:"index"`

	AgentID uint64

	// The Dest can be a domain name or an IP address
	Dest string
	// The Src can be an interface name or an IP address
	Src string

	TimeoutMS  uint32
	IntervalMS uint32

	Tag datatypes.JSONType[map[string]string]

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type PingTaskRepoImpl struct {
	db *gorm.DB
}

func (r *PingTaskRepoImpl) Find(ctx context.Context, uid uint64) (*PingTask, error) {
	var task PingTask
	if err := r.db.Where("uid = ?", uid).WithContext(ctx).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PingTaskRepoImpl) FindByAgentID(ctx context.Context, agentUID uint64) ([]PingTask, error) {
	var tasks []PingTask
	if err := r.db.Where("agent_id = ?", agentUID).WithContext(ctx).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

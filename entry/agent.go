package entry

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	db *gorm.DB

	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// agent 名称 用于给 agent 提供一个可读的记号
	Name string
	// 控制器侧的 ping 命令版本号
	ControllerPingCommandVersion string
	// agent 侧的 ping 命令版本号 用于和和控制器侧的对比，如果不一致说明需要重新下发命令
	AgentPingCommandVersion string
	// 控制器侧的 tcp ping 命令版本号
	ControllerTcpPingCommandVersion string
	// agent 侧的 tcp ping 命令版本号 用于和和控制器侧的对比，如果不一致说明需要重新下发命令
	AgentTcpPingCommandVersion string
	// 上次活跃时间，每次控制器收到来自 agent 的请求都更新此值
	LastActiveTime time.Time
}

func (a *Agent) PingTargets(ctx context.Context) ([]PingTarget, error) {
	var targets []PingTarget
	if err := a.db.WithContext(ctx).Where("id = ?", a.ID).Find(&targets).Error; err != nil {
		return nil, err
	}

	for i := range targets {
		targets[i].db = a.db
	}

	return targets, nil
}

type AgentRepo interface {
	Find(ctx context.Context, id uint) (*Agent, error)
}

type AgentRepoImpl struct {
	DB *gorm.DB
}

func (i *AgentRepoImpl) Find(ctx context.Context, id uint) (*Agent, error) {
	a := Agent{db: i.DB}
	if err := i.DB.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

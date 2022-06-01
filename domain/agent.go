package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID        uint64 `gorm:"primarykey"`
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

	PingTargets    []PingTarget
	TcpPingTargets []TcpPingTarget
}

// ActivateByGetPingComm 将上次活动时间设置为现在，并设置命令版本号
func (a *Agent) ActivateByGetPingComm(commVer string) {
	a.LastActiveTime = time.Now()
	a.AgentPingCommandVersion = commVer
}

func (a *Agent) ActivateByGetTcpPingComm(commVer string) {
	a.LastActiveTime = time.Now()
	a.AgentTcpPingCommandVersion = commVer
}

type AgentRepo interface {
	FindWithPingTargets(ctx context.Context, id uint64) (*Agent, error)
	FindWithTcpPingTargets(ctx context.Context, id uint64) (*Agent, error)
	Save(ctx context.Context, a *Agent) error
}

func NewAgentRepo(db *gorm.DB) *AgentRepoImpl {
	return &AgentRepoImpl{db: db}
}

type AgentRepoImpl struct {
	db *gorm.DB
}

func (i *AgentRepoImpl) FindWithPingTargets(ctx context.Context, id uint64) (*Agent, error) {
	var a Agent
	if err := i.db.WithContext(ctx).Where("id = ?", id).Preload("PingTargets").First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func (i *AgentRepoImpl) FindWithTcpPingTargets(ctx context.Context, id uint64) (*Agent, error) {
	var a Agent
	if err := i.db.WithContext(ctx).Where("id = ?", id).Preload("TcpPingTargets").First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func (i *AgentRepoImpl) Save(ctx context.Context, a *Agent) error {
	if err := i.db.WithContext(ctx).Save(a).Error; err != nil {
		return err
	}

	return nil
}

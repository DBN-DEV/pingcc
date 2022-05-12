package domain

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"time"

	"gorm.io/gorm"
)

type PingTarget struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 被 ping 地址
	IP string `gorm:"index"`
	// 超时时间 单位毫秒
	TimeoutMS uint32
	// 间隔时间 单位毫秒
	IntervalMS uint32

	// 各种标签
	// 类型 表示内网或外网 等
	Type string
	// 地域 华北华东…………
	Region string
	// 省份
	Province string
	// 运营商
	ISP string

	AgentID uint
	Agent   Agent
}

type PingTargetRepo interface {
	Find(ctx context.Context, id uint)
}

type PingTargetRepoImpl struct {
	DB *gorm.DB

	cache *ristretto.Cache
}

package domain

import (
	"context"
	"github.com/hhyhhy/tsdb"
	"github.com/outcaste-io/ristretto"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type TcpPingTarget struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 被探测的地址 可以是域名，要带端口号
	Address string
	// 超时时间 单位毫秒
	TimeoutMS uint32
	// 间隔时间 单位毫秒
	IntervalMS uint32

	// 地域 华北华东…………
	Region string

	AgentID uint64
	Agent   Agent
}

func (t *TcpPingTarget) Tags() []tsdb.Tag {
	return []tsdb.Tag{
		{Key: "agent_id", Value: strconv.Itoa(int(t.AgentID))},
		{Key: "measurement", Value: "tcp"},
		{Key: "address", Value: t.Address},
		{Key: "region", Value: t.Region},
	}
}

type TcpPingTargetRepo interface {
	Find(ctx context.Context, id uint64) (*TcpPingTarget, error)
}

type TcpPingTargetRepoImpl struct {
	db    *gorm.DB
	cache *ristretto.Cache
}

func NewTcpPingTargetRepo(db *gorm.DB) *TcpPingTargetRepoImpl {
	var maxItem int64 = 10000
	cfg := ristretto.Config{
		NumCounters: maxItem * 10,
		MaxCost:     maxItem,
		BufferItems: 64,
	}
	// cfg is static can not err
	cache, _ := ristretto.NewCache(&cfg)
	return &TcpPingTargetRepoImpl{db: db, cache: cache}
}

func (i *TcpPingTargetRepoImpl) Find(ctx context.Context, id uint64) (*TcpPingTarget, error) {
	if t, ok := i.cache.Get(id); ok {
		return t.(*TcpPingTarget), nil
	}

	var t TcpPingTarget
	if err := i.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		return nil, err
	}

	i.cache.SetWithTTL(id, &t, 1, _defaultTTL)

	return &t, nil
}

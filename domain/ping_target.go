package domain

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hhyhhy/tsdb"
	"github.com/outcaste-io/ristretto"

	"gorm.io/gorm"
)

type PingTarget struct {
	ID        uint64 `gorm:"primarykey"`
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

func (t PingTarget) Tags() []tsdb.Tag {
	return []tsdb.Tag{
		{Key: "agent_id", Value: strconv.Itoa(int(t.AgentID))},
		{Key: "type", Value: t.Type},
		{Key: "region", Value: t.Region},
		{Key: "province", Value: t.Province},
		{Key: "isp", Value: t.ISP},
		{Key: "ip", Value: t.IP},
	}
}

func (t PingTarget) SeriesWithoutIP() string {
	var b strings.Builder
	b.WriteString("agent_id=")
	b.WriteString(strconv.Itoa(int(t.AgentID)))
	b.WriteString(";type=")
	b.WriteString(t.Type)
	b.WriteString(";region=")
	b.WriteString(t.Region)
	b.WriteString(";province=")
	b.WriteString(t.Province)
	b.WriteString(";isp=")
	b.WriteString(t.ISP)

	return b.String()
}

func (t PingTarget) TagsWithoutIP() []tsdb.Tag {
	return []tsdb.Tag{
		{Key: "agent_id", Value: strconv.Itoa(int(t.AgentID))},
		{Key: "type", Value: t.Type},
		{Key: "region", Value: t.Region},
		{Key: "province", Value: t.Province},
		{Key: "isp", Value: t.ISP},
	}
}

type PingTargetRepo interface {
	Find(ctx context.Context, id uint64) (*PingTarget, error)
}

func NewPingTargetRepo(db *gorm.DB) *PingTargetRepoImpl {
	var maxItem int64 = 10000
	cfg := ristretto.Config{
		NumCounters: maxItem * 10,
		MaxCost:     maxItem,
		BufferItems: 64,
	}
	// cfg is static can not err
	cache, _ := ristretto.NewCache(&cfg)
	return &PingTargetRepoImpl{db: db, cache: cache}
}

type PingTargetRepoImpl struct {
	db    *gorm.DB
	cache *ristretto.Cache
}

func (i *PingTargetRepoImpl) Find(ctx context.Context, id uint64) (*PingTarget, error) {
	if t, ok := i.cache.Get(id); ok {
		return t.(*PingTarget), nil
	}

	var t PingTarget
	if err := i.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		return nil, err
	}

	i.cache.SetWithTTL(id, &t, 1, _defaultTTL)

	return &t, nil
}

package domain

import (
	"time"
)

type TcpPingTarget struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 被探测的地址 可以是域名，要带端口号
	Address string `gorm:"index"`
	// 超时时间 单位毫秒
	TimeoutMS uint32
	// 间隔时间 单位毫秒
	IntervalMS uint32

	// 地域 华北华东…………
	Region string

	AgentID uint
	Agent   Agent
}

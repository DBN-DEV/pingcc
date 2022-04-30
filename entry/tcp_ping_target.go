package entry

import (
	"time"

	"gorm.io/gorm"
)

type TcpPingTarget struct {
	db *gorm.DB

	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 被探测的地址 可以是域名，要带端口号
	Address string `gorm:"index:,type:hash"`
	// 超时时间 单位毫秒
	TimeoutMS uint32
	// 间隔时间 单位毫秒
	IntervalMS uint32

	// 地域 华北华东…………
	Region string

	AgentID uint
}

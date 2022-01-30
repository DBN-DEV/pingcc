package controller

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model

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

	PingTargets []PingTarget
}

type PingTarget struct {
	gorm.Model

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

type TcpPingTarget struct {
	gorm.Model

	// 被探测的地址 可以是域名，要带端口号
	Address string `gorm:"index:,type:hash"`
	// 超时时间 单位毫秒
	TimeoutMS uint32
	// 间隔时间 单位毫秒
	IntervalMS uint32

	// 地域 华北华东…………
	Region string

	AgentID uint
	Agent   Agent
}

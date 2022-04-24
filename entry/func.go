package entry

import "gorm.io/gorm"

// InitTables 初始化数据库表
func InitTables(db *gorm.DB) error {
	tables := []any{&Agent{}, &PingTarget{}, &TcpPingTarget{}}
	return db.AutoMigrate(tables...)
}

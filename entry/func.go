package entry

import (
	"gorm.io/gorm"
)

// InitTables Init database.
func InitTables(db *gorm.DB) error {
	tables := []any{&Agent{}, &PingTarget{}, &TcpPingTarget{}}
	return db.AutoMigrate(tables...)
}

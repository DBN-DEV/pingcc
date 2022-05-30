package domain

import (
	"time"

	"gorm.io/gorm"
)

const _defaultTTL = 30 * time.Minute

// InitTables Init database.
func InitTables(db *gorm.DB) error {
	tables := []any{&Agent{}, &PingTarget{}, &TcpPingTarget{}}
	return db.AutoMigrate(tables...)
}

package infra

import (
	"sync/atomic"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db atomic.Value

func InitDB(dsn string, logger logger.Interface) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return err
	}

	_db.Store(db)

	return nil
}

func DB() *gorm.DB {
	return _db.Load().(*gorm.DB)
}

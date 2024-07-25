package domain

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type R struct {
	AgentRepo       *AgentRepoImpl
	PingTaskRepo    *PingTaskRepoImpl
	TcpPingTaskRepo *TcpPingTaskImpl
}

// InitTables Init database.
func InitTables(db *gorm.DB) error {
	tables := []any{&Agent{}, &PingTask{}, &TcpPingTask{}}
	return db.AutoMigrate(tables...)
}

var _r R

func Repo() *R {
	return &_r
}

func InitDB(dsn string, logger logger.Interface) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return err
	}

	r := R{
		AgentRepo:       &AgentRepoImpl{db: db},
		PingTaskRepo:    &PingTaskRepoImpl{db: db},
		TcpPingTaskRepo: &TcpPingTaskImpl{db: db},
	}

	_r = r

	return nil
}

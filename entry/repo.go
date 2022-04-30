package entry

import "gorm.io/gorm"

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		AgentRepo:      &AgentRepoImpl{DB: db},
		PingTargetRepo: &PingTargetRepoImpl{DB: db},
	}
}

type Repo struct {
	AgentRepo      AgentRepo
	PingTargetRepo PingTargetRepo
}

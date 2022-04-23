package entry

import "gorm.io/gorm"

type AgentRepo interface {
	FindPreloadPingCommand(id int) (*Agent, error)
}

func NewAgentRepo(db *gorm.DB) AgentRepo {
	return &agentRepoImpl{db: db}
}

type agentRepoImpl struct {
	db *gorm.DB
}

func (r *agentRepoImpl) FindPreloadPingCommand(id int) (*Agent, error) {
	var a Agent

	if err := r.db.Where("id = ?", id).Preload("PingTargets").First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

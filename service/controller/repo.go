package controller

import "gorm.io/gorm"

type agentRepo interface {
	findPreloadPingCommand(id int) (*Agent, error)
}

type agentRepoImpl struct {
	db *gorm.DB
}

func (r *agentRepoImpl) findPreloadPingCommand(id int) (*Agent, error) {
	var a Agent

	if err := r.db.Where("id = ?", id).Preload("PingTargets").First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

// internal/process/manager.go
package process

import (
	"nexus/internal/config"
	"nexus/pkg/models"
)

type Collector interface {
	Collect([]models.Process) ([]models.Process, error)
}

type Manager struct {
	collectors []Collector
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		collectors: []Collector{
			NewBasicCollector(cfg),
		},
	}
}

func (m *Manager) CollectAll() ([]models.Process, error) {
	var procs []models.Process
	var err error

	for _, c := range m.collectors {
		procs, err = c.Collect(procs)
		if err != nil {
			return procs, err
		}
	}
	return procs, nil
}

package service

import "github.com/gbadopt/gbadopt/internal/model"

// AdoptionStats aggregates application statistics for dashboards.
type AdoptionStats struct {
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"by_status"`
	Approved     int            `json:"approved"`
	PendingTasks int            `json:"pending_tasks"`
}

// ComputeStats derives stats from a list of applications.
func ComputeStats(apps []model.AdoptionApplication) AdoptionStats {
	s := AdoptionStats{Total: len(apps), ByStatus: map[string]int{}}
	for _, a := range apps {
		s.ByStatus[a.Status]++
		if a.Status == "approved" {
			s.Approved++
		}
	}
	return s
}

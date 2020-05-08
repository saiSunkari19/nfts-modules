package common

import "time"

type Info struct {
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Duration `json:"end_date"`
	Territories string        `json:"territories"`
	Exclusivity bool          `json:"exclusivity"`
}

type RightsDetails map[string]Info

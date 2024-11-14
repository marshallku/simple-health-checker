package types

import "time"

type CheckResult struct {
	URL         string    `json:"url"`
	StatusCode  int       `json:"statusCode"`
	TimeTaken   string    `json:"timeTaken"`
	Status      bool      `json:"status"`
	LastChecked time.Time `json:"lastChecked"`
}

type History struct {
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

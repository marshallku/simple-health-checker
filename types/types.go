package types

import "time"

type CheckResult struct {
	URL         string
	StatusCode  int
	TimeTaken   string
	Status      bool
	LastChecked time.Time
}

type History struct {
	URL       string
	Status    string
	Timestamp time.Time
	Message   string
}

package models

type ProcessInfo struct {
	PID       int    `json:"pid"`
	Name      string `json:"name"`
	State     string `json:"state"`
	ParentPID int    `json:"parent_pid"`
}

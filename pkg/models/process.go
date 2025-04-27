package models

type Process struct {
	PID       int    `json:"pid"`
	Name      string `json:"name"`
	State     string `json:"state"`
	ParentPID int    `json:"parent_pid"`
}

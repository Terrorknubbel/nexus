package config

// Config enthält globale Einstellungen für alle Kollektoren.
type Config struct {
	ProcRoot string
}

func NewDefault() *Config {
	return &Config{ProcRoot: "/proc"}
}

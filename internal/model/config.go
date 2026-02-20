package model

// AppConfig holds application configuration.
type AppConfig struct {
	RefreshInterval int    `json:"refresh_interval"`
	TempUnit        string `json:"temp_unit"`
	HistorySize     int    `json:"history_size"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() AppConfig {
	return AppConfig{
		RefreshInterval: 2,
		TempUnit:        "C",
		HistorySize:     60,
	}
}

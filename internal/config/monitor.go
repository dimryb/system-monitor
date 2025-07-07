package config

type (
	MonitorConfig struct {
		Log Log `yaml:"log"`
	}
)

func NewMonitorConfig(configPath string) (*MonitorConfig, error) {
	cfg := &MonitorConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}

package config

type (
	MonitorConfig struct {
		Log  Log  `yaml:"log"`
		GRPC GRPC `yaml:"grpc"`
	}
)

func NewMonitorConfig(configPath string) (*MonitorConfig, error) {
	cfg := &MonitorConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}

package config

type (
	NetworkMetrics struct {
		Enabled bool `yaml:"enabled" env:"NETWORK_METRICS_ENABLED"`

		TopTalkersByProtocol bool `yaml:"topTalkersByProtocol" env:"NETWORK_TOP_TALKERS_BY_PROTOCOL"`
		TopTalkersByTraffic  bool `yaml:"topTalkersByTraffic" env:"NETWORK_TOP_TALKERS_BY_TRAFFIC"`
	}

	CPUMetrics struct {
		Enabled bool `yaml:"enabled" env:"CPU_METRICS_ENABLED"`

		CPUUsagePercent      bool `yaml:"usagePercent" env:"CPU_USAGE_PERCENT"`
		CPUUserModePercent   bool `yaml:"userMode" env:"CPU_USER_MODE"`
		CPUSystemModePercent bool `yaml:"systemMode" env:"CPU_SYSTEM_MODE"`
		CPUIdlePercent       bool `yaml:"idle" env:"CPU_IDLE"`
	}

	DiskMetrics struct {
		Enabled bool `yaml:"enabled" env:"DISK_METRICS_ENABLED"`

		DiskTPS      bool `yaml:"tps" env:"DISK_TPS"`
		DiskKBPerSec bool `yaml:"kbPerSec" env:"DISK_KB_PER_SEC"`
		DiskUsage    bool `yaml:"usage" env:"DISK_USAGE"`
	}

	Metrics struct {
		CPU     CPUMetrics     `yaml:"cpu" env:"METRICS_CPU"`
		Disk    DiskMetrics    `yaml:"disk" env:"METRICS_DISK"`
		Network NetworkMetrics `yaml:"network" env-prefix:"NETWORK_"`
	}

	MonitorConfig struct {
		Log     Log     `yaml:"log"`
		GRPC    GRPC    `yaml:"grpc"`
		Metrics Metrics `yaml:"metrics"`
	}
)

func NewMonitorConfig(configPath string) (*MonitorConfig, error) {
	cfg := &MonitorConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}

package config

type ProfilerConfig struct {
	Server  string `yaml:"server" env:"PROFILER_SERVER" env-default:"http://localhost:4040"`
	Enabled bool   `yaml:"enabled" env:"PROFILER_ENABLED" env-default:"false"`
}

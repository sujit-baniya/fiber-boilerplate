package config

type JwtConfig struct {
	Secret string `yaml:"secret"`
	Expire string `yaml:"expire"`
}

type JwtSecrets struct {
	App JwtConfig `yaml:"app"`
	Api JwtConfig `yaml:"api"`
}

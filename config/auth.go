package config

type AuthConfig struct {
	Type string `yaml:"type" env:"AUTH_TYPE"`
	Driver DatabaseDriver `yaml:"driver"`
}

func (c *AuthConfig) Setup() {


}

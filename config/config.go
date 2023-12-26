package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		MySQL            `yaml:"mysql"`
		Jwt              `yaml:"jwt"`
		Redis            `yaml:"redis"`
		Setting          `yaml:"setting"`
		Cookie           `yaml:"cookie"`
		LoadTestSettings `yaml:"load_test_settings"`
	}

	MySQL struct {
		DBConn string `env-required:"true" yaml:"db_conn"env :"DB_CONN"`
	}

	Jwt struct {
		Secret   string `yaml:"secret" env:"JWT_SECRET"`
		Issuer   string `yaml:"issuer" env:"JWT_ISSUER"`
		Audience string `yaml:"audience" env:"JWT_AUDIENCE"`
	}

	Redis struct {
		CONN string `env-required:"true" yaml:"conn" env:"REDIS_CONN"`
	}

	Setting struct {
		Environment string `yaml:"environment" env:"SETTING_ENVIRONMENT"`
	}

	Cookie struct {
		Key  string `yaml:"key" env:"ADMIN_COOKIE_KEY"`
		Host string `yaml:"host" env:"ADMIN_COOKIE_HOST"`
	}

	LoadTestSettings struct {
		Enable string `yaml:"enable" env:"LOAD_TEST_ENABLE"`
	}
)

func (s Setting) IsDevelopment() bool {
	return s.Environment == "Development"
}

func (s LoadTestSettings) IsMock() bool {
	return s.Enable == "true"
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

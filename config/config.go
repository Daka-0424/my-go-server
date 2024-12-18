package config

import (
	"encoding/hex"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	PlatformUnknown = iota
	PlatformAndroid
	PlatformIOS
	PlatformWebgl
	PlatformWindows
)

func PlatformName(platformNumber uint) string {
	switch platformNumber {
	case PlatformAndroid:
		return "Android"
	case PlatformIOS:
		return "iOS"
	case PlatformWebgl:
		return "WebGL"
	case PlatformWindows:
		return "Windows"
	default:
		return "Unknown"
	}
}

type (
	Config struct {
		MySQL           `yaml:"mysql"`
		Jwt             `yaml:"jwt"`
		Redis           `yaml:"redis"`
		Settings        `yaml:"settings"`
		Cookie          `yaml:"cookie"`
		Kpi             `yaml:"kpi"`
		Appstore        `yaml:"appstore"`
		SandboxAppstore `yaml:"sandbox_appstore"`
		GooglePlay      `yaml:"google_play"`
		MultiDevice     `yaml:"multi_device"`
		RequestKeyIv    `yaml:"request_key_iv"`
		Admin           `yaml:"admin"`
	}

	MySQL struct {
		DBConn string `env-required:"true" yaml:"db_conn" env:"DB_CONN"`
	}

	Jwt struct {
		Secret   string `yaml:"secret" env:"JWT_SECRET"`
		Issuer   string `yaml:"issuer" env:"JWT_ISSUER"`
		Audience string `yaml:"audience" env:"JWT_AUDIENCE"`
	}

	Redis struct {
		CONN string `env-required:"true" yaml:"conn" env:"REDIS_CONN"`
	}

	Settings struct {
		Environment string `yaml:"environment" env:"SETTING_ENVIRONMENT"`
		BaseDomain  string `yaml:"base_domain" env:"SETTING_BASE_DOMAIN"`
	}

	Cookie struct {
		Key  string `yaml:"key" env:"ADMIN_COOKIE_KEY"`
		Host string `yaml:"host" env:"ADMIN_COOKIE_HOST"`
	}

	Kpi struct {
		ProjectID string `yaml:"project_id" env:"KPI_PROJECT_ID"`
	}

	ReviewVersion struct {
		IOS     string `yaml:"ios" env:"REVIEW_VERSION_IOS"`
		Android string `yaml:"android" env:"REVIEW_VERSION_ANDROID"`
	}

	RequirementVersion struct {
		IOS     string `yaml:"ios" env:"REQUIREMENT_VERSION_IOS"`
		Android string `yaml:"android" env:"REQUIREMENT_VERSION_ANDROID"`
	}

	Appstore struct {
		KeyContent string `yaml:"appstore_key_content" env:"APPSTORE_KEY_CONTENT"`
		KeyID      string `yaml:"appstore_private_key" env:"APPSTORE_KEY_ID"`
		BundleID   string `yaml:"appstore_bundle_id" env:"APPSTORE_BUNDLE_ID"`
		IssuerID   string `yaml:"appstore_issuer_id" env:"APPSTORE_ISSUER_ID"`
	}

	SandboxAppstore struct {
		SandboxKeyContent string `yaml:"sandbox_appstore_key_content" env:"SANDBOX_APPSTORE_KEY_CONTENT"`
		SandboxKeyID      string `yaml:"sandbox_appstore_private_key" env:"SANDBOX_APPSTORE_KEY_ID"`
		SandboxBundleID   string `yaml:"sandbox_appstore_bundle_id" env:"SANDBOX_APPSTORE_BUNDLE_ID"`
		SandboxIssuerID   string `yaml:"sandbox_appstore_issuer_id" env:"SANDBOX_APPSTORE_ISSUER_ID"`
	}

	GooglePlay struct {
		Base64EncodedPublicKey       string `yaml:"base64_encoded_public_key" env:"GOOGLE_PLAY_BASE64_ENCODED_PUBLIC_KEY"`
		GoogleApplicationCredentials string `yaml:"google_application_credentials" env:"GOOGLE_STORE_APPLICATION_CREDENTIALS"`
	}

	MultiDevice struct {
		Access string `yaml:"access" env:"multi_device_access"`
	}

	RequestKeyIv struct {
		DefaultKey    []byte
		DefaultIv     []byte
		DefaultKeyStr string `yaml:"default_key_str" env:"REQUEST_DEFAULT_KEY_STR"`
		DefaultIvStr  string `yaml:"default_iv_str" env:"REQUEST_DEFAULT_IV_STR"`
	}

	Admin struct {
		RegisterEmailSender string `yaml:"register_email_sender" env:"REGISTER_EMAIL_SENDER"`
		RegisterEmailPass   string `yaml:"register_email_pass" env:"REGISTER_EMAIL_PASS"`
	}
)

func (s Settings) IsDevelopment() bool {
	return s.Environment == "Development"
}

func (s MultiDevice) IsMultiDeviceAccess() bool {
	return s.Access == "Actiivation"
}

func (c *Config) setDefaultKeyIV() error {
	key, err := hex.DecodeString(c.RequestKeyIv.DefaultKeyStr)
	if err != nil {
		return err
	}
	iv, err := hex.DecodeString(c.RequestKeyIv.DefaultIvStr)
	if err != nil {
		return err
	}

	c.RequestKeyIv.DefaultKey = key
	c.RequestKeyIv.DefaultIv = iv
	return nil
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

	if err := cfg.setDefaultKeyIV(); err != nil {
		return nil, err
	}

	return cfg, nil
}

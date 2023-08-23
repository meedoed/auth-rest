package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute
)

type Config struct {
	Mongo MongoConfig
	HTTP  HTTPConfig
	Auth  AuthConfig
}

type MongoConfig struct {
	URI      string
	User     string
	Password string
	Name     string `mapstructure:"databaseName"`
}

type AuthConfig struct {
	JWT                    JWTConfig
	PasswordSalt           string
	VerificationCodeLength int `mapstructure:"verificationCodeLength"`
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	SigningKey      string
}

type HTTPConfig struct {
	Host               string        `mapstructure:"host"`
	Port               string        `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
}

func Init(configPath string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	return viper.UnmarshalKey("auth", &cfg.Auth.JWT)
}

func setFromEnv(cfg *Config) {

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}

package utils

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	MinUsernameLength int `mapstructure:"MIN_USERNAME_LEN"`
	MaxUsernameLength int `mapstructure:"MAX_USERNAME_LEN"`
	MinPwdLength      int `mapstructure:"MIN_PASSWORD_LEN"`
	MaxPwdLength      int `mapstructure:"MAX_PASSWORD_LEN"`

	MigrationURL string `mapstructure:"MIGRATION_URL"`

	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	AllowedOrigins    []string `mapstructure:"ALLOWED_ORIGINS"`
	GRPCServerAddress string   `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func (c *Config) GetDBSource() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", c.DBDriver, c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}

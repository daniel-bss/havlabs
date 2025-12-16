package utils

import (
	"fmt"

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

	MinioRootUser     string `mapstructure:"MINIO_ROOT_USER"`
	MinioRootPassword string `mapstructure:"MINIO_ROOT_PASSWORD"`
	MinioServerPort   string `mapstructure:"X_MINIO_SERVER_PORT"`
	MinioBrowserPort  string `mapstructure:"X_MINIO_BROWSER_PORT"`
	MinioHost         string `mapstructure:"X_MINIO_HOST"`
	MinioUseSSL       string `mapstructure:"X_MINIO_USESSL"`
	MinioBucketName   string `mapstructure:"X_MINIO_BUCKETNAME"`

	MigrationURL string `mapstructure:"MIGRATION_URL"`

	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func (c *Config) GetDBSource() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", c.DBDriver, c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	DSN string
}

func LoadServerConfig() ServerConfig {
	server_config := ServerConfig{}
	viper.SetConfigFile("internal/config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if host := viper.GetString("server.host"); host != "" {
		server_config.Host = host
	} else {
		server_config.Host = "localhost"
	}

	if port := viper.GetString("server.port"); port != "" {
		server_config.Port = port
	} else {
		server_config.Port = "8000"
	}

	return server_config
}

func LoadDatabaseConfig() DatabaseConfig {
	database_config := DatabaseConfig{}
	viper.SetConfigFile("internal/config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if dsn := viper.GetString("database.dsn"); dsn != "" {
		database_config.DSN = dsn
	} else {
		database_config.DSN = "postgres://dev:123456789@localhost/flotify"
	}

	return database_config
}

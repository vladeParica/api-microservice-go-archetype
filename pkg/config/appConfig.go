package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Configuration AppConfig

type AppConfig struct {
	MicroserviceName    string   `mapstructure:"microserviceName" yaml:"microserviceName"`
	MicroserviceServer  string   `mapstructure:"microserviceServer" yaml:"microserviceServer"`
	MicroservicePort    string   `mapstructure:"microservicePort" yaml:"microservicePort"`
	MicroserviceVersion string   `mapstructure:"microserviceVersion" yaml:"microserviceVersion"`
	Environment         string   `mapstructure:"environment" yaml:"environment"`
	WorkspaceFolder     string   `mapstructure:"workspaceFolder" yaml:"workspaceFolder"`
	Log                 Log      `mapstructure:"log" yaml:"log"`
	Database            Database `mapstructure:"database" yaml:"database"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type Database struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DbName   string `mapstructure:"dbName" yaml:"dbName"`
	SslMode  string `mapstructure:"sslMode" yaml:"sslMode"`
}

func LoadConfigurationMicroservice(path string) {
	fmt.Println("Loading configuration from file [application.yml]")
	viper.SetConfigName("application")
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading dbconfig file, %s", err)
		panic("Error reading dbconfig file")
	}

	// Set undefined variables
	viper.SetDefault("microserviceServer", "localhost")
	viper.SetDefault("microservicePathRoot", "./")

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}

	fmt.Println("Configuration loaded")
}

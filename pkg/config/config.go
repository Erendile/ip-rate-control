package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	ApplicationConfig
}

type ApplicationConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func NewConfiguration() *Config {
	applicationConfig := &ApplicationConfig{}
	applicationConfig.readApplicationConfig()

	return &Config{ApplicationConfig: *applicationConfig}
}

func (c *ApplicationConfig) readApplicationConfig() {
	env, found := os.LookupEnv("ACTIVE_PROFILE")

	if !found {
		env = "local"
	}

	print("ACTIVE_PROFILE: ", env, "\n")

	v := viper.New()
	v.SetTypeByDefaultValue(true)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")

	readConfigErr := v.ReadInConfig()
	if readConfigErr != nil {
		panic("Couldn't load application configuration, cannot start. Terminating. : " + readConfigErr.Error())
	}

	sub := v.Sub(env)

	unMarshallErr := sub.Unmarshal(c)

	if unMarshallErr != nil {
		panic("Configuration cannot deserialize. Terminating. : " + unMarshallErr.Error())
	}
}

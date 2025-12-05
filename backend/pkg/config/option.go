package config

import (
	"github.com/spf13/viper"
	"invoice-scan/backend/pkg"
	"invoice-scan/backend/pkg/log"
	"io"
	"strings"
)

// Option config
type Option func(*Config)

// WithDefaultEnvVars default setup for env variables source
func WithDefaultEnvVars(envPrefix string) Option {
	return func(c *Config) {
		// source: env. variables
		if pkg.IsStringNotEmpty(envPrefix) {
			c.SetEnvPrefix(envPrefix)
		}
		c.AutomaticEnv()
		replacer := strings.NewReplacer(".", "_")
		c.SetEnvKeyReplacer(replacer)
	}
}

// WithDefaultConfigFile default setup for file source
func WithDefaultConfigFile(appName, configFile string) Option {
	return func(c *Config) {
		// set config file by name or path
		configFile = pkg.StringTrimSpace(configFile)
		if pkg.IsStringEmpty(configFile) {
			// load config by name. By default, looking for config.yaml file
			c.SetConfigName("config")
			c.SetConfigType("yaml")
		} else {
			// load config by file path
			c.SetConfigFile(configFile)
		}

		// declare paths to look for the config file in
		appName = pkg.StringTrimSpace(appName)
		if pkg.IsStringNotEmpty(appName) {
			c.AddConfigPath("/etc/" + appName + "/") // /etc/appName/
			c.AddConfigPath("$HOME/." + appName)     // $HOME/.appName
		}
		c.AddConfigPath(".") // working directory

		// Find and read the config file
		if err := c.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Info("Config file not found, app will not load configs from file")
			} else {
				log.Errorw("Config file was found but another error was produced", "error", err)
			}
		}
	}
}

// WithReader load config from io.Reader for testing conveniently by cfgType format
func WithReader(cfgType string, reader io.Reader) Option {
	return func(c *Config) {
		c.SetConfigType(cfgType)
		if err := c.ReadConfig(reader); err != nil {
			log.Error("Cannot read config from reader", "error", err)
		}
	}
}

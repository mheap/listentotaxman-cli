// Package config handles configuration file loading and management.
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Defaults Defaults `mapstructure:"defaults"`
}

// Defaults holds default values for CLI flags
type Defaults struct {
	Region        string `mapstructure:"region"`
	Year          string `mapstructure:"year"`
	Age           string `mapstructure:"age"`
	Pension       string `mapstructure:"pension"`
	StudentLoan   string `mapstructure:"student-loan"`
	TaxCode       string `mapstructure:"tax-code"`
	Extra         int    `mapstructure:"extra"`
	Period        string `mapstructure:"period"`
	Income        int    `mapstructure:"income"`
	Married       bool   `mapstructure:"married"`
	Blind         bool   `mapstructure:"blind"`
	NoNI          bool   `mapstructure:"no-ni"`
	PartnerIncome int    `mapstructure:"partner-income"`
}

// Load loads the configuration file
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "listentotaxman")
	configFile := filepath.Join(configPath, "config.yaml")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("defaults.region", "uk")
	viper.SetDefault("defaults.age", "0")
	viper.SetDefault("defaults.pension", "")
	viper.SetDefault("defaults.student-loan", "")
	viper.SetDefault("defaults.tax-code", "")
	viper.SetDefault("defaults.extra", 0)
	viper.SetDefault("defaults.year", "")
	viper.SetDefault("defaults.period", "yearly")
	viper.SetDefault("defaults.income", 0)
	viper.SetDefault("defaults.married", false)
	viper.SetDefault("defaults.blind", false)
	viper.SetDefault("defaults.no-ni", false)
	viper.SetDefault("defaults.partner-income", 0)

	// Read config file if it exists
	if _, err := os.Stat(configFile); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetString returns a string config value with fallback
func GetString(key, fallback string) string {
	val := viper.GetString(key)
	if val == "" {
		return fallback
	}
	return val
}

// GetInt returns an int config value with fallback
func GetInt(key string, fallback int) int {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetInt(key)
}

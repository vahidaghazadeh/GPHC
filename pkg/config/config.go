package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Commit analysis settings
	MaxCommitsToAnalyze int `mapstructure:"max_commits_to_analyze"`
	
	// Branch analysis settings
	StaleBranchThresholdDays int `mapstructure:"stale_branch_threshold_days"`
	
	// Commit message settings
	MaxCommitMessageLength int `mapstructure:"max_commit_message_length"`
	
	// Commit size settings
	MaxCommitSizeLines int `mapstructure:"max_commit_size_lines"`
	
	// Scoring weights
	Weights Weights `mapstructure:"weights"`
}

// Weights holds the scoring weights for different categories
type Weights struct {
	Documentation int `mapstructure:"documentation"`
	Commits       int `mapstructure:"commits"`
	Hygiene       int `mapstructure:"hygiene"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxCommitsToAnalyze:     50,
		StaleBranchThresholdDays: 60,
		MaxCommitMessageLength:   72,
		MaxCommitSizeLines:       500,
		Weights: Weights{
			Documentation: 3,
			Commits:       4,
			Hygiene:       2,
		},
	}
}

// LoadConfig loads configuration from file and environment
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("gphc")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.gphc")
		viper.AddConfigPath("/etc/gphc")
	}

	// Set default values
	viper.SetDefault("max_commits_to_analyze", 50)
	viper.SetDefault("stale_branch_threshold_days", 60)
	viper.SetDefault("max_commit_message_length", 72)
	viper.SetDefault("max_commit_size_lines", 500)
	viper.SetDefault("weights.documentation", 3)
	viper.SetDefault("weights.commits", 4)
	viper.SetDefault("weights.hygiene", 2)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, use defaults
	}

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetStaleThreshold returns the stale branch threshold as a duration
func (c *Config) GetStaleThreshold() time.Duration {
	return time.Duration(c.StaleBranchThresholdDays) * 24 * time.Hour
}

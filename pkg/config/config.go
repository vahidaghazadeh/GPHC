package config

import (
	"os"
	"path/filepath"
	"strings"
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
	Structure     int `mapstructure:"structure"`
	Security      int `mapstructure:"security"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxCommitsToAnalyze:      50,
		StaleBranchThresholdDays: 60,
		MaxCommitMessageLength:   72,
		MaxCommitSizeLines:       500,
		Weights: Weights{
			Documentation: 3,
			Commits:       4,
			Hygiene:       2,
			Structure:     2,
			Security:      5,
		},
	}
}

// LoadConfig loads configuration from file and environment
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("gphc")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".gphc"))
		}
		v.AddConfigPath("/etc/gphc")
	}

	v.SetEnvPrefix("GPHC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetDefault("max_commits_to_analyze", 50)
	v.SetDefault("stale_branch_threshold_days", 60)
	v.SetDefault("max_commit_message_length", 72)
	v.SetDefault("max_commit_size_lines", 500)
	v.SetDefault("weights.documentation", 3)
	v.SetDefault("weights.commits", 4)
	v.SetDefault("weights.hygiene", 2)
	v.SetDefault("weights.structure", 2)
	v.SetDefault("weights.security", 5)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, use defaults
	}

	// Unmarshal config
	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetStaleThreshold returns the stale branch threshold as a duration
func (c *Config) GetStaleThreshold() time.Duration {
	return time.Duration(c.StaleBranchThresholdDays) * 24 * time.Hour
}

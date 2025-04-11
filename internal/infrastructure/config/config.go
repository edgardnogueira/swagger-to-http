package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ConfigProvider implements the ConfigProvider interface
type ConfigProvider struct {
	viper *viper.Viper
}

// NewConfigProvider creates a new ConfigProvider
func NewConfigProvider() *ConfigProvider {
	v := viper.New()
	
	// Set default configuration values
	setDefaults(v)
	
	// Configure viper to read from environment variables
	v.SetEnvPrefix("STH") // STH = Swagger To Http
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Try to find and read config file
	v.SetConfigName("swagger-to-http")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.swagger-to-http")
	v.AddConfigPath("/etc/swagger-to-http")
	
	// Silently ignore if config file is not found
	_ = v.ReadInConfig()
	
	return &ConfigProvider{
		viper: v,
	}
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("output.directory", "http-requests")
	v.SetDefault("generator.indent_json", true)
	v.SetDefault("generator.include_auth", false)
	v.SetDefault("generator.auth_header", "Authorization")
	v.SetDefault("generator.default_tag", "default")
	v.SetDefault("snapshots.directory", "snapshots")
	v.SetDefault("snapshots.update_on_difference", false)
}

// GetString retrieves a string configuration value
func (c *ConfigProvider) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetInt retrieves an integer configuration value
func (c *ConfigProvider) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func (c *ConfigProvider) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetStringMap retrieves a string map configuration value
func (c *ConfigProvider) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

// GetStringSlice retrieves a string slice configuration value
func (c *ConfigProvider) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// GetConfigFilePath returns the path to the used config file
func (c *ConfigProvider) GetConfigFilePath() string {
	return c.viper.ConfigFileUsed()
}

// SaveConfig saves the current configuration to a file
func (c *ConfigProvider) SaveConfig(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	
	return c.viper.WriteConfigAs(filePath)
}

// Set sets a configuration value
func (c *ConfigProvider) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

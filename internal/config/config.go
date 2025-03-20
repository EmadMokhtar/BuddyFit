package config

import (
	"fmt"
	"os"
	"strings"
)

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type AIProviderConfig struct {
	Name string
	Host string
	Port string
	Key  string
}

func NewOpenAIConfig(key string) *AIProviderConfig {
	return &AIProviderConfig{
		Name: "openai",
		Key:  key,
	}
}

func NewOllamaConfig(host, port string) *AIProviderConfig {
	return &AIProviderConfig{
		Name: "ollama",
		Host: host,
		Port: port,
	}

}

func (c *AIProviderConfig) GetOllamaAPIURL() string {
	if strings.Contains(c.Host, "http://") {
		return fmt.Sprintf("%s/api/chat", c.Host)
	}
	return fmt.Sprintf("http://%s:%s/api/chat", c.Host, c.Port)
}

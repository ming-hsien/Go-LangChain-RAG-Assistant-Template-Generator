package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GitHubToken    string `mapstructure:"GITHUB_TOKEN"`
	LLMModel       string `mapstructure:"LLM_MODEL"`
	LLMBaseURL     string `mapstructure:"LLM_BASE_URL"`
	QdrantURL      string `mapstructure:"QDRANT_URL"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	AppName        string `mapstructure:"APP_NAME"`
	CollectionName string `mapstructure:"COLLECTION_NAME"`
	SystemPrompt   string `mapstructure:"SYSTEM_PROMPT"`
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}

	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("GITHUB_TOKEN", "")
	viper.SetDefault("LLM_MODEL", "gpt-4o")
	viper.SetDefault("LLM_BASE_URL", "https://models.inference.ai.azure.com/v1")
	viper.SetDefault("QDRANT_URL", "http://localhost:6333")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("APP_NAME", "AI Assistant")
	viper.SetDefault("COLLECTION_NAME", "documents")
	viper.SetDefault("SYSTEM_PROMPT", "You are a helpful AI Assistant. Answer the question based ONLY on the provided context")

	AppConfig = &Config{
		GitHubToken:    viper.GetString("GITHUB_TOKEN"),
		LLMModel:       viper.GetString("LLM_MODEL"),
		LLMBaseURL:     viper.GetString("LLM_BASE_URL"),
		QdrantURL:      viper.GetString("QDRANT_URL"),
		ServerPort:     viper.GetString("SERVER_PORT"),
		AppName:        viper.GetString("APP_NAME"),
		CollectionName: viper.GetString("COLLECTION_NAME"),
		SystemPrompt:   viper.GetString("SYSTEM_PROMPT"),
	}

	if AppConfig.GitHubToken == "" {
		log.Println("WARNING: GITHUB_TOKEN is not set")
	}
}

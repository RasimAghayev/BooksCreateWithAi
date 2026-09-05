package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type ChannelConfig struct {
	Name   string
	Token  string
	ChatID string
	Rules  map[string]interface{}
}

type Config struct {
	AppPort      string
	DatabaseURL  string
	UploadDir    string
	ChannelsDir  string
	APITimeout   int
	LogLevel     string
	Channels     map[string]ChannelConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	exeDir, err := getExecutableDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable directory: %w", err)
	}

	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "./uploads.db"),
		UploadDir:   resolvePath(exeDir, getEnv("UPLOAD_DIR", "../../telegram_upload")),
		ChannelsDir: resolvePath(exeDir, getEnv("CHANNELS_DIR", "../channels")),
		APITimeout:  60,
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Channels:    make(map[string]ChannelConfig),
	}

	if err := cfg.loadChannels(); err != nil {
		return nil, fmt.Errorf("failed to load channels: %w", err)
	}

	return cfg, nil
}

func (c *Config) loadChannels() error {
	entries, err := os.ReadDir(c.ChannelsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		channelName := entry.Name()
		channelDir := filepath.Join(c.ChannelsDir, channelName)

		configPath := filepath.Join(channelDir, "config.json")
		envPath := filepath.Join(channelDir, ".env")

		_, configErr := os.Stat(configPath)
		_, envErr := os.Stat(envPath)

		if os.IsNotExist(configErr) || os.IsNotExist(envErr) {
			continue
		}

		envMap, err := godotenv.Read(envPath)
		if err != nil {
			continue
		}

		token := envMap["TELEGRAM_BOT_TOKEN"]
		chatID := envMap["TELEGRAM_CHAT_ID"]

		if token == "" || chatID == "" {
			continue
		}

		rules, err := loadRules(configPath)
		if err != nil {
			fmt.Printf("Warning: failed to load rules for %s: %v\n", channelName, err)
			rules = make(map[string]interface{})
		}

		c.Channels[channelName] = ChannelConfig{
			Name:   channelName,
			Token:  token,
			ChatID: chatID,
			Rules:  rules,
		}
	}

	return nil
}

func loadRules(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	rules, ok := raw["rules"].(map[string]interface{})
	if !ok {
		return make(map[string]interface{}), nil
	}

	return rules, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func resolvePath(baseDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(baseDir, path))
}

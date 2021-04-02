package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"os"
)

// This struct should match the config.json, used for parsing
//
// Session contains an AuthKey and EncryptionKey
// AuthKey should be 32 or 64 bytes
// EncryptionKey must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
type AppConfig struct {
	Session struct {
		AuthKey       []byte //note will be serialized as base64
		EncryptionKey []byte
	}
	Database struct {
		Type string // one of mysql, sqlite
		DSN  string
	}
}

func GetConfig() (AppConfig, error) {
	configPath := os.Getenv("APP_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		cwd, _ := os.Getwd()
		return AppConfig{}, fmt.Errorf("unable to open config file at %s note cwd %s - %w", configPath, cwd, err)
	}
	cfg := AppConfig{}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return AppConfig{}, fmt.Errorf("failed to unmarshal %s - %w", configPath, err)
	}
	//@todo validate

	return cfg, nil
}

func GetExampleConfigJson() (string, error) {
	cfg := AppConfig{
		Session: struct {
			AuthKey       []byte
			EncryptionKey []byte
		}{
			AuthKey:       securecookie.GenerateRandomKey(64),
			EncryptionKey: securecookie.GenerateRandomKey(32),
		},
	}
	indent, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return "", fmt.Errorf("can't serialize Config object - %w", err)
	}
	return string(indent), nil
}

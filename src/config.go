package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	CurrentAddress string
	GeoCodeSucceeded bool
}

func ConfPath() string {
	home, _ :=  os.UserHomeDir()
	configDir := filepath.Join(home, "/go-weather")
	configPath := filepath.Join(configDir, "config.json")
	return configPath
}

func MakeConfIfDoesntExist() {
	home, _ :=  os.UserHomeDir()
	configDir := filepath.Join(home, "/go-weather")
	if _, err := os.Stat(ConfPath()); os.IsNotExist(err) {
		os.Mkdir(configDir, os.ModePerm)
		os.WriteFile(ConfPath(), []byte("{\n\"currentAddress\": \"\"\n}"), os.ModePerm)
	}
}

func LoadConfig() Config {
	MakeConfIfDoesntExist()
	data, _ := os.ReadFile(ConfPath())

	var config Config
	err := json.Unmarshal(data, &config)
	if (err != nil) {
		// Config in invalid format
		fmt.Println("Warn: could not parse config, might be in invalid format")
		return *new(Config)
	}

	return config
}

func SaveConfig(newConfig Config) {
	MakeConfIfDoesntExist()
	configData, _ := json.Marshal(newConfig)
	os.WriteFile(ConfPath(), configData, os.ModePerm)
}
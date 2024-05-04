package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultPath string   `json:"DefaultPath"`
	URLStore    []string `json:"URLStore"`
}

func EnsureConfigExisit() {
	configDir := GetConfigDir()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create config dir: %s", err))
		}
	}

	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := Config{
			DefaultPath: "",
			URLStore: []string{
				"https://deviantart.com",
			},
		}

		data, err := json.MarshalIndent(defaultConfig, "", " ")

		if err != nil {
			panic("Failed to searialize default config")
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			panic("Failed to write default config file")
		}
	}
}

func GetConfigDir() string {
	configDir := os.Getenv("LOCALAPPDATA")

	return filepath.Join(configDir, "cbg")
}

func getConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

func ReadCofnig() Config {
	EnsureConfigExisit()

	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)

	if err != nil {
		return Config{
			DefaultPath: "",
			URLStore: []string{
				"https://deviantart.com",
			},
		}
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		panic("Error rading config")
	}

	return config
}

func WriteConfig(config Config) {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(config, "", " ")

	if err != nil {
		panic("Failed to serialize config")
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		panic("Failed to wrtie config file")
	}
}

func setDefaultPath(value string) {
	config := ReadCofnig()
	config.DefaultPath = value

	WriteConfig(config)

	fmt.Println("Config updated!")
}

func AddUrlToStore(value string) {
	config := ReadCofnig()

	for _, url := range config.URLStore {
		if url == value {
			fmt.Printf("URL already exists: %s\n", value)
			return
		}
	}

	config.URLStore = append(config.URLStore, value)

	WriteConfig(config)

	fmt.Printf("Added: %s to store", value)
}

func RemoveUrlFromStore(value string) {
	config := ReadCofnig()
	indexToRemove := -1

	for i, url := range config.URLStore {
		if url == value {
			indexToRemove = i
			break
		}
	}

	if indexToRemove != -1 {
		config.URLStore = append(config.URLStore[:indexToRemove], config.URLStore[indexToRemove+1:]...)
		WriteConfig(config)
		fmt.Printf("Removed %s from store", value)
	} else {
		fmt.Println("URL not found in store.")
	}
}

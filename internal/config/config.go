package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DB_URL string `json:"db_url"`
	USER   string `json:"user"`
}

const configFileName = ".gatorconfig.json"
const db_url = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not determine the home directory: %q", err)
	}
	completePath := filepath.Join(homeDir, configFileName)
	return completePath, nil
}

func Read() (Config, error) {
	configpath, _ := getConfigPath()
	file, err := os.Open(configpath)
	if err != nil {
		return Config{}, fmt.Errorf("Could not open file: %q", err)
	}
	defer file.Close()
	c := Config{}

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding file: %q", err)
	}
	return c, nil
}

func (c *Config) SetUser(userName string) error {
	c.USER = userName
	c.DB_URL = db_url
	return write(c)
}

func write(c *Config) error {
	fullpath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("Error reading file path: %q", err)
	}
	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("Error creating file: %q", err)
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(&c)
	if err != nil {
		return fmt.Errorf("Error encoding file: %q", err)
	}
	return nil

}

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	fileutil "github.com/zan8in/pins/file"
	"gopkg.in/yaml.v3"
)

var (
	defaultName   = "zoomeye"
	defaultConfig = defaultName + ".yaml"
)

type Config struct {
	Zoomeye Zoomeye `yaml:"zoomeye"`
}

type Zoomeye struct {
	ApiKey []string `yaml:"api-key"`
}

var (
	ValidKeys   []string
	InvalidKeys []string
)

func NewConfig() (*Config, error) {
	if configFile, err := getConfig(); err != nil {
		return nil, err
	} else {
		config, err := readConfig(configFile)
		if err != nil {
			return nil, err
		}
		if len(GetApiKey()) == 0 {
			return nil, fmt.Errorf("no api key found")
		}
		return config, nil
	}
}

func getConfig() (configFile string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return configFile, err
	}

	configDir := filepath.Join(homeDir, ".config", defaultName)
	if !fileutil.FolderExists(configDir) {
		if err = os.MkdirAll(configDir, 0755); err != nil {
			return configFile, err
		}
	}

	configFile = filepath.Join(configDir, defaultConfig)
	if !fileutil.FileExists(configFile) {
		config := &Config{}
		config.Zoomeye.ApiKey = []string{DefaultKey}
		if err = createConfig(configFile, config); err != nil {
			return configFile, err
		}
	}

	return configFile, nil
}

func createConfig(configFile string, config *Config) error {
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(configYaml); err != nil {
		return err
	}

	return nil
}

func readConfig(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	if len(config.Zoomeye.ApiKey) == 0 {
		return nil, errors.New("api key not found")
	}

	ValidKeys = append(ValidKeys, config.Zoomeye.ApiKey...)

	return config, nil
}

func GetApiKey() string {
	for _, v := range ValidKeys {
		if v == DefaultKey {
			continue
		}
		if !IsApiKeyNull(v) {
			return v
		}
	}
	return ""
}

func SetApiKeyNull(apikey string) {
	if !IsApiKeyNull(apikey) {
		InvalidKeys = append(InvalidKeys, apikey)
	}
}

func IsApiKeyNull(apikey string) bool {
	for _, v := range InvalidKeys {
		if v == apikey {
			return true
		}
	}
	return false
}

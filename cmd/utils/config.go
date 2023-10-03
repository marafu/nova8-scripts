package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Checkmarx struct {
		AccessToken string `yaml:"access_token"`
		BaseUrl     string `yaml:"base_url"`
		AuthUrl     string `yaml:"auth_url"`
		ClientID    string `yaml:"client_id"`
		TenantName  string `yaml:"tenant_name"`
	} `yaml:"checkmarx"`
	General struct {
		Proxy string `yaml:"proxy"`
	} `yaml:"general"`
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

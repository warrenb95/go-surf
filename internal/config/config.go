package config

import (
	"context"
	"fmt"
	"github/warrenb95/go-surf/internal/gosurf"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServingPort         uint          `yaml:"servingPort"`
	Spots               []gosurf.Spot `yaml:"spots"`
	AccountSID          string        `yaml:"accountSid"`
	MessagingServiceSid string        `yaml:"messagingServiceSid"`
	TargetMobileNumber  string        `yaml:"targetMobileNumber"`
	StormglassURL       string        `yaml:"stormglassURL"`
}

func Parse(ctx context.Context, filename string) (Config, error) {
	yamlfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read file %v, error: %v", filename, err)
	}

	var cfg Config
	err = yaml.Unmarshal(yamlfile, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal file %v, error: %v", filename, err)
	}

	return cfg, nil
}

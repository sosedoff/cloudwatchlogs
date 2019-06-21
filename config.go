package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Host         string `long:"host"`
	Port         string `long:"port"`
	AuthUser     string `long:"auth-user"`
	AuthPassword string `long:"auth-pass"`
	AccessKey    string `long:"access-key"`
	SecretKey    string `long:"secret-key"`
	Region       string `long:"region"`
	Profile      string `long:"profile"`
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func readConfig() (*Config, error) {
	args := os.Args
	config := &Config{}

	if _, err := flags.ParseArgs(config, args); err != nil {
		return nil, err
	}

	if config.Host == "" {
		config.Host = "0.0.0.0"
	}
	if config.Port == "" {
		config.Port = "5555"
	}

	if config.Region == "" {
		config.Region = os.Getenv("AWS_REGION")
	}
	if config.Region == "" {
		config.Region = "us-east-1"
	}
	if config.AccessKey == "" {
		config.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	}
	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("AWS_SECRET_KEY")
	}
	if config.Profile == "" {
		config.Profile = os.Getenv("AWS_PROFILE")
	}

	return config, nil
}

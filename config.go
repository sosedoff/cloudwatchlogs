package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Host         string `long:"host" default:"0.0.0.0" description:"Server host"`
	Port         string `long:"port" default:"5555" description:"Server port"`
	AuthUser     string `long:"auth-user" description:"User name for basic authentication"`
	AuthPassword string `long:"auth-pass" description:"User password for basic authentication"`
	AccessKey    string `long:"access-key" description:"AWS access key"`
	SecretKey    string `long:"secret-key" description:"AWS secret key"`
	Region       string `long:"region" description:"AWS region"`
	Profile      string `long:"profile" description:"AWS CLI profile"`
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

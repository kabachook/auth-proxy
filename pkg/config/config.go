/*
Package config implements config structure and loading
*/
package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// AuthnConfig : Authentication config
type AuthnConfig struct {
	JWT struct {
		Secret string `yaml: "secret"`
	} `yaml: "jwt"`
}

// AuthzConfig : Authorization config
type AuthzConfig struct {
	LDAP struct {
		URL      string `yaml: "url"`
		Username string `yaml: "username"`
		Password string `yaml: "password"`
	} `yaml: "ldap"`
}

// Backend : backend struct
type Backend struct {
	Name     string `yaml: "name"`
	Endpoint string `yaml: "endpoint"`
}

// Route : proxy route
type Route struct {
	Match struct {
		Prefix string `yaml: "prefix"`
	} `yaml: "match"`
	Backend string `yaml: "backend"`
}

// Config : main config struct
type Config struct {
	Listen   string      `yaml: "listen"`
	Authn    AuthnConfig `yaml: "authn"`
	Authz    AuthzConfig `yaml: "authz"`
	Backends []Backend   `yaml: "backends"`
	Routes   []Route     `yaml: "routes"`
}

// Load : loads config from file
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

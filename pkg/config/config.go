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
		Secret string `yaml:"secret"`
		Field  string `yaml:"field"`
	} `yaml:"jwt"`
}

type ACLEntry struct {
	Action string   `yaml:"action"`
	Users  []string `yaml:"users"`
	Routes []string `yaml:"routes"`
}

type LDAP struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	BaseDN   string `yaml:"base_dn"`
	Filter   string `yaml:"filter"`
}

// AuthzConfig : Authorization config
type AuthzConfig struct {
	LDAP LDAP `yaml:"ldap"`
	ACL []ACLEntry `yaml:"acl"`
}

// Backend : backend struct
type Backend struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Scheme string `yaml:"scheme"`
}

// Route : proxy route
type Route struct {
	Match struct {
		Host string `yaml:"host"`
	} `yaml:"match"`
	Backend string `yaml:"backend"`
}

// Config : main config struct
type Config struct {
	Listen   string      `yaml:"listen"`
	Authn    AuthnConfig `yaml:"authn"`
	Authz    AuthzConfig `yaml:"authz"`
	Backends []Backend   `yaml:"backends"`
	Routes   []Route     `yaml:"routes"`
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

// BackendsToMap coverts []Backend to map[Backend.Name]Backend
func BackendsToMap(backends []Backend) map[string]Backend {
	m := make(map[string]Backend)
	for _, b := range backends {
		m[b.Name] = b
	}
	return m
}

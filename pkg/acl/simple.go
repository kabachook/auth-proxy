package acl

import "github.com/kabachook/auth-proxy/pkg/config"

// SimpleACL: default ACL implementation
// Supports only `allow` action, by default is a blocklist
type SimpleACL struct {
	entries []config.ACLEntry
}

func NewSimpleACL(cfg []config.ACLEntry) *SimpleACL {
	return &SimpleACL{entries: cfg}
}

func (a *SimpleACL) Check(username string, route string) bool {
	// TODO: make it faster (e.g. hashmap user[rules])
	for _, entry := range a.entries {
		if entry.Action == "allow" && contains(username, entry.Users) { // TODO: optimize search
			if contains(route, entry.Routes) {
				return true
			}
		}
	}
	return false
}

func contains(key string, arr []string) bool {
	for _, value := range arr {
		if key == value {
			return true
		}
	}
	return false
}

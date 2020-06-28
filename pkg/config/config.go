package config

type AuthnConfig struct {
	JWT struct {
		Secret string `yaml: "secret"`
	} `yaml: "jwt"`
}

type AuthzConfig struct {
	LDAP struct {
		URL string `yaml: "url"`
		Username string `yaml: "username"`
		Password string `yaml: "password"`
	} `yaml: "ldap"`
}

type Backend struct {
	Name string `yaml: "name"`
	Endpoint string `yaml: "endpoint"`
}

type Route struct {
	Match struct {
		Prefix string `yaml: "prefix"`
	} `yaml: "match"`
	Bakcend string `yaml: "backend"`
}

type Config struct {
	Listen string `yaml: "listen"`
	Authn AuthnConfig `yaml: "authn"`
	Authz AuthzConfig `yaml: "authz"`
	Backends []Backend `yaml: "backends"`
	Routes []Route `yaml: "routes"`
}
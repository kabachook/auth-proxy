package authz

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/kabachook/auth-proxy/pkg/config"
)

// LDAPAuthz is a basic LDAP Authorizer
// Authorizes user if it exists in LDAP
type LDAPAuthz struct {
	conn ldap.Conn
	cfg  config.AuthzConfig
}

func NewLDAPAuthz(cfg config.AuthzConfig) (*LDAPAuthz, error) {
	// Not sure if connection should happen in the constructor
	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		return nil, err
	}

	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: cfg.LDAP.Username,
		Password: cfg.LDAP.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LDAPAuthz{
		conn: *conn,
		cfg:  cfg,
	}, nil
}

func (a *LDAPAuthz) Authorize(username string) (bool, error) {
	req := ldap.NewSearchRequest(
		a.cfg.LDAP.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(uid=%s)", username), // TODO: probably unhardcode filter
		[]string{"dn", "cn"},
		nil,
	)

	sr, err := a.conn.Search(req)
	if err != nil && err.(*ldap.Error).ResultCode != ldap.LDAPResultNoSuchObject {
		switch e := err.(*ldap.Error).ResultCode; {
		case e == ldap.LDAPResultNoSuchObject:
			return false, nil
		default:
			return false, err
		}
	} else if err.(*ldap.Error).ResultCode == ldap.LDAPResultNoSuchObject {
		return false, nil
	}

	if len(sr.Entries) < 1 {
		return false, nil
	}

	return true, nil
}

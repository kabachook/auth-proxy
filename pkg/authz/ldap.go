package authz

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/kabachook/auth-proxy/pkg/config"
)

// LDAPAuthz is a basic LDAP Authorizer
// Authorizes user if LDAP search query with filter returns non-zero result
type LDAPAuthz struct {
	conn *ldap.Conn
	cfg  config.LDAP
}

func NewLDAPAuthz(cfg config.LDAP) (*LDAPAuthz, error) {
	return &LDAPAuthz{
		conn: nil,
		cfg:  cfg,
	}, nil
}

func (a *LDAPAuthz) Open() error {
	conn, err := ldap.DialURL(a.cfg.URL)
	if err != nil {
		return err
	}

	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: a.cfg.Username,
		Password: a.cfg.Password,
	})
	if err != nil {
		return err
	}
	a.conn = conn
	return nil
}

func (a *LDAPAuthz) Close() {
	if a.conn != nil {
		a.conn.Close()
	}
}

func (a *LDAPAuthz) Authorize(username string) (bool, error) {
	if a.conn == nil {
		return false, fmt.Errorf("connection is not opened")
	}

	req := ldap.NewSearchRequest(
		a.cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf(a.cfg.Filter, username),
		[]string{"dn", "cn"},
		nil,
	)

	sr, err := a.conn.Search(req)
	if err != nil {
		switch e := err.(*ldap.Error).ResultCode; {
		case e == ldap.LDAPResultNoSuchObject:
			return false, nil
		default:
			return false, err
		}
	}

	if len(sr.Entries) < 1 {
		return false, nil
	}

	return true, nil
}

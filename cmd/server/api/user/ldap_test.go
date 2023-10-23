package user

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/zhenorzz/goploy/config"
	"testing"
)

func TestLdap(t *testing.T) {

	ldapConfig := config.LDAPConfig{
		Enabled:    true,
		URL:        "ldap://ldap.forumsys.com",
		BindDN:     "cn=read-only-admin,dc=example,dc=com",
		Password:   "password",
		BaseDN:     "dc=example,dc=com",
		UID:        "uid",
		Name:       "cn",
		UserFilter: "",
	}

	searchAccount := "riemann"
	searchPwd := "password"

	conn, err := ldap.DialURL(ldapConfig.URL)

	if err != nil {
		t.Fatal(err)
	}

	if ldapConfig.BindDN != "" {
		if err := conn.Bind(ldapConfig.BindDN, ldapConfig.Password); err != nil {
			t.Fatal(err)
		}
	}

	filter := fmt.Sprintf("(%s=%s)", ldapConfig.UID, searchAccount)
	if ldapConfig.UserFilter != "" {
		filter = fmt.Sprintf("(&(%s)%s)", ldapConfig.UserFilter, filter)
	}

	searchRequest := ldap.NewSearchRequest(
		ldapConfig.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{ldapConfig.UID, ldapConfig.Name},
		nil)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		t.Fatal(err)
	}

	if len(sr.Entries) != 1 {
		t.Fatal(err)
	}
	if err := conn.Bind(sr.Entries[0].DN, searchPwd); err != nil {
		t.Fatal(err)
	}

	for _, attr := range sr.Entries[0].Attributes {
		fmt.Printf("%s: %v\n", attr.Name, attr.Values)
	}
}

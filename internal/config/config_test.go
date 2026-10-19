package config

import (
	"net/url"
	"testing"
)

func TestDSN_EscapesSpecialCharsInPassword(t *testing.T) {
	const password = "p@ss/w:rd#1?x"

	c := DBConfig{
		Host: "localhost", Port: 5432, User: "postgres",
		Name: "warehouse_db", SSLMode: "disable", Password: password,
	}

	u, err := url.Parse(c.DSN())
	if err != nil {
		t.Fatalf("DSN is not a valid URL: %v", err)
	}

	got, _ := u.User.Password()
	if got != password {
		t.Errorf("password did not survive round trip: got %q", got)
	}

	if u.Hostname() != "localhost" {
		t.Errorf("host was corrupted by password: got %q", u.Hostname())
	}

	if u.Path != "/warehouse_db" {
		t.Errorf("unexpected path %q", u.Path)
	}
}
package main

import (
	"errors"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/bcrypt"
)

type configuration struct {
	Host      string                    `toml:"host"`
	Port      int                       `toml:"port"`
	Username  string                    `toml:"username"`
	Password  string                    `toml:"password"`
	PlainText bool                      `toml:"plaintext"`
	Paths     map[string]toml.Primitive `toml:"path"`
}

var config = configuration{
	Host:     "192.168.0.1",
	Port:     9921,
	Username: "scanner",
	Password: "$2y$12$3TwvitKJL3L4/4XVMFFgAOYVCsnj6jZ/cxRBF2/ynbrQPYOEUzqEm", // scanme
}

func (c configuration) CheckPasswd(username, password string) (bool, error) {
	if username == c.Username {
		if c.PlainText {
			if c.Password == password {
				return true, nil
			}
		} else if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err == nil {
			return true, nil
		}
	}
	return false, errors.New("invalid credentials")
}

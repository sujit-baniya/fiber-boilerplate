package config

import (
	"github.com/alexedwards/argon2id"
)

type Hash struct {
	// Argon2id configuration
	Params *argon2id.Params
}

func (d *Hash) Create(password string) (hash string, err error) {
	if d.Params == nil {
		d.Params = argon2id.DefaultParams
	}
	return argon2id.CreateHash(password, d.Params)
}

func (d *Hash) Match(password string, hash string) (match bool, err error) {
	if d.Params == nil {
		d.Params = argon2id.DefaultParams
	}
	return argon2id.ComparePasswordAndHash(password, hash)
}

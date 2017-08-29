package model

import (
	"crypto/rand"
	"errors"

	"github.com/oklog/ulid"
)

type GUID string

func NewGUID() GUID {
	return GUID(ulid.MustNew(ulid.Now(), rand.Reader).String())
}

func (g *GUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		*g = GUID(src)
	case []byte:
		*g = GUID(string(src))
	default:
		return errors.New("invalid source type for GUID scan")
	}

	return nil
}

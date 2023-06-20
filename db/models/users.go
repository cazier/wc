package models

import (
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email   string `gorm:"unique"`
	Salt    Base64 `gorm:"type:string"`
	Hash    Base64 `gorm:"type:string"`
	Session Cookie `gorm:"embedded;embeddedPrefix:session_"`
}

func (u User) IsNil() bool {
	return u.Email == "" && u.Salt == nil && u.Hash == nil && u.Session == Cookie{} && u.Model == gorm.Model{}
}

type Cookie struct {
	Name      string
	Value     string
	CreatedAt time.Time
}

func NewCookie(name, value string) Cookie {
	return Cookie{Name: name, Value: value, CreatedAt: time.Now()}
}

func (c Cookie) IsTooOld(lifetime time.Duration) bool {
	return c.CreatedAt.Add(lifetime).Before(time.Now())
}

type Base64 []byte

func (b *Base64) Scan(value interface{}) error {
	bytes, ok := value.(string)

	if !ok {
		return fmt.Errorf("could not unmarshall database value: %s", value)
	}

	decode, err := base64.StdEncoding.DecodeString(bytes)

	if err == nil {
		*b = Base64(decode)
	}

	return err
}

func (b Base64) Value() (driver.Value, error) {
	resp := base64.StdEncoding.EncodeToString(b)
	return resp, nil
}

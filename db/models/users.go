package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`

	Name    string `gorm:"unique" json:"name"`
	Email   string `gorm:"unique" json:"email"`
	Salt    Base64 `gorm:"type:string" json:"-"`
	Hash    Base64 `gorm:"type:string" json:"-"`
	Session Token  `gorm:"embedded;embeddedPrefix:session_" json:"-"`
	Csrf    Token  `gorm:"embedded;embeddedPrefix:csrf_" json:"-"`
}

func (u User) IsNil() bool {
	return u.Email == "" && u.Salt == nil && u.Hash == nil && u.Session == Token{} && u.Model == gorm.Model{}
}

func (u User) Serialize() map[string]string {
	data := make(map[string]string)

	j, _ := json.Marshal(u)
	json.Unmarshal(j, &data)

	return data
}

type Token struct {
	Name      string
	Value     string
	CreatedAt time.Time
}

func NewToken(name, value string) Token {
	return Token{Name: name, Value: value, CreatedAt: time.Now()}
}

func (t Token) IsValid(lifetime time.Duration) bool {
	return t.CreatedAt.Before(time.Now()) && t.CreatedAt.Add(lifetime).After(time.Now())
}

type Base64 []byte

func (b *Base64) Scan(value interface{}) error {
	plaintext, ok := value.(string)

	if !ok {
		return fmt.Errorf("could not unmarshall database value: %s", value)
	}

	decode, err := base64.StdEncoding.DecodeString(plaintext)

	if err == nil {
		*b = Base64(decode)
	}

	return err
}

func (b Base64) Value() (driver.Value, error) {
	resp := base64.StdEncoding.EncodeToString(b)
	return resp, nil
}

package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserIsNil(t *testing.T) {
	assert := assert.New(t)

	assert.True(User{}.IsNil())
	assert.False(User{Email: "user@email.com"}.IsNil())
}

func TestNewCookie(t *testing.T) {
	assert := assert.New(t)

	start := time.Now()

	cookie := NewToken("session", "secret")
	assert.Equal("secret", cookie.Value)
	assert.WithinDuration(start, cookie.CreatedAt, time.Since(start))
}
func TestTimedTokenIsValid(t *testing.T) {
	assert := assert.New(t)

	duration, _ := time.ParseDuration("-1h")

	assert.False(Token{CreatedAt: time.Now()}.IsValid(duration))
	assert.True(Token{CreatedAt: time.Now()}.IsValid(duration.Abs()))
}

func TestUserSerialize(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Name:    "serialize",
		Email:   "serialize@email.com",
		Salt:    []byte("Salt"),
		Hash:    []byte("Hash"),
		Session: Token{Name: "Session", Value: "Value"},
		Csrf:    Token{Name: "Csrf", Value: "Value"},
	}

	assert.False(user.IsNil())
	assert.EqualValues(map[string]string{"name": "serialize", "email": "serialize@email.com"}, user.Serialize())
}

func TestScanner(t *testing.T) {
	assert := assert.New(t)

	plaintext := "Scanner"
	encoded := "U2Nhbm5lcg=="

	container := Base64(plaintext)
	value, err := container.Value()

	assert.Nil(err)
	assert.EqualValues(encoded, value)

	container = Base64{}
	err = container.Scan(encoded)

	assert.Nil(err)
	assert.EqualValues(plaintext, container)

	container = Base64{}
	err = container.Scan(10)
	assert.ErrorContains(err, "could not unmarshall")
}

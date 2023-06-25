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

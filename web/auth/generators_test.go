package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSalt(t *testing.T) {
	assert := assert.New(t)

	salt := generateSalt()

	assert.NotNil(salt)
	assert.Len(salt, saltLength)
}

func TestGenerateHash(t *testing.T) {
	assert := assert.New(t)

	hash := generateHash("password", generateSalt())

	assert.NotNil(hash)
	assert.Len(hash, hashLength)
}

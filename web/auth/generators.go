package auth

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const memory = 64
const threads = 4
const saltLength = 16
const hashLength = 64
const iterations = 10

func generateSalt() []byte {
	salt := make([]byte, saltLength)
	rand.Read(salt)

	return salt
}

func generateHash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, iterations, memory*1024, threads, hashLength)
}

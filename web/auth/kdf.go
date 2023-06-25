package auth

import (
	"bytes"
	"errors"

	"crypto/rand"

	"github.com/cazier/wc/db/models"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

const memory = 64
const threads = 4
const saltLength = 16
const hashLength = 64
const iterations = 10

var ErrAccountExists error = errors.New("an account with this name or email address already exists")

func create(name, email, password string) (models.User, error) {
	var dest models.User

	tx := db.Where(models.User{Name: name, Email: email}).FirstOrCreate(&dest)

	if (tx.Error == nil && tx.RowsAffected != 1) || errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return dest, ErrAccountExists
	}

	dest.Salt = generateSalt()
	dest.Hash = generateHash(password, dest.Salt)

	save(dest)

	return dest, nil
}

func save(user models.User) {
	db.Save(&user)
}

func retrieve(search models.User) models.User {
	var dest models.User

	tx := db.First(&dest, &search)

	if (tx.Error != nil) || (tx.RowsAffected != 1) {
		return models.User{}
	}

	return dest
}

func isValid(email, password string) bool {
	user := retrieve(models.User{Email: email})

	if user.Hash == nil {
		return false
	}

	if bytes.Equal(generateHash(password, user.Salt), user.Hash) {
		return true
	}

	return false
}

func generateSalt() []byte {
	salt := make([]byte, saltLength)
	rand.Read(salt)

	return salt
}

func generateHash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, iterations, memory*1024, threads, hashLength)
}

package auth

import (
	"bytes"
	"errors"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/web/exceptions"
	"gorm.io/gorm"
)

type AuthenticationDatabase interface {
	create(name, email, password string) (models.User, error)
	save(user models.User)
	retrieve(search models.User) models.User
	isValid(email, password string) bool
}

type RuntimeAuthenticationDatabase struct {
	db *gorm.DB
}

func (r *RuntimeAuthenticationDatabase) create(name, email, password string) (models.User, error) {
	var dest models.User

	tx := r.db.Where(models.User{Name: name, Email: email}).FirstOrCreate(&dest)

	if (tx.Error == nil && tx.RowsAffected != 1) || errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return dest, exceptions.ErrAccountExists
	}

	updateHashSalt(&dest, password)
	r.save(dest)

	return dest, nil
}

func (r *RuntimeAuthenticationDatabase) save(user models.User) {
	r.db.Save(&user)
}

func (r *RuntimeAuthenticationDatabase) retrieve(search models.User) models.User {
	var dest models.User

	tx := r.db.First(&dest, &search)

	if (tx.Error != nil) || (tx.RowsAffected != 1) {
		return models.User{}
	}

	return dest
}

func (r *RuntimeAuthenticationDatabase) isValid(email, password string) bool {
	user := r.retrieve(models.User{Email: email})

	if user.Hash == nil {
		return false
	}

	if bytes.Equal(generateHash(password, user.Salt), user.Hash) {
		return true
	}

	return false
}

func updateHashSalt(u *models.User, password string) {
	salt := generateSalt()
	hash := generateHash(password, salt)

	u.Salt = salt
	u.Hash = hash
}

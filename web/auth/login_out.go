package auth

import (
	"errors"
	"net/http"

	"github.com/cazier/wc/db/models"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) (int, map[string]any) {
	type form struct {
		Name     string `form:"name" binding:"required"`
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
		Confirm  string `form:"confirm" binding:"required"`
		// Csrf     string `form:"csrf"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		return http.StatusNotAcceptable, gin.H{"message": "invalid username or password"}
	}

	if data.Password != data.Confirm {
		return http.StatusNotAcceptable, gin.H{"message": "passwords do not match"}
	}

	user, err := create(data.Name, data.Email, data.Password)

	if err != nil && errors.Is(err, ErrAccountExists) {
		return http.StatusNotAcceptable, gin.H{"message": err.Error()}
	}

	withCookie(c, user)
	return http.StatusFound, nil
}

func Login(c *gin.Context) (int, map[string]any) {
	type form struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
		// Csrf     string `form:"csrf"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		return http.StatusNotAcceptable, gin.H{"message": "invalid username or password"}
	}

	user := retrieve(models.User{Email: data.Email})

	if isValid(user.Email, data.Password) {
		withCookie(c, user)
		return http.StatusFound, nil
	}

	return http.StatusUnauthorized, gin.H{"message": "invalid username or password"}
}

func withCookie(c *gin.Context, user models.User) {
	user.Session = models.NewToken("session", generateCookie())
	save(user)
	addSessionCookie(c, user.Session.Value)
}

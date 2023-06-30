package auth

import (
	"errors"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/web/exceptions"
	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	type form struct {
		Name     string `form:"name" binding:"required"`
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
		Confirm  string `form:"confirm" binding:"required"`
		// Csrf     string `form:"csrf"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		c.Error(exceptions.ErrAccountDetailsInvalid)
		return
	}

	if data.Password != data.Confirm {
		c.Error(exceptions.ErrAccountPasswordMismatch)
		return
	}

	user, err := create(data.Name, data.Email, data.Password)

	if err != nil && errors.Is(err, exceptions.ErrAccountExists) {
		c.Error(exceptions.ErrAccountExists)
		return
	}

	withCookie(c, user)
}

func Login(c *gin.Context) {
	type form struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
		// Csrf     string `form:"csrf"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		c.Error(exceptions.ErrAccountDetailsInvalid)
		return
	}

	user := retrieve(models.User{Email: data.Email})

	if isValid(user.Email, data.Password) {
		withCookie(c, user)
		return
	}

	c.Error(exceptions.ErrUnauthorized)
}

func Update(c *gin.Context) {
	type form struct {
		Email       string `form:"email"`
		Password    string `form:"password" binding:"required"`
		NewPassword string `form:"new"`
		Confirm     string `form:"confirm"`
		Csrf        string `form:"csrf" binding:"required"`
	}

	return
}

func withCookie(c *gin.Context, user models.User) {
	user.Session = models.NewToken("session", generateCookie())
	save(user)
	addSessionCookie(c, user.Session.Value)
}

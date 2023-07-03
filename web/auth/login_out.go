package auth

import (
	"errors"

	"github.com/cazier/wc/db/models"
	"github.com/cazier/wc/web/exceptions"
	"github.com/gin-gonic/gin"
)

type Authentication interface {
	Create(c *gin.Context)
	Login(c *gin.Context)
	Update(c *gin.Context)
	withCookie(c *gin.Context, user models.User)
}

type RuntimeAuthentication struct{}

func (r *RuntimeAuthentication) Create(c *gin.Context) {
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

	user, err := db.create(data.Name, data.Email, data.Password)

	if err != nil && errors.Is(err, exceptions.ErrAccountExists) {
		c.Error(exceptions.ErrAccountExists)
		return
	}

	r.withCookie(c, user)
}

func (r *RuntimeAuthentication) Login(c *gin.Context) {
	type form struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		c.Error(exceptions.ErrAccountDetailsInvalid)
		return
	}

	user := db.retrieve(models.User{Email: data.Email})

	if db.isValid(user.Email, data.Password) {
		r.withCookie(c, user)
		return
	}

	c.Error(exceptions.ErrUnauthorized)
}

func (r *RuntimeAuthentication) Update(c *gin.Context) {
	type form struct {
		Email       string `form:"email"`
		Password    string `form:"password" binding:"required"`
		NewPassword string `form:"new"`
		Confirm     string `form:"confirm"`
		Csrf        string `form:"csrf" binding:"required"`
	}

	data := form{}

	if err := c.ShouldBind(&data); err != nil {
		c.Error(exceptions.ErrAccountDetailsInvalid)
		return
	}

	// user, _ := getUser(c)

	// if isValid(user.Email, data.Password) {
	// 	if data.Password == data.Confirm {
	// 		user.
	// 	}

	// }

}

func (r *RuntimeAuthentication) withCookie(c *gin.Context, user models.User) {
	user.Session = models.NewToken("session", generateCookie())
	db.save(user)
	addSessionCookie(c, user.Session.Value)
}

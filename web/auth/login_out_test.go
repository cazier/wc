package auth

import (
	"testing"

	"github.com/cazier/wc/web/exceptions"
	"github.com/stretchr/testify/assert"

	test "github.com/cazier/wc/testing"
)

// LOGIN_OUT

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	mf := test.NewMockForm()
	type msa map[string]any

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "pass"}))
	assert.Empty(mf.Errors())

	Create(mf.Form(msa{"email": "create@email.com", "confirm": "password"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountDetailsInvalid.Error())

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "badpass"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountPasswordMismatch.Error())

	Create(mf.Form(msa{"name": "create", "email": "create@email.com", "password": "pass", "confirm": "pass"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountExists.Error())
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	mf := test.NewMockForm()
	type msa map[string]any

	Create(mf.Form(msa{"email": "login@email.com", "password": "password", "confirm": "password"}))

	Login(mf.Form(msa{"email": "login@email.com", "password": "password"}))
	assert.Empty(mf.Errors())

	Login(mf.Form(msa{"email": "login@email.com"}))
	assert.Contains(mf.Errors(), exceptions.ErrAccountDetailsInvalid.Error())

	Login(mf.Form(msa{"email": "login@email.com", "password": "wrongpassword"}))
	assert.Contains(mf.Errors(), exceptions.ErrUnauthorized.Error())
}

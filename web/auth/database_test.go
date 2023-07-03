package auth

import (
	"testing"

	"github.com/cazier/wc/db/models"
	test "github.com/cazier/wc/testing"
	"github.com/stretchr/testify/assert"
)

var tf AuthenticationDatabase

func init() {
	//TODO Make m.OpenDB more generic so that I don't need to sue this whole mock.
	m = test.NewMock(
		&test.MockOptions{
			Models: []any{
				&models.User{},
			},
		},
	)

	tf = &RuntimeAuthenticationDatabase{m.Database}
}

func TestDbCreate(t *testing.T) {
	assert := assert.New(t)

	_, err := tf.create("dbcreate", "dbcreate@email.com", "password")
	assert.Nil(err)

	_, err = tf.create("new_dbcreate", "dbcreate@email.com", "passwordpassword")
	assert.Errorf(err, "an account with this name or email address already exists")

	_, err = tf.create("dbcreate", "new_dbcreate@email.com", "passwordpassword")
	assert.Errorf(err, "an account with this name or email address already exists")

	_, err = tf.create("dbcreate", "dbcreate@email.com", "password")
	assert.Errorf(err, "an account with this name or email address already exists")
}

func TestDbSave(t *testing.T) {
	assert := assert.New(t)

	tf.create("", "dbsave@email.com", "password")
	model := tf.retrieve(models.User{Email: "dbsave@email.com"})

	model.Name = "dbsave"

	tf.save(model)

	model = tf.retrieve(models.User{Name: "dbsave"})
	assert.False(model.IsNil())
}

func TestRetrieve(t *testing.T) {
	assert := assert.New(t)

	user := tf.retrieve(models.User{Email: "retrieve@email.com"})
	assert.Nil(user.Salt)
	assert.Nil(user.Hash)

	user, _ = tf.create("retrieve", "retrieve@email.com", "password")
	assert.NotNil(user.Salt)
	assert.NotNil(user.Hash)

	assert.False(user.IsNil())

	retrieved := tf.retrieve(models.User{Email: "retrieve@email.com"})
	assert.EqualExportedValues(user, retrieved)
}

func TestIsValid(t *testing.T) {
	assert := assert.New(t)

	tf.create("valid", "isvalid@email.com", "password")

	assert.True(tf.isValid("isvalid@email.com", "password"))

	assert.False(tf.isValid("isvalid@email.com", "wrongpassword"))
	assert.False(tf.isValid("isinvalid@email.com", "password"))
}

func TestUpdateHashSalt(t *testing.T) {
	assert := assert.New(t)

	user := models.User{}
	assert.Nil(user.Hash)
	assert.Nil(user.Salt)

	updateHashSalt(&user, "password")

	assert.Len(user.Hash, hashLength)
	assert.Len(user.Salt, saltLength)
}

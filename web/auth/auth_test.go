package auth

import (
	"os"
	"testing"

	test "github.com/cazier/wc/testing"

	"github.com/cazier/wc/db/models"
)

var m test.Mock

func TestMain(tm *testing.M) {
	m = test.NewMock(
		&test.MockOptions{
			Callback: Init,
			Models: []any{
				&models.User{},
			},
		},
	)

	os.Exit(tm.Run())
}

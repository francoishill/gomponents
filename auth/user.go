package auth

import (
	"github.com/francoishill/gomponents/user"
)

//User extends the user.User interface
type User interface {
	user.User

	PasswordHash() string
	MagicLoginToken() *string
}

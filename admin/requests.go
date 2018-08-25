package admin

import "github.com/francoishill/gomponents/auth"

type RequestFactory interface {
	AddUser() AddUserRequest
}

type AddUserRequest interface {
	Validate() error
	ToUser(passwordHash string) auth.User
}

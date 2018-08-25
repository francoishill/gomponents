package admin

import "github.com/francoishill/gomponents/user"

type ResponseFactory interface {
	User(user user.User) UserResponse
}

type UserResponse interface{}

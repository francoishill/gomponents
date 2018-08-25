package anonymous

import "github.com/francoishill/gomponents/auth"

type ResponseFactory interface {
	LoggedIn(user auth.User, token string) LoggedInResponse
}

type LoggedInResponse interface{}

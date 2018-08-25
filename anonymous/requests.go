package anonymous

import "github.com/francoishill/gomponents/auth"

type RequestFactory interface {
	Register() RegisterRequest
	Login() LoginRequest
	MagicLogin() MagicLoginRequest
}

type RegisterRequest interface {
	Validate() error
	ToUser() auth.User
}

type LoginRequest interface {
	Validate() error
	Password() string
	ToUser() auth.User
}

type MagicLoginRequest interface {
	Validate() error
	Token() string
	ToUser() auth.User
}

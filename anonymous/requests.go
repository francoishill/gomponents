package anonymous

import "github.com/francoishill/gomponents/auth"

type RequestFactory interface {
	Register() RegisterRequest
	Login() LoginRequest
	MagicLogin() MagicLoginRequest
}

type RegisterRequest interface {
	Validate() error
	LoadUser() (auth.User, error)
}

type LoginRequest interface {
	Validate() error
	Password() string
	LoadUser() (auth.User, error)
}

type MagicLoginRequest interface {
	Validate() error
	Token() string
	LoadUser() (auth.User, error)
}

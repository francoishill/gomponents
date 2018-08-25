package user

type User interface {
	ID() string

	IsAdmin() bool
}

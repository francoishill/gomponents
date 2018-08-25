package user

type Validation interface {
	User(user User) error
}

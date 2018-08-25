package user

type Factory interface {
	Repo() Repo
}

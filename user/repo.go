package user

type Repo interface {
	ApplySchema() error

	IsDupErr(err error) bool

	Add(user User) error
	Get(id string) (User, error)
	List() ([]User, error)
}

type RepoFactory interface {
	Repo() Repo
}

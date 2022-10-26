package user

// Struct for DB
type User struct {
	Id   int64
	Name string
}

// Interface required for each new implementation
type UserStorage interface {
	GetUser(id int64) (*User, error)
	GetUsers() []*User
	AddUser(name string) error
	UpdateUser(name string) error
	DeleteUser(id int64) error
}

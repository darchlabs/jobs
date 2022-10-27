package user

// Struct for DB
type User struct {
	Id   string
	Name string
}

// Interface required for each new implementation
type UserStorage interface {
	GetUser(id string) (*User, error)
	GetUsers() []*User
	AddUser(name string) error
	UpdateUser(name string) error
	DeleteUser(id string) error
}

package userstorage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/darchlabs/jobs/internal/storage"
	"github.com/darchlabs/jobs/internal/user"
	"github.com/teris-io/shortid"
)

type US struct {
	storage *storage.S
}

func New(s *storage.S) *US {
	return &US{
		storage: s,
	}
}

func (us *US) GetUser(id string) (*user.User, error) {
	data, err := us.storage.DB.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}

	var u *user.User

	err = json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (us *US) GetUsers() []*user.User {
	users := make([]*user.User, 0)

	iter := us.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var p *user.User
		err := json.Unmarshal(iter.Value(), &p)

		if err != nil {
			log.Printf("Error while iterating db: %v \n", err)
			return nil
		}

		users = append(users, p)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		log.Printf("Error while iterating db: %v \n", err)
		return nil
	}

	return users
}

func (us *US) AddUser(name string) error {
	// Validate param
	if name == "" {
		return fmt.Errorf("%s", "name param string is empty")
	}

	// Generate new id (the params is for always start with a letter)
	id, err := shortid.New(1, "-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_", 2342)
	if err != nil {
		return err
	}

	u := &user.User{
		Id:   id.String(),
		Name: name,
	}

	b, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = us.storage.DB.Put([]byte(u.Id), b, nil)
	if err != nil {
		return err
	}

	return nil
}

func (us *US) UpdateUser(id string, name string) error {
	// Validate param
	if name == "" {
		return fmt.Errorf("%s", "name param string is empty")
	}

	// Validate the user exists
	data, err := us.storage.DB.Get([]byte(id), nil)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("%s", "No User exists for the given id")
	}

	// Convert user to bytes
	var u *user.User
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}

	// Update in DB
	err = us.storage.DB.Put([]byte(id), b, nil)
	if err != nil {
		return err
	}

	return nil
}

func (us *US) DeleteUser(id string) error {
	// Delete user from DB
	err := us.storage.DB.Delete([]byte(id), nil)
	if err != nil {
		return err
	}

	return nil
}

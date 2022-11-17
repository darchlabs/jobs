package storage

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type S struct {
	DB *leveldb.DB
}

func New(filepath string) (*S, error) {
	if filepath == "" {
		return nil, fmt.Errorf("%s", "Empty filepath")
	}

	// read db from file
	db, err := leveldb.OpenFile(fmt.Sprintf("./%s", filepath), nil)
	if err != nil {
		return nil, err
	}

	return &S{
		DB: db,
	}, nil
}

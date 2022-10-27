package jobstorage

import (
	"encoding/json"
	"fmt"

	"github.com/darchlabs/jobs/internal/storage"
	providerstorage "github.com/darchlabs/jobs/internal/storage/provider"
	userstorage "github.com/darchlabs/jobs/internal/storage/user"
)

// Main DB table

type JobStorage struct {
	storage         *storage.S
	SmartContractId string
	Event           string
	UserId          string
	ProviderId      uint8
	SynchronizerId  string // Indicates the id of the synchronizer db

	/// @notice: Provider-User related fields
	Setup        bool
	Working      bool
	NeedsFunding bool
}

func New(s *storage.S) *JobStorage {
	return &JobStorage{
		storage: s,
	}
}

func checkStrings(params []string) error {
	for k, v := range params {
		if v == "" {
			return fmt.Errorf("%d param is empty", k)
		}
	}

	return nil
}

func (js *JobStorage) AddUserProvider(
	scAddress string,
	event string,
	userId string,
	providerId uint8) error {
	// Validate params
	paramsArr := make([]string, 0)
	paramsArr = append(paramsArr, scAddress, event, userId)
	err := checkStrings(paramsArr)
	if err != nil {
		return err
	}

	// Check smart contract exists
	sc, err := js.getContract(scAddress)
	if err != nil {
		return err
	}

	if sc == "" {
		return fmt.Errorf("%s", "smart contract doesn't exists for given id")
	}

	// Check user and provider exist for the given id
	us := userstorage.New(js.storage)
	user, err := us.GetUser(userId)
	if err != nil {
		return err
	}

	// TODO(nb): Don't know if this is necessary. If is nil, it returns an err?
	if user == nil {
		return fmt.Errorf("%s", "User not found for given id")
	}

	ps := providerstorage.New(js.storage)
	provider, err := ps.GetImplementation(providerId)
	if err != nil {
		return err
	}

	if provider == nil {
		return fmt.Errorf("%s", "Provider not found for given id")
	}

	setup := provider.Implemetation.Setup(scAddress)
	if setup != true {
		return fmt.Errorf("%s", "Error in setup of the provider implementation")
	}

	working, _, needsFunding := provider.Implemetation.GetState(provider.Name)

	jobStorage := JobStorage{
		SmartContractId: sc,
		Event:           event,
		UserId:          userId,
		ProviderId:      providerId,
		SynchronizerId:  "", // TODO(nb): How to link this to syncho=rnoizer?
		Setup:           setup,
		Working:         working,
		NeedsFunding:    needsFunding,
	}

	// JSON.stringify
	b, err := json.Marshal(jobStorage)
	if err != nil {
		return err
	}

	// Save in database
	err = js.storage.DB.Put([]byte(scAddress), b, nil)
	if err != nil {
		return err
	}

	return nil

}

func (js *JobStorage) getContract(id string) (string, error) {
	return "", nil
}

// func (js *JobStorage) UpdateUserProvider(
// 	scAddress string,
// 	event string,
// 	userId string,
// 	providerId int8) error {
// 	// Validate params
// 	err := checkStrings(scAddress, "scAddress")
// 	if err != nil {
// 		return err
// 	}

// 	err = checkStrings(event, "event")
// 	if err != nil {
// 		return err
// 	}

// 	// Check user and provider exist for the given id
// 	us := userstorage.New(js.storage)
// 	user, err := us.GetUser(userId)
// 	if err != nil {
// 		return err
// 	}

// 	// TODO(nb): Don't know if this is necessary. If is nil, it returns an err?
// 	if user == nil {
// 		return fmt.Errorf("%s", "User not found for given id")
// 	}

// 	ps := providerstorage.New(js.storage)
// 	provider, err := ps.GetImplementation(providerId)
// 	if err != nil {
// 		return err
// 	}

// 	if provider == nil {
// 		return fmt.Errorf("%s", "Provider not found for given id")
// 	}

// 	// TODO(nb): put implementation
// 	return nil
// }

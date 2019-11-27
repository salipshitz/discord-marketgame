package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// User data
type User struct {
	mux     sync.Mutex
	ID      string          `json:"id"`
	Balance float64         `json:"balance"`
	Shares  map[string]uint `json:"shares"`
}

var memUsers map[string]*User = map[string]*User{}

func getUserFilePath(id string) (*string, error) {
	userFilePath, err := filepath.Abs(filepath.Join("users", id+".json"))

	if err != nil {
		return nil, err
	}

	return &userFilePath, nil
}

func GetUser(id string) (*User, error) {
	memUser, ok := memUsers[id]

	var user User

	if ok {
		return memUser, nil
	}

	userFilePath, err := getUserFilePath(id)

	if err != nil {
		return nil, err
	}

	userFile, err := os.Open(*userFilePath)

	defer userFile.Close()

	if err != nil {
		user = User{
			ID:      id,
			mux:     sync.Mutex{},
			Balance: 100000,
			Shares:  map[string]uint{},
		}

		err := user.Save()

		memUsers[id] = &user

		if err != nil {
			return &user, nil
		}

		return nil, err
	}

	userFileContent, err := ioutil.ReadAll(userFile)

	if err != nil {
		return nil, err
	}

	user = User{
		mux: sync.Mutex{},
	}

	err = json.Unmarshal(userFileContent, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Save user data
func (u *User) Save() error {
	userFilePath, err := getUserFilePath(u.ID)

	if err != nil {
		return err
	}

	file, err := os.Create(*userFilePath)
	defer file.Close()

	if err != nil {
		return err
	}

	marshalled, err := json.Marshal(u)

	if err != nil {
		return err
	}

	_, err = file.Write(marshalled)

	if err != nil {
		return err
	}

	return nil
}

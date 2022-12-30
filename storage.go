package main

import (
	"encoding/json"
	"os"

	"github.com/michimani/gotwi/resources"
)


func SaveUsers(users []resources.User, filename string) {
	json, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, json, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadUsers(filename string) []resources.User {
	users := []resources.User{}
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &users)
	return users
}

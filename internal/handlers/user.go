package handlers

import (
	"fmt"
	"net/http"
)

type User struct{}

func (o *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creates a user")
}

func (o *User) ListUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all users")
}

func (o *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Find a user by their ID")
}

func (o *User) UpdateUserbyID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updates userinfo")
}

func (o *User) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deletes a user")
}

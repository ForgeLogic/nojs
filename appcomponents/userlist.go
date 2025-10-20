//go:build js || wasm
// +build js wasm

package appcomponents

import (
	"github.com/vcrobe/nojs/runtime"
)

type User struct {
	ID   int
	Name string
}

type UserList struct {
	runtime.ComponentBase
	Users []User
	Title string
}

func (u *UserList) AddUser() {
	if u.Users == nil {
		u.Users = []User{}
	}
	// Add a new user to the list
	newID := len(u.Users) + 1
	u.Users = append(u.Users, User{
		ID:   newID,
		Name: "User " + string(rune('A'+newID-1)),
	})
	u.StateHasChanged()
}

func (u *UserList) ClearUsers() {
	u.Users = []User{}
	u.StateHasChanged()
}

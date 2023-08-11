package types

import (
	"time"
)

type Account struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    int `json:"gender"`
	Dni       int64 `json:"dni"`
	Cbu       int64 `json:"cbu"`
	Balance   int64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAccount(firstName string, lastName string, gender int, dni int64) *Account {
return &Account{
	FirstName: firstName,
	LastName: lastName,
	Gender: gender,
	Dni: dni,
}
}
package types

import (
	"fmt"
	"time"

	"github.com/Gonzapepe/bank-api/helper"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Gender int `json:"gender"`
	Dni int64 `json:"dni"`
	Password string `json:"password"`
}

type Account struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    int `json:"gender"`
	Dni       int64 `json:"dni"`
	Password string `json:"password"`
	Cuit       int64 `json:"cuit"`
	Balance   int64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Dni int64 `json:"dni"`
	Password string `json:"password"`
}

func NewAccount(firstName string, lastName string, gender int, dni int64, password string) *Account {

cuit, err := helper.Cuit(dni, gender)

if err != nil {
	fmt.Println(err.Error())
}

return &Account{
	FirstName: firstName,
	LastName: lastName,
	Gender: gender,
	Dni: dni,
	Cuit: cuit,
	Password: password,
}
}
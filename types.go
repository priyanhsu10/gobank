package main

import (
	"math/rand"
	"time"
)

type AccountCreateReuqust struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type AccountUdpateRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}
type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"fistname"`
	LastName  string    `json:"lastname"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}

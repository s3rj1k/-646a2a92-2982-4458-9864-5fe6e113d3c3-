package main

import (
	"golang.org/x/crypto/bcrypt"
)

// Account describes data needed to create user account in DB.
type Account struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// GetSecureAccount returns Account object with hashed and salted password.
func (account Account) GetSecureAccount() Account {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(account.Password),
		bcrypt.MinCost,
	)
	if err != nil {
		panic(err)
	}

	return Account{
		Email:    account.Email,
		Password: string(hash),
	}
}

package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestLetter(t *testing.T) *Letter {
	return &Letter{
		UserID: 1,
		Name:   "Письмо 1",
		Body:   "Тело письма",
	}
}

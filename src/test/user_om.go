package test

import (
	"github.com/fernetbalboa/arqweb/src/domain"
)

func DefaultUserString() string {
	return `{
	"username": "testName",
	"password": "testPass"
}`
}

func DefaultUser() *domain.User {
	return &domain.User{
		Username: "testName",
		Password: "testPass",
		IsAdmin:  false,
	}
}

func DefaultUserWithId() *domain.User {
	return &domain.User{
		Id:       NewDocumentId(),
		Username: "testName",
		Password: "testPass",
		IsAdmin:  false,
	}
}

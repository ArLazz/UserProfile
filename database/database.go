package database

import (
	"userprofile/models"

	"github.com/go-openapi/strfmt"
)

type Database interface {
	AddUser(user *models.User) error
	GetUserByID(id strfmt.UUID) (*models.User, error)
	UpdateUser(id strfmt.UUID, newUser *models.User) error
	DeleteUser(id strfmt.UUID) error
	ListUsers() []*models.User
	AuthenticateUser(username, password string) (*models.User, error)
}
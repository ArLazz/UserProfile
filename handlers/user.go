package handlers

import (
	"fmt"
	"userprofile/database"
	"userprofile/models"
	"userprofile/restapi/operations"

	"github.com/google/uuid"
)

var db = database.NewInMemoryDB()

func GetUser(params operations.GetUserParams) []*models.User {
	return db.ListUsers()
}

func AddUser(params operations.PostUserParams) error {
	if err := checkUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	if err := db.AddUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	return nil
}

func GetUserByID(params operations.GetUserIDParams) (*models.User, error) {
	if !isValidUUID(string(params.ID)) {
		return nil, fmt.Errorf("error to get user by id: invalid uuid")
	}

	user, err := db.GetUserByID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("error to get user by id: %w", err)
	}

	return user, nil
}

func DeleteUserByID(params operations.DeleteUserIDParams) error {
	if !isValidUUID(string(params.ID)) {
		return fmt.Errorf("error to delete user by id: invalid uuid")
	}

	if err := db.DeleteUser(params.ID); err != nil {
		return fmt.Errorf("error to delete user by id: %w", err)
	}

	return nil
}

func UpdateUserByID(params operations.PutUserIDParams) error {
	if !isValidUUID(string(params.ID)) {
		return fmt.Errorf("error to update user by id: invalid uuid")
	}

	if err := checkUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	if err := db.UpdateUser(params.ID, params.User); err != nil {
		return fmt.Errorf("error to update user by id: %w", err)
	}

	return nil
}

func AuthenticateUser(username, password string) (interface{}, error) {
	user, err := db.AuthenticateUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("error to authenticate user: %w", err)
	}
	return user, nil
}

func checkUser(user *models.User) error {
	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("empty username or password")
	}

	if !isValidUUID(string(user.ID)) {
		return fmt.Errorf("empty ID")
	}

	return nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

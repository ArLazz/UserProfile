package handlers

import (
	"fmt"
	"userprofile/database"
	"userprofile/models"
	"userprofile/restapi/operations"

	"github.com/google/uuid"
)
type Handler struct {
	db database.Database
}

func NewHandler(db database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetUser(params operations.GetUserParams) []*models.User {
	return h.db.ListUsers()
}

func (h *Handler) AddUser(params operations.PostUserParams) error {
	if err := checkUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	if err := h.db.AddUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	return nil
}

func (h *Handler) GetUserByID(params operations.GetUserIDParams) (*models.User, error) {
	if !isValidUUID(string(params.ID)) {
		return nil, fmt.Errorf("error to get user by id: invalid uuid")
	}

	user, err := h.db.GetUserByID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("error to get user by id: %w", err)
	}

	return user, nil
}

func (h *Handler) DeleteUserByID(params operations.DeleteUserIDParams) error {
	if !isValidUUID(string(params.ID)) {
		return fmt.Errorf("error to delete user by id: invalid uuid")
	}

	if err := h.db.DeleteUser(params.ID); err != nil {
		return fmt.Errorf("error to delete user by id: %w", err)
	}

	return nil
}

func (h *Handler) UpdateUserByID(params operations.PutUserIDParams) error {
	if !isValidUUID(string(params.ID)) {
		return fmt.Errorf("error to update user by id: invalid uuid")
	}

	if err := checkUser(params.User); err != nil {
		return fmt.Errorf("error to add user: %s", err)
	}

	if err := h.db.UpdateUser(params.ID, params.User); err != nil {
		return fmt.Errorf("error to update user by id: %w", err)
	}

	return nil
}

func (h *Handler) AuthenticateUser(username, password string) (interface{}, error) {
	user, err := h.db.AuthenticateUser(username, password)
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

package database

import (
	"errors"
	"fmt"
	"sync"
	"userprofile/models"

	"github.com/google/uuid"

	"github.com/go-openapi/strfmt"
)

type InMemoryDB struct {
    mu      sync.RWMutex
    users   map[strfmt.UUID]*models.User
	usernames map[string]strfmt.UUID
}

func NewInMemoryDB() *InMemoryDB {
    db := &InMemoryDB{
        users:   make(map[strfmt.UUID]*models.User),
		usernames: make(map[string]strfmt.UUID),
    }
    return db.Init()
}

func (db *InMemoryDB) Init() *InMemoryDB{
    adminID := strfmt.UUID(uuid.New().String())
    db.users[adminID] = &models.User{
        ID:       adminID,
        Email:    "admin@example.com",
        Username: "admin",
        Password: "adminpass",
        Admin:    true,
    }
    db.usernames["admin"] = adminID
    return db
}

func (db *InMemoryDB) AddUser(user *models.User) error {
    db.mu.Lock()
    defer db.mu.Unlock()
    
    if db.usernames[user.Username] != ""{
        return fmt.Errorf("username already exists")
    }

	db.usernames[user.Username] = user.ID
    db.users[user.ID] = user
    return nil
}

func (db *InMemoryDB) GetUserByID(id strfmt.UUID) (*models.User, error) {
    db.mu.RLock()
    defer db.mu.RUnlock()

    if user, exists := db.users[id]; exists {
        return user, nil
    }

    return nil, fmt.Errorf("user doesn't exists") 
}

func (db *InMemoryDB) UpdateUser(id strfmt.UUID, newUser *models.User) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if oldUser, exists := db.users[id]; exists {
		delete(db.usernames, oldUser.Username)
		db.usernames[newUser.Username] = id
        db.users[id] = newUser

        return nil
    }
    return fmt.Errorf("user doesn't exists")
}

func (db *InMemoryDB) DeleteUser(id strfmt.UUID) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if user, exists := db.users[id]; exists {
        delete(db.usernames, user.Username)
        delete(db.users, id)
        return nil
    }
    return fmt.Errorf("user doesn't exists")
}

func (db *InMemoryDB) ListUsers() []*models.User {
    db.mu.RLock()
    defer db.mu.RUnlock()

    users := make([]*models.User, 0, len(db.users))
    for _, user := range db.users {
        users = append(users, user)
    }
    return users
}

func (db*InMemoryDB) AuthenticateUser(username, password string) (*models.User, error) {
    db.mu.RLock()
    defer db.mu.RUnlock()

    userID, exists := db.usernames[username]
    if !exists || db.users[userID].Password != password {
        return nil, errors.New("invalid username or password")
    }

    return db.users[userID], nil
}
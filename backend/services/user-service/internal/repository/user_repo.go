package repository

import (
	"context"
	"user-service/internal/database"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID      int
	Auth0ID string
	Email   string
}

func GetUserByAuth0ID(auth0ID string) (*User, error) {
	db := database.GetDB()
	row := db.QueryRow(context.Background(), "SELECT id, auth0_id, email FROM users WHERE auth0_id = $1", auth0ID)

	var user User
	err := row.Scan(&user.ID, &user.Auth0ID, &user.Email)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(auth0ID, email string) error {
	db := database.GetDB()
	_, err := db.Exec(context.Background(), "INSERT INTO users (auth0_id, email) VALUES ($1, $2)", auth0ID, email)
	return err
}

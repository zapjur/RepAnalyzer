package repository

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/database"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID      int
	Auth0ID string
	Email   string
}

type Video struct {
	URL       string
	CreatedAt time.Time
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

func SaveUploadedVideo(auth0ID, url, exercise string) error {
	db := database.GetDB()

	var userID int
	err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id = $1", auth0ID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("could not find user with auth0_id %s: %w", auth0ID, err)
	}

	_, err = db.Exec(context.Background(), `
		INSERT INTO videos (user_id, url, exercise_name)
		VALUES ($1, $2, $3)
	`, userID, url, exercise)

	return err
}

func GetUserVideosByExercise(auth0ID, exercise string) ([]Video, error) {
	db := database.GetDB()

	var userID int
	err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id = $1", auth0ID).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user with auth0_id %s: %w", auth0ID, err)
	}

	rows, err := db.Query(context.Background(), `
		SELECT url
		FROM videos
		WHERE user_id = $1 AND exercise_name = $2
		ORDER BY created_at DESC
	`, userID, exercise)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err = rows.Scan(&v.URL, &v.CreatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}

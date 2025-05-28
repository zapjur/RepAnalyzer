package repository

import (
	"context"
	"db-service/internal/database"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID      int
	Auth0ID string
	Email   string
}

type Video struct {
	ObjectKey string
	Bucket    string
	CreatedAt time.Time
	ID        int64
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

func SaveUploadedVideo(auth0ID, bucket, objectKey, exercise string) (int64, error) {
	db := database.GetDB()

	var userID int
	err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id = $1", auth0ID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("could not find user with auth0_id %s: %w", auth0ID, err)
	}

	var videoID int64
	err = db.QueryRow(context.Background(), `
		INSERT INTO videos (user_id, bucket, object_key, exercise_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, bucket, objectKey, exercise).Scan(&videoID)
	if err != nil {
		return 0, fmt.Errorf("could not insert video: %w", err)
	}

	return videoID, nil
}

func GetUserVideosByExercise(auth0ID, exercise string) ([]Video, error) {
	db := database.GetDB()

	var userID int
	err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE auth0_id = $1", auth0ID).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("could not find user with auth0_id %s: %w", auth0ID, err)
	}

	rows, err := db.Query(context.Background(), `
		SELECT object_key, bucket, created_at, id
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
		if err = rows.Scan(&v.ObjectKey, &v.Bucket, &v.CreatedAt, &v.ID); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}

func SaveAnalysis(videoID int64, analysisType, bucket, objectKey string) (int64, error) {
	db := database.GetDB()

	var analysisID int64
	err := db.QueryRow(context.Background(), `
		INSERT INTO video_analysis (video_id, type, bucket, object_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, videoID, analysisType, bucket, objectKey).Scan(&analysisID)
	if err != nil {
		return 0, fmt.Errorf("could not insert analysis for video_id %d: %w", videoID, err)
	}

	return analysisID, nil
}

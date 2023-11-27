package services

import (
	"context"
	"errors"
	"net/http"
)

type User struct{}

func (u *User) Follow(followId string, sessionId string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	if err != nil {
		return http.StatusUnauthorized, errors.New("unauthorized")
	}

	if userId == followId {
		return http.StatusBadRequest, errors.New("cantFollowSelf")
	}

	query := `
		INSERT INTO follow_relationships (followee_id, follower_id)
		VALUES ($1, $2)
	`

	_, err = db.ExecContext(ctx, query, followId, userId)

	if err != nil {
		return http.StatusInternalServerError, errors.New("unableToFollow")
	}

	return http.StatusOK, nil
}

func (u *User) Unfollow(unfollowId string, sessionId string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	if err != nil {
		return http.StatusUnauthorized, errors.New("unauthorized")
	}

	if userId == unfollowId {
		return http.StatusBadRequest, errors.New("cantUnfollowSelf")
	}

	query := `
		DELETE FROM follow_relationships
		WHERE followee_id = $1 AND follower_id = $2
	`

	_, err = db.ExecContext(ctx, query, unfollowId, userId)

	if err != nil {
		return http.StatusInternalServerError, errors.New("unableToUnfollow")
	}

	return http.StatusOK, nil
}

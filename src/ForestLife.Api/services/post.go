package services

import (
	"context"
	"errors"
	"time"
)

var auth Auth

type Post struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) LikePost(postId string, sessionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	query := `
		INSERT INTO post_likes (post_id, user_id)
		VALUES ($1, $2)
	`

	_, err = db.ExecContext(ctx, query, postId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (p *Post) UnlikePost(postId string, sessionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	query := `
		DELETE FROM post_likes
		WHERE post_id = $1 AND user_id = $2
	`

	_, err = db.ExecContext(ctx, query, postId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (p *Post) GetPosts(authorId string) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, text, image, author_id, created_at, updated_at
		FROM posts
		WHERE author_id = $1
	`

	rows, err := db.QueryContext(ctx, query, authorId)

	if err != nil {
		return nil, err
	}

	var posts []*Post

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.Image,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *Post) GetPostById(id string) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, text, image, author_id, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	row := db.QueryRowContext(ctx, query, id)

	var post Post
	err := row.Scan(
		&post.ID,
		&post.Text,
		&post.Image,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *Post) CreatePost(post Post, sessionId string) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	if err != nil {
		return nil, errors.New("unauthorized")
	}

	post.AuthorID = userId

	query := `
		INSERT INTO posts (text, image, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		returning *
	`

	_, err = db.ExecContext(
		ctx,
		query,
		post.Text,
		post.Image,
		userId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *Post) UpdatePost(id string, body Post, sessionId string) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	if err != nil {
		return nil, errors.New("unauthorized")
	}

	query := `
		UPDATE posts
		SET text = $1, image = $2, updated_at = $3
		WHERE id = $4 AND author_id = $5
		returning *
	`

	_, err = db.ExecContext(
		ctx,
		query,
		body.Text,
		body.Image,
		time.Now(),
		id,
		userId,
	)

	if err != nil {
		return nil, err
	}

	return &body, nil
}

func (p *Post) DeletePost(id string, sessionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	userId, err := auth.GetUserId(ctx, sessionId)

	if err != nil {
		return errors.New("unauthorized")
	}

	query := `
		DELETE FROM posts
		WHERE id = $1 AND author_id = $2
	`

	_, err = db.ExecContext(ctx, query, id, userId)

	if err != nil {
		return err
	}

	return nil
}

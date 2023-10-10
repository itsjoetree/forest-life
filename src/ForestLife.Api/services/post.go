package services

import (
	"context"
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TODO: Research how to properly handle query params in services
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

func (p *Post) CreatePost(post Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO posts (text, image, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		post.Text,
		post.Image,
		post.AuthorID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *Post) UpdatePost(id string, body Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE posts
		SET text = $1, image = $2, updated_at = $3
		WHERE id = $4
		returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		body.Text,
		body.Image,
		time.Now(),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &body, nil
}

func (p *Post) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	_, err := db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

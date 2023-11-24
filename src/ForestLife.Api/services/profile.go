package services

import (
	"context"
	"time"
)

type ProfileTheme string

const (
	Forest   ProfileTheme = "forest"
	Dark                  = "dark"
	Standard              = "standard"
)

type Profile struct {
	ID        string       `json:"id"`
	Username  string       `json:"username"`
	Nickname  string       `json:"nickname"`
	Email     string       `json:"email"`
	Theme     ProfileTheme `json:"theme"`
	Password  string       `json:"password,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (p *Profile) GetProfileByUserId(userId string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT profiles.id, username, nickname, email, theme
		FROM profiles
		INNER JOIN users ON profiles.id = users.profile_id
		WHERE users.id = $1
	`

	row := db.QueryRowContext(ctx, query, userId)

	var profile Profile
	err := row.Scan(
		&profile.ID,
		&profile.Username,
		&profile.Nickname,
		&profile.Email,
		&profile.Theme,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

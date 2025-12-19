package model

import "time"

type OAuthProvider struct {
	ID           int       `db:"id"`
	StudentID    int       `db:"student_id"`
	ProviderName string    `db:"provider_name"`
	ExternalID   string    `db:"external_id"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
}

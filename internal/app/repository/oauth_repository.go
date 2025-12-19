package repository

import (
	"database/sql"
	"fmt"
	"student-portal/internal/app/model"
)

type OAuthRepository struct {
	db *sql.DB
}

func NewOAuthRepository(db *sql.DB) *OAuthRepository {
	return &OAuthRepository{db: db}
}

func (r *OAuthRepository) FindByProviderAndExternalID(provider, externalID string) (*model.OAuthProvider, error) {
	query := `SELECT id, student_id, provider_name, external_id, email, created_at 
              FROM oauth_providers 
              WHERE provider_name = $1 AND external_id = $2`

	var oauth model.OAuthProvider
	err := r.db.QueryRow(query, provider, externalID).Scan(
		&oauth.ID, &oauth.StudentID, &oauth.ProviderName,
		&oauth.ExternalID, &oauth.Email, &oauth.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска OAuth: %w", err)
	}

	return &oauth, nil
}

func (r *OAuthRepository) Create(oauth *model.OAuthProvider) error {
	query := `INSERT INTO oauth_providers 
              (student_id, provider_name, external_id, email) 
              VALUES ($1, $2, $3, $4) 
              RETURNING id, created_at`

	return r.db.QueryRow(
		query, oauth.StudentID, oauth.ProviderName,
		oauth.ExternalID, oauth.Email,
	).Scan(&oauth.ID, &oauth.CreatedAt)
}

func (r *OAuthRepository) FindByStudentAndProvider(studentID int, provider string) (*model.OAuthProvider, error) {
	query := `SELECT id, student_id, provider_name, external_id, email, created_at 
              FROM oauth_providers 
              WHERE student_id = $1 AND provider_name = $2`

	var oauth model.OAuthProvider
	err := r.db.QueryRow(query, studentID, provider).Scan(
		&oauth.ID, &oauth.StudentID, &oauth.ProviderName,
		&oauth.ExternalID, &oauth.Email, &oauth.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска OAuth: %w", err)
	}

	return &oauth, nil
}

package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type CertificateRepository struct {
	db *sql.DB
}

func NewCertificateRepository(db *sql.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

func (r *CertificateRepository) FindByUserID(userID uuid.UUID) ([]models.Certificate, error) {
	query := `
		SELECT id, user_id, pathway_id, title, description,
			certificate_number, verification_url, download_url,
			is_verified, verified_at, issued_date, created_at
		FROM certificates
		WHERE user_id = $1
		ORDER BY issued_date DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query certificates: %w", err)
	}
	defer rows.Close()

	var certificates []models.Certificate
	for rows.Next() {
		var c models.Certificate
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PathwayID,
			&c.Title,
			&c.Description,
			&c.CertificateNumber,
			&c.VerificationURL,
			&c.DownloadURL,
			&c.IsVerified,
			&c.VerifiedAt,
			&c.IssuedDate,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan certificate: %w", err)
		}
		certificates = append(certificates, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return certificates, nil
}

func (r *CertificateRepository) FindByID(id uuid.UUID) (*models.Certificate, error) {
	query := `
		SELECT id, user_id, pathway_id, title, description,
			certificate_number, verification_url, download_url,
			is_verified, verified_at, issued_date, created_at
		FROM certificates
		WHERE id = $1
	`
	cert := &models.Certificate{}
	err := r.db.QueryRow(query, id).Scan(
		&cert.ID,
		&cert.UserID,
		&cert.PathwayID,
		&cert.Title,
		&cert.Description,
		&cert.CertificateNumber,
		&cert.VerificationURL,
		&cert.DownloadURL,
		&cert.IsVerified,
		&cert.VerifiedAt,
		&cert.IssuedDate,
		&cert.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find certificate by ID: %w", err)
	}
	return cert, nil
}

func (r *CertificateRepository) Create(certificate *models.Certificate) error {
	query := `
		INSERT INTO certificates (
			id, user_id, pathway_id, title, description,
			certificate_number, verification_url, download_url,
			is_verified, verified_at, issued_date, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`
	if certificate.ID == uuid.Nil {
		certificate.ID = uuid.New()
	}
	now := time.Now()
	err := r.db.QueryRow(
		query,
		certificate.ID,
		certificate.UserID,
		certificate.PathwayID,
		certificate.Title,
		certificate.Description,
		certificate.CertificateNumber,
		certificate.VerificationURL,
		certificate.DownloadURL,
		certificate.IsVerified,
		certificate.VerifiedAt,
		certificate.IssuedDate,
		now,
	).Scan(&certificate.ID, &certificate.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}
	return nil
}

func (r *CertificateRepository) Update(certificate *models.Certificate) error {
	query := `
		UPDATE certificates
		SET
			title = $2,
			description = $3,
			certificate_number = $4,
			verification_url = $5,
			download_url = $6,
			is_verified = $7,
			verified_at = $8,
			issued_date = $9
		WHERE id = $1
		RETURNING id
	`
	_, err := r.db.Exec(
		query,
		certificate.ID,
		certificate.Title,
		certificate.Description,
		certificate.CertificateNumber,
		certificate.VerificationURL,
		certificate.DownloadURL,
		certificate.IsVerified,
		certificate.VerifiedAt,
		certificate.IssuedDate,
	)
	if err != nil {
		return fmt.Errorf("failed to update certificate: %w", err)
	}
	return nil
}

func (r *CertificateRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM certificates WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("certificate not found")
	}
	return nil
}
package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	UserID             uuid.UUID  `json:"user_id" db:"user_id"`
	PathwayID          uuid.UUID  `json:"pathway_id" db:"pathway_id"`
	Title              string     `json:"title" db:"title"`
	Description        *string    `json:"description,omitempty" db:"description"`
	CertificateNumber  string     `json:"certificate_number" db:"certificate_number"`
	VerificationURL    *string    `json:"verification_url,omitempty" db:"verification_url"`
	DownloadURL        *string    `json:"download_url,omitempty" db:"download_url"`
	IsVerified         bool       `json:"is_verified" db:"is_verified"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	IssuedDate         time.Time  `json:"issued_date" db:"issued_date"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}
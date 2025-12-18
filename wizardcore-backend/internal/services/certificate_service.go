package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type CertificateService struct {
	certificateRepo *repositories.CertificateRepository
}

func NewCertificateService(certificateRepo *repositories.CertificateRepository) *CertificateService {
	return &CertificateService{certificateRepo: certificateRepo}
}

func (s *CertificateService) GetUserCertificates(userID uuid.UUID) ([]models.Certificate, error) {
	return s.certificateRepo.FindByUserID(userID)
}

func (s *CertificateService) GetCertificateByID(id uuid.UUID) (*models.Certificate, error) {
	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}
	if cert == nil {
		return nil, fmt.Errorf("certificate not found")
	}
	return cert, nil
}

func (s *CertificateService) CreateCertificate(certificate *models.Certificate) error {
	if certificate.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if certificate.PathwayID == uuid.Nil {
		return fmt.Errorf("pathway_id is required")
	}
	if certificate.Title == "" {
		return fmt.Errorf("title is required")
	}
	if certificate.CertificateNumber == "" {
		return fmt.Errorf("certificate_number is required")
	}
	if certificate.IssuedDate.IsZero() {
		return fmt.Errorf("issued_date is required")
	}
	return s.certificateRepo.Create(certificate)
}

func (s *CertificateService) UpdateCertificate(certificate *models.Certificate) error {
	if certificate.ID == uuid.Nil {
		return fmt.Errorf("certificate ID is required")
	}
	return s.certificateRepo.Update(certificate)
}

func (s *CertificateService) DeleteCertificate(id uuid.UUID) error {
	return s.certificateRepo.Delete(id)
}
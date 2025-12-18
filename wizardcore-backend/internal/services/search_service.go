package services

import (
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type SearchService struct {
	searchRepo *repositories.SearchRepository
}

func NewSearchService(searchRepo *repositories.SearchRepository) *SearchService {
	return &SearchService{searchRepo: searchRepo}
}

func (s *SearchService) Search(query string, limit, offset int) (*models.SearchResults, error) {
	return s.searchRepo.Search(query, limit, offset)
}
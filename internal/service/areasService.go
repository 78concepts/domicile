package service

import (
	"78concepts.com/domicile/internal/model"
	"78concepts.com/domicile/internal/repository"
	"context"
	"github.com/gofrs/uuid"
)

func NewAreasService(areasRepository repository.IAreasRepository) *AreasService {
	return &AreasService{areasRepository: areasRepository}
}

type AreasService struct {
	areasRepository repository.IAreasRepository
}

func (s *AreasService) GetAreas(ctx context.Context) ([]model.Area, error) {
	return s.areasRepository.GetAreas(ctx)
}

func (s *AreasService) GetArea(ctx context.Context, uuid uuid.UUID) (*model.Area, error) {
	return s.areasRepository.GetArea(
		ctx,
		uuid,
	)
}
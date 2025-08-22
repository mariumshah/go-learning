package converter

import (
	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/pkg/models"
)

func AuthorToDomain(m *models.Author) *domain.Author {
	return &domain.Author{
		ID:   int(m.ID),
		Name: m.Name,
	}

}

func AuthorToModel(u *domain.Author) *models.Author {
	return &models.Author{
		ID:   int(u.ID),
		Name: u.Name,
	}
}

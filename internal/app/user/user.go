package converter

import (
	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/pkg/models"
)

func UserToDomain(m *models.User) *domain.User {
	return &domain.User{
		ID:    int(m.ID),
		Email: m.Email,
		Hash:  m.PasswordHash,
	}
}

func UserToModel(u *domain.User) *models.User {
	return &models.User{
		ID:           int(u.ID),
		Email:        u.Email,
		PasswordHash: u.Hash,
	}
}

// TODO TEST: Write a test converting back-and-forth; assert equality.

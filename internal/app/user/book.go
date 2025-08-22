package converter

import (
	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/pkg/models"
)

func BookToDomain(m *models.Book) *domain.Book {
	if m == nil {
		return nil
	}
	return &domain.Book{
		ID:     int(m.ID),
		Title:  m.Title,
		Author: int(m.AuthorID), // caller can populate author name if needed (see below)
		// ISBN:    m.ISBN.String, // if sqlboiler produced sql.NullString; adjust type if needed
		// AddedAt: m.CreatedAt, // depending on your generated fields
	}
}

func BookToModel(u *domain.Book) *models.Book {
	if u == nil {
		return nil
	}
	return &models.Book{
		ID:       int(u.ID),
		Title:    u.Title,
		AuthorID: int(u.Author),
		// ISBN:    sql.NullString{String: u.ISBN, Valid: u.ISBN != ""}, // adjust type if needed
		// AddedAt: u.AddedAt, // depending on your generated fields
	}
}

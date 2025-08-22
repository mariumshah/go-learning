// package repository

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"

// 	"github.com/aarondl/sqlboiler/v4/boil"
// 	"github.com/aarondl/sqlboiler/v4/queries/qm"

// 	"github.com/playground/userapi/internal/app/domain"
// 	converter "github.com/playground/userapi/internal/app/user"
// 	"github.com/playground/userapi/pkg/models"
// )

// // LibraryRepo manages adding/listing/removing books for user libraries.
// type LibraryRepo interface {
// 	AddAuthor (ctx context.Context, name string) (int, error)
// 	FindOrCreateAuthorByName(ctx context.Context, name string) (*domain.Author, error)

// 	AddBook(ctx context.Context, userID int, b *domain.Book) (*domain.Book, error)
// 	ListBooks(ctx context.Context, userID int) ([]*domain.Book, error)
// 	RemoveBook(ctx context.Context, userID, bookID int) error
// }

// type libraryRepo struct{ db boil.ContextExecutor }

// func NewLibraryRepo(db boil.ContextExecutor) LibraryRepo {
// 	return &libraryRepo{db: db}
// }

// func (r *libraryRepo) AddBook(ctx context.Context, userID int, b *domain.Book) (*domain.Book, error) {
// 	// 1) Ensure author exists (require valid author id)
// 	author, err := models.Authors(models.AuthorWhere.ID.EQ(b.Author)).One(ctx, r.db)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, fmt.Errorf("author with id %d not found", b.Author)
// 		}
// 		return nil, fmt.Errorf("find author: %w", err)
// 	}

// 	// 2) Find existing book by title+author(+isbn if provided)
// 	var book *models.Book
// 	book, err = models.Books(qm.Where("title = ? AND author_id = ?", b.Title, author.ID)).One(ctx, r.db)

// 	// if b.ISBN != "" {
// 	// 	book, err = models.Books(qm.Where("title = ? AND author_id = ? AND isbn = ?", b.Title, author.ID, b.ISBN)).One(ctx, r.db)
// 	// } else {
// 	// 	book, err = models.Books(qm.Where("title = ? AND author_id = ?", b.Title, author.ID)).One(ctx, r.db)
// 	// }
// 	if err != nil && !errors.Is(err, sql.ErrNoRows) {
// 		return nil, fmt.Errorf("find book: %w", err)
// 	}
// 	if errors.Is(err, sql.ErrNoRows) {
// 		book = &models.Book{
// 			Title:    b.Title,
// 			AuthorID: author.ID,
// 		}
// 		// optional fields if present on your model
// 		// if b.ISBN != "" {
// 		// 	book.ISBN = b.ISBN
// 		// }
// 		// if b.PublicationYear != 0 {
// 		// 	// adjust field name if your model uses a different one
// 		// 	book.PublicationYear = b.PublicationYear
// 		// }

// 		if err := book.Insert(ctx, r.db, boil.Infer()); err != nil {
// 			return nil, fmt.Errorf("insert book: %w", err)
// 		}
// 	}

// 	// 3) Insert into user_books (junction) — will fail if already present
// 	ub := &models.UserBook{
// 		UserID: userID,
// 		BookID: book.ID,
// 	}
// 	if err := ub.Insert(ctx, r.db, boil.Infer()); err != nil {
// 		// If duplicate (user already has it), be idempotent: return the existing book
// 		// A better approach is to detect duplicate-key error from the driver; for now return book.
// 		return converter.BookToDomain(book), nil
// 	}

// 	// 4) Update total_libraries counter (best-effort)
// 	// Try to do an atomic update if underlying db supports ExecContext
// 	if dbsql, ok := r.db.(*sql.DB); ok {
// 		if _, err := dbsql.ExecContext(ctx, "UPDATE books SET total_libraries = total_libraries + 1 WHERE id = ?", book.ID); err != nil {
// 			// not fatal, proceed
// 		}
// 	} else {
// 		// best-effort fallback: increment and update via model
// 		// TotalLibraries is a nullable-type with fields Int and Valid — handle accordingly.
// 		// If your generated type is slightly different adjust these field names.
// 		if book.TotalLibraries.Valid {
// 			book.TotalLibraries.Int = book.TotalLibraries.Int + 1
// 			book.TotalLibraries.Valid = true
// 		} else {
// 			book.TotalLibraries.Int = 1
// 			book.TotalLibraries.Valid = true
// 		}
// 		_, _ = book.Update(ctx, r.db, boil.Infer())
// 	}

// 	// return domain model
// 	return converter.BookToDomain(book), nil
// }

// func (r *libraryRepo) ListBooks(ctx context.Context, userID int) ([]*domain.Book, error) {
// 	// Join books with user_books to get user's books ordered by added_at desc
// 	ms, err := models.Books(
// 		qm.InnerJoin("user_books ub on ub.book_id = books.id"),
// 		qm.Where("ub.user_id = ?", userID),
// 		qm.OrderBy("ub.added_at DESC"),
// 	).All(ctx, r.db)
// 	if err != nil {
// 		return nil, fmt.Errorf("list books: %w", err)
// 	}
// 	out := make([]*domain.Book, 0, len(ms))
// 	for _, m := range ms {
// 		d := converter.BookToDomain(m)
// 		// Try to load author id (or name) if needed
// 		if a, err := models.Authors(models.AuthorWhere.ID.EQ(m.AuthorID)).One(ctx, r.db); err == nil {
// 			d.Author = a.ID
// 		}
// 		out = append(out, d)
// 	}
// 	return out, nil
// }

// func (r *libraryRepo) RemoveBook(ctx context.Context, userID, bookID int) error {
// 	// Delete row from user_books where both match
// 	count, err := models.UserBooks(
// 		models.UserBookWhere.UserID.EQ(userID),
// 		models.UserBookWhere.BookID.EQ(bookID),
// 	).DeleteAll(ctx, r.db)
// 	if err != nil {
// 		return fmt.Errorf("delete user_book: %w", err)
// 	}
// 	if count == 0 {
// 		return sql.ErrNoRows
// 	}

//		// Decrement book.total_libraries best-effort
//		if dbsql, ok := r.db.(*sql.DB); ok {
//			_, _ = dbsql.ExecContext(ctx, "UPDATE books SET total_libraries = total_libraries - 1 WHERE id = ?", bookID)
//		} else {
//			// fallback: try load book and decrement its nullable counter
//			if book, err := models.FindBook(ctx, r.db, bookID); err == nil {
//				if book.TotalLibraries.Valid && book.TotalLibraries.Int > 0 {
//					book.TotalLibraries.Int = book.TotalLibraries.Int - 1
//					book.TotalLibraries.Valid = true
//				} else {
//					// keep it zero/invalid
//					book.TotalLibraries.Int = 0
//					book.TotalLibraries.Valid = false
//				}
//				_, _ = book.Update(ctx, r.db, boil.Infer())
//			}
//		}
//		return nil
//	}
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"github.com/playground/userapi/internal/app/domain"
	converter "github.com/playground/userapi/internal/app/user"
	"github.com/playground/userapi/pkg/models"
)

// LibraryRepo manages adding/listing/removing books for user libraries.
type LibraryRepo interface {
	// Author operations
	AddAuthor(ctx context.Context, name string) (int, error)
	FindOrCreateAuthorByName(ctx context.Context, name string) (*domain.Author, error)

	// Book/library operations
	AddBook(ctx context.Context, userID int, b *domain.Book) (*domain.Book, error)
	ListBooks(ctx context.Context, userID int) ([]*domain.Book, error)
	RemoveBook(ctx context.Context, userID, bookID int) error
	ListAllBooks(ctx context.Context, userID int) ([]*domain.Book, error)
}

type libraryRepo struct {
	db    boil.ContextExecutor
	sqldb *sql.DB
}

func NewLibraryRepo(db boil.ContextExecutor) LibraryRepo {
	return &libraryRepo{db: db}
}

// AddAuthor inserts a new author (if not exists) and returns its id.
// Uses converter.AuthorToModel / AuthorToDomain to convert between domain and model.
func (r *libraryRepo) AddAuthor(ctx context.Context, name string) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("author name required")
	}

	// try find existing
	existing, err := models.Authors(models.AuthorWhere.Name.EQ(name)).One(ctx, r.db)
	if err == nil {
		return existing.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("find author: %w", err)
	}

	// create via domain->model converter so we use a consistent conversion path
	dom := &domain.Author{Name: name}
	na := converter.AuthorToModel(dom)

	if err := na.Insert(ctx, r.db, boil.Infer()); err != nil {
		// possible race: another request inserted same name concurrently.
		// Try to re-query before returning a hard error.
		if a2, err2 := models.Authors(models.AuthorWhere.Name.EQ(name)).One(ctx, r.db); err2 == nil {
			return a2.ID, nil
		}
		return 0, fmt.Errorf("insert author: %w", err)
	}

	return na.ID, nil
}

// FindOrCreateAuthorByName returns existing author or inserts and returns a domain.Author.
// Uses converters so callers always get domain types.
func (r *libraryRepo) FindOrCreateAuthorByName(ctx context.Context, name string) (*domain.Author, error) {
	if name == "" {
		return nil, fmt.Errorf("author name required")
	}

	a, err := models.Authors(models.AuthorWhere.Name.EQ(name)).One(ctx, r.db)
	if err == nil {
		return converter.AuthorToDomain(a), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("find author: %w", err)
	}

	// insert new author via converter
	dom := &domain.Author{Name: name}
	na := converter.AuthorToModel(dom)
	if err := na.Insert(ctx, r.db, boil.Infer()); err != nil {
		// race handling: try re-query
		if a2, err2 := models.Authors(models.AuthorWhere.Name.EQ(name)).One(ctx, r.db); err2 == nil {
			return converter.AuthorToDomain(a2), nil
		}
		return nil, fmt.Errorf("insert author: %w", err)
	}
	return converter.AuthorToDomain(na), nil
}

func (r *libraryRepo) AddBook(ctx context.Context, userID int, b *domain.Book) (*domain.Book, error) {
	// requires that the author already exists (two-step flow):
	// client calls POST /authors -> receives author_id, then POST /library/books with author_id
	author, err := models.Authors(models.AuthorWhere.ID.EQ(b.Author)).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("author with id %d not found", b.Author)
		}
		return nil, fmt.Errorf("find author: %w", err)
	}

	// 2) Find existing book by title+author (+isbn if you want)
	var book *models.Book
	mods := []qm.QueryMod{
		models.BookWhere.Title.EQ(b.Title),
		models.BookWhere.AuthorID.EQ(b.Author),
	}
	book, err = models.Books(mods...).One(ctx, r.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("find book: %w", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		book = &models.Book{
			Title:    b.Title,
			AuthorID: author.ID,
		}
		// if b.ISBN != "" {
		// 	// adjust if your model's field name differs
		// 	book.ISBN = b.ISBN
		// }
		if err := book.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("insert book: %w", err)
		}
	}

	// 3) Insert into user_books (junction)
	ub := &models.UserBook{
		UserID: userID,
		BookID: book.ID,
	}

	txn, err := r.sqldb.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := ub.Insert(ctx, txn, boil.Infer()); err != nil {
		// idempotent: if duplicate (user already has it) return existing domain book
		return converter.BookToDomain(book), nil
	}

	// 4) Update total_libraries counter (best-effort)
	if dbsql, ok := r.db.(*sql.DB); ok {
		_, _ = dbsql.ExecContext(ctx, "UPDATE books SET total_libraries = total_libraries + 1 WHERE id = ?", book.ID)
	} else {
		// model fallback (nullable handling)
		if book.TotalLibraries.Valid {
			book.TotalLibraries.Int = book.TotalLibraries.Int + 1
			book.TotalLibraries.Valid = true
		} else {
			book.TotalLibraries.Int = 1
			book.TotalLibraries.Valid = true
		}
		_, _ = book.Update(ctx, r.db, boil.Infer())
	}

	// convert models.Book -> domain.Book via your existing converter (keeps consistency)
	return converter.BookToDomain(book), nil
}

func (r *libraryRepo) ListBooks(ctx context.Context, userID int) ([]*domain.Book, error) {
	ms, err := models.Books(
		qm.InnerJoin("user_books ub on ub.book_id = books.id"),
		qm.Where("ub.user_id = ?", userID),
		qm.OrderBy("ub.added_at DESC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	out := make([]*domain.Book, 0, len(ms))
	for _, m := range ms {
		d := converter.BookToDomain(m)
		// if you want author as domain object here, convert it:
		if a, err := models.Authors(models.AuthorWhere.ID.EQ(m.AuthorID)).One(ctx, r.db); err == nil {
			// either attach id (existing behavior) or attach full author domain
			d.Author = a.ID
			// optional: expose AuthorName if domain.Book has that field:
			// d.AuthorName = a.Name
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *libraryRepo) ListAllBooks(ctx context.Context, userID int) ([]*domain.Book, error) {
	ms, err := models.Books(
		qm.Load(models.BookRels.Author),
		qm.OrderBy("books.created_at DESC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}
	out := make([]*domain.Book, 0, len(ms))
	for _, m := range ms {
		d := converter.BookToDomain(m)

		if a, err := models.Authors(models.AuthorWhere.ID.EQ(m.AuthorID)).One(ctx, r.db); err == nil {
			d.Author = a.ID
			// optional: expose AuthorName if domain.Book has that field:
			// d.AuthorName = a.Name
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *libraryRepo) RemoveBook(ctx context.Context, userID, bookID int) error {
	count, err := models.UserBooks(
		models.UserBookWhere.UserID.EQ(userID),
		models.UserBookWhere.BookID.EQ(bookID),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return fmt.Errorf("delete user_book: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	if dbsql, ok := r.db.(*sql.DB); ok {
		_, _ = dbsql.ExecContext(ctx, "UPDATE books SET total_libraries = total_libraries - 1 WHERE id = ?", bookID)
	} else {
		if book, err := models.FindBook(ctx, r.db, bookID); err == nil {
			if book.TotalLibraries.Valid && book.TotalLibraries.Int > 0 {
				book.TotalLibraries.Int = book.TotalLibraries.Int - 1
				book.TotalLibraries.Valid = true
			} else {
				book.TotalLibraries.Int = 0
				book.TotalLibraries.Valid = false
			}
			_, _ = book.Update(ctx, r.db, boil.Infer())
		}
	}
	return nil
}

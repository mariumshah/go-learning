// package repository

// import (
// 	"context"
// 	"errors"

// 	"github.com/aarondl/sqlboiler/v4/boil"
// 	"github.com/playground/userapi/internal/app/domain"
// 	converter "github.com/playground/userapi/internal/app/user"
// 	"github.com/playground/userapi/pkg/models"
// )
// var ErrEmailExists = errors.New("email already exists")

// type UserRepo interface {
// 	Create(ctx context.Context, u *domain.User) error
// 	GetByEmail(ctx context.Context, email string) (*domain.User, error)
// 	GetByEmailForUpdate(ctx context.Context, email string) (*domain.User, error) // SELECT ... FOR UPDATE
// 	Update(ctx context.Context, u *domain.User) error
// }
// type userRepo struct{ db boil.ContextExecutor }

// func NewUserRepo(db boil.ContextExecutor) UserRepo { return &userRepo{db} }
// func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
// 	m := converter.UserToModel(u)
// 	return m.Insert(ctx, r.db, boil.Infer())
// }

// // GetByEmail uses qm.Where for filtering
// func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
// 	m, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, r.db)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return converter.UserToDomain(m), nil
// }
// func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
// 	m := converter.UserToModel(u)
// 	_, err := m.Update(ctx, r.db, boil.Infer())
// 	return err
// }

package repository

import (
	"context"
	// "database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	mysql "github.com/go-sql-driver/mysql"
	"github.com/playground/userapi/internal/app/domain"
	converter "github.com/playground/userapi/internal/app/user"
	"github.com/playground/userapi/pkg/models"
)

// ErrEmailExists is returned when attempting to create a user with a duplicate email.
var ErrEmailExists = errors.New("email already exists")

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByEmailForUpdate(ctx context.Context, email string) (*domain.User, error) // SELECT ... FOR UPDATE
	Update(ctx context.Context, u *domain.User) error
}

type userRepo struct{ db boil.ContextExecutor }

func NewUserRepo(db boil.ContextExecutor) UserRepo { return &userRepo{db} }

// Create inserts a new user. It will:
//   - check (fast) for an existing email to return a friendly ErrEmailExists,
//   - attempt the insert and copy the DB-assigned ID back to u.ID,
//   - detect duplicate-key DB errors as a fallback for race conditions.
func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	// quick uniqueness check (helps give a nicer error)
	if existing, _ := r.GetByEmail(ctx, u.Email); existing != nil {
		return ErrEmailExists
	}

	m := converter.UserToModel(u) // leave ID zero for DB auto-increment

	if err := m.Insert(ctx, r.db, boil.Infer()); err != nil {
		// detect MySQL duplicate-key error (race-safety)
		if isMySQLDuplicateErr(err) {
			return ErrEmailExists
		}
		return fmt.Errorf("insert user: %w", err)
	}

	// copy assigned ID back into domain object so caller sees it
	u.ID = m.ID
	return nil
}

// GetByEmail returns domain user or sql.ErrNoRows if not found
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	m, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, r.db)
	if err != nil {
		// keep underlying sql.ErrNoRows behavior to let callers distinguish "not found"
		return nil, err
	}
	return converter.UserToDomain(m), nil
}

// GetByEmailForUpdate demonstrates a SELECT ... FOR UPDATE using qm.For("UPDATE").
// Useful when you need to lock the row before making read-modify-write changes in a transaction.
func (r *userRepo) GetByEmailForUpdate(ctx context.Context, email string) (*domain.User, error) {
	m, err := models.Users(qm.Where("email = ?", email), qm.For("UPDATE")).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return converter.UserToDomain(m), nil
}

// Update persists changes. For SQLBoiler Update() we convert domain->model and run Update.
func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	// require valid ID for update
	if u.ID == 0 {
		return errors.New("missing id for update")
	}

	m := converter.UserToModel(u)
	_, err := m.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// isMySQLDuplicateErr inspects an error and returns true for MySQL duplicate-key errors.
// If you later support Postgres, add a Postgres branch (pq.Error with Code "23505") here.
func isMySQLDuplicateErr(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		// MySQL duplicate entry error number is 1062
		return me.Number == 1062
	}
	// For other drivers, you can inspect error strings or error types.
	return false
}

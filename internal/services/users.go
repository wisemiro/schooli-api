package services

import (
	"context"
	"errors"
	"log"
	"schooli-api/internal/models"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/resterrors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type UserService interface {
	CreateUser(ctx context.Context, um models.User) (*models.User, error)
	UpdateUser(ctx context.Context, um models.User) error
	DeleteUser(ctx context.Context, um models.User) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	OneUser(ctx context.Context, userID int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

func (sq *SQLStore) CreateUser(ctx context.Context, um models.User) (*models.User, error) {
	if err := sq.store.CreateUser(
		ctx,
		db.CreateUserParams{
			Email:        um.Email,
			PhoneNumber:  um.PhoneNumber,
			PasswordHash: um.PasswordHash,
		},
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "users_email_key" {
				log.Println(*pgErr)
				return nil, errors.New(resterrors.ErrorEmailExists)
			} else if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "users_phone_number_key" {
				log.Println(pgErr.ColumnName)
				return nil, errors.New(resterrors.ErrorPhoneExists)
			}
		}
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.Create")
	}
	return nil, nil
}
func (sq *SQLStore) UpdateUser(ctx context.Context, um models.User) error {
	if err := sq.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:          um.ID,
		Email:       um.Email,
		PhoneNumber: um.PhoneNumber,
	},
	); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.Update")
	}
	return nil
}
func (sq *SQLStore) DeleteUser(ctx context.Context, um models.User) error {
	if err := sq.store.DeleteUser(ctx,
		um.ID,
	); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.Delete")
	}
	return nil
}
func (sq *SQLStore) ListUsers(ctx context.Context) ([]*models.User, error) {
	users, err := sq.store.ListUsers(ctx)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.List")

	}
	us := make([]*models.User, len(users))
	for i, user := range users {
		us[i] = &models.User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt.Time,
			UpdatedAt:   user.UpdatedAt.Time,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			IsVerified:  user.IsVerified.Bool,
		}
	}
	return us, nil
}
func (sq *SQLStore) OneUser(ctx context.Context, userID int64) (*models.User, error) {
	user, err := sq.store.GetUser(ctx, userID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.One")
	}
	return &models.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified.Bool,
	}, nil
}

func (sq *SQLStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := sq.store.UserByEmail(ctx, email)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "UsersService.FindByEmail")
	}
	return &models.User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		IsVerified:   user.IsVerified.Bool,
		PasswordHash: user.PasswordHash,
	}, nil
}

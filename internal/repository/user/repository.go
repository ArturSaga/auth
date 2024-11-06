package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/auth/internal/client/db"
	"github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/repository"
	"github.com/ArturSaga/auth/internal/repository/user/converter"
	modelRepo "github.com/ArturSaga/auth/internal/repository/user/model"
	serviceErr "github.com/ArturSaga/auth/internal/service_error"
)

const (
	hashCost = 10

	tableName          = "users"
	idColumn           = "id"
	nameColumn         = "name"
	emailColumn        = "email"
	passwordHashColumn = "password_hash"
	roleColumn         = "role"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

type repo struct {
	db db.Client
}

// NewUserRepository - публичный метод, создащий сущность репозитория, для работы с данными сущности в бд
func NewUserRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// CreateUser - публичный метод, создания пользователя в бд
func (r *repo) CreateUser(ctx context.Context, userInfo *model.UserInfo) (int64, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), hashCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v", err)
		return 0, err
	}

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordHashColumn, roleColumn).
		Values(userInfo.Name, userInfo.Email, hashPassword, userInfo.Role.String()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		fmt.Printf("failed to insert user: %v", err)
		return 0, err
	}

	log.Printf("inserted user with id: %d", userID)

	return userID, nil
}

// GetUser - публичный метод, получения пользователя из бд
func (r *repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, passwordHashColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, err
	}
	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		fmt.Printf("failed to get user: %v", err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser - публичный метод, обновления пользователя в бд
func (r *repo) UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) (emptypb.Empty, error) {
	builderUserInfoUpdate := sq.Update(tableName)
	hasUpdates := false

	if updateUserInfo.Name != nil {
		builderUserInfoUpdate = builderUserInfoUpdate.Set(nameColumn, *updateUserInfo.Name)
		hasUpdates = true
	}

	if updateUserInfo.Role != nil {
		roleStr := fmt.Sprintf("%s", *updateUserInfo.Role)
		builderUserInfoUpdate = builderUserInfoUpdate.Set(roleColumn, roleStr)
		hasUpdates = true
	}

	if updateUserInfo.Password != nil {
		hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(*updateUserInfo.Password), hashCost)
		if err != nil {
			fmt.Printf("failed to hash password: %v", err)
			return emptypb.Empty{}, err
		}
		builderUserInfoUpdate = builderUserInfoUpdate.Set(passwordHashColumn, hashNewPassword)
		hasUpdates = true
	}

	if hasUpdates {
		query, args, err := builderUserInfoUpdate.PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": updateUserInfo.UserID}).ToSql()
		if err != nil {
			return emptypb.Empty{}, fmt.Errorf("failed to build query: %v", err)
		}

		q := db.Query{
			Name:     "user_repository.Update",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			return emptypb.Empty{}, err
		}
	} else {
		return emptypb.Empty{}, serviceErr.ErrUpdateUser
	}

	return emptypb.Empty{}, nil
}

// DeleteUser - публичный метод, удаления пользователя в бд
func (r *repo) DeleteUser(ctx context.Context, id int64) (emptypb.Empty, error) {
	builderUserDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderUserDelete.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return emptypb.Empty{}, err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		fmt.Printf("failed to delete roleID: %v", err)
	}
	return emptypb.Empty{}, nil
}

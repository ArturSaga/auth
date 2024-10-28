package user

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ArturSaga/auth/internal/client/db"
	"github.com/ArturSaga/auth/internal/model"
	"github.com/ArturSaga/auth/internal/repository"
	"github.com/ArturSaga/auth/internal/repository/user/converter"
	modelRepo "github.com/ArturSaga/auth/internal/repository/user/model"
)

const hashCost = 10

const tableName = "users"

const idColumn = "id"
const nameColumn = "name"
const emailColumn = "email"
const passwordHashColumn = "password_hash"
const roleColumn = "role"
const createdAtColumn = "created_at"
const updatedAtColumn = "updated_at"

type repo struct {
	db db.Client
}

func NewUserRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, userInfo *model.UserInfo) (int64, error) {
	if userInfo.Password != userInfo.PasswordConfirm {
		return 0, errors.New("confirm password not equal to password")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), hashCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v", err)
		return 0, err
	}

	now := time.Now()
	fmt.Println(userInfo.Role)
	// Делаем запрос на вставку записи в таблицу user
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordHashColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(userInfo.Name, userInfo.Email, hashPassword, userInfo.Role.String(), now, now).
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
	var userInfo modelRepo.UserInfo
	var role string
	var createdAt, updatedAt time.Time
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &userInfo.Name, &userInfo.Email, &userInfo.Password, &role, &createdAt, &updatedAt)
	if err != nil {
		fmt.Printf("failed to get user: %v", err)
		return nil, err
	}
	userInfo.Role = modelRepo.RoleFromString(role)
	user.Info = userInfo
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) (emptypb.Empty, error) {
	user, err := r.GetUser(ctx, updateUserInfo.UserID)
	fmt.Println(updateUserInfo.Role)
	if err != nil {
		return emptypb.Empty{}, err
	}

	if updateUserInfo.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Info.Password), []byte(updateUserInfo.OldPassword.Value))
		if err != nil {
			return emptypb.Empty{}, fmt.Errorf("old password not equal to current password: %v", err)
		}

		if updateUserInfo.Password != updateUserInfo.PasswordConfirm {
			return emptypb.Empty{}, errors.New("confirm password not equal to password")
		}
	}

	builderUserInfoUpdate := sq.Update(tableName)
	hasUpdates := false // Флаг для проверки наличия обновлений

	if updateUserInfo.Name != nil {
		builderUserInfoUpdate = builderUserInfoUpdate.Set("name", updateUserInfo.Name.Value)
		hasUpdates = true
	}

	fmt.Println(updateUserInfo.Role)
	fmt.Println(user.Info.Role)

	if updateUserInfo.Role != user.Info.Role {
		builderUserInfoUpdate = builderUserInfoUpdate.Set("role", updateUserInfo.Role.String())
		hasUpdates = true
	}

	if updateUserInfo.Password != nil {
		hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserInfo.Password.Value), hashCost)
		if err != nil {
			fmt.Printf("failed to hash password: %v", err)
			return emptypb.Empty{}, err
		}
		builderUserInfoUpdate = builderUserInfoUpdate.Set("password", hashNewPassword)
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
			return emptypb.Empty{}, fmt.Errorf("failed to update user: %v", err)
		}
	} else {
		return emptypb.Empty{}, errors.New("No fields to update")
	}

	return emptypb.Empty{}, nil
}

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
		fmt.Printf("failed to get roleID: %v", err)
	}
	return emptypb.Empty{}, nil
}

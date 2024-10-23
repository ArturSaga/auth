package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ArturSaga/auth/api/grpc/pkg/user_v1"
	"github.com/ArturSaga/auth/internal/config"
	"github.com/ArturSaga/auth/internal/config/env"
)

const hashCost = 10

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserApiServer
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserApiServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser - публичный метод, который создает пользователя.
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	if req.Info.Password != req.Info.PasswordConfirm {
		return nil, errors.New("confirm password not equal to password")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Info.Password), hashCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v", err)
		return nil, err
	}

	hashConfirmPassword, err := bcrypt.GenerateFromPassword([]byte(req.Info.PasswordConfirm), hashCost)
	if err != nil {
		fmt.Printf("failed to hash confirm password: %v", err)
		return nil, err
	}

	roleID, err := s.getRoleID(ctx, req.Info.Role.String())
	if err != nil {
		fmt.Printf("failed to get RoleID: %v", err)
		return nil, err
	}

	// Делаем запрос на вставку записи в таблицу note
	builderInsert := sq.Insert("user_info").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role_id").
		Values(req.Info.Name, req.Info.Email, hashPassword, hashConfirmPassword, roleID).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, err
	}

	var userInfoID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userInfoID)
	if err != nil {
		fmt.Printf("failed to insert user_info: %v", err)
		return nil, err
	}

	// Делаем запрос на вставку записи в таблицу note
	builderUserInsert := sq.Insert("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Columns("info_id", "created_at", "updated_at").
		Values(userInfoID, time.Now(), time.Now()).
		Suffix("RETURNING id")

	query, args, err = builderUserInsert.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		fmt.Printf("failed to insert user: %v", err)
		return nil, err
	}

	log.Printf("inserted note with id: %d", userID)

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUser - публичный метод, который позволяет получить данные пользователя.
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	builderUserSelect := sq.Select("id", "info_id", "created_at", "updated_at").
		From("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderUserSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var user desc.User
	var userInfoID int64
	var createdAt, updatedAt time.Time
	err = s.pool.QueryRow(ctx, query, args...).Scan(&user.Id, &userInfoID, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	userInfo, err := s.getUserInfo(ctx, userInfoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_info: %v", err)
	}

	return &desc.GetUserResponse{
		User: &desc.User{
			Id: user.Id,
			Info: &desc.UserInfo{
				Name:            userInfo.Name,
				Email:           userInfo.Email,
				Password:        userInfo.Password,
				PasswordConfirm: userInfo.PasswordConfirm,
				Role:            desc.Role_USER,
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// UpdateUser - публичный метод, который обновляет данные пользователя.
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	user, userInfoID := s.getUser(ctx, req.Info.UserID)
	if user == nil {
		return nil, errors.New("User not found by this ID")
	}
	userInfo, err := s.getUserInfo(ctx, userInfoID)
	if err != nil {
		return nil, errors.New("User info not found for this user")
	}

	if req.Info.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.Info.OldPassword.Value))
		if err != nil {
			return nil, fmt.Errorf("old password not equal to current password: %v", err)
		}

		if req.Info.Password.Value != req.Info.PasswordConfirm.Value {
			return nil, errors.New("confirm password not equal to password")
		}
	}

	builderUserInfoUpdate := sq.Update("user_info")
	hasUpdates := false // Флаг для проверки наличия обновлений

	if req.Info.Name != nil {
		builderUserInfoUpdate = builderUserInfoUpdate.Set("name", req.Info.Name.Value)
		hasUpdates = true
	}

	if req.Info.Role != desc.Role_UNKNOWN {
		roleID, err := s.getRoleID(ctx, req.Info.Role.String())
		if err != nil {
			fmt.Printf("failed to get roleID: %v", err)
			return nil, err
		}
		builderUserInfoUpdate = builderUserInfoUpdate.Set("role_id", roleID)
		hasUpdates = true
	}

	if req.Info.Password != nil {
		hashNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.Info.Password.Value), hashCost)
		if err != nil {
			fmt.Printf("failed to hash password: %v", err)
			return nil, err
		}
		builderUserInfoUpdate = builderUserInfoUpdate.Set("password", hashNewPassword)
		hasUpdates = true
	}

	if req.Info.PasswordConfirm != nil {
		hashConfirmPassword, err := bcrypt.GenerateFromPassword([]byte(req.Info.PasswordConfirm.Value), hashCost)
		if err != nil {
			fmt.Printf("failed to hash password: %v", err)
			return nil, err
		}
		builderUserInfoUpdate = builderUserInfoUpdate.Set("password_confirm", hashConfirmPassword)
		hasUpdates = true
	}

	if hasUpdates {
		query, args, err := builderUserInfoUpdate.PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": userInfoID}).ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query: %v", err)
		}

		_, err = s.pool.Exec(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %v", err)
		}

		builderUserUpdate := sq.Update("\"user\"").
			PlaceholderFormat(sq.Dollar).
			Set("updated_at", time.Now()).
			Where(sq.Eq{"id": req.Info.UserID})
		query, args, err = builderUserUpdate.ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %v", err)
		}
		_, err = s.pool.Exec(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %v", err)
		}
	} else {
		return nil, errors.New("No fields to update")
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser - публичный метод, который удаляет пользователя.
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, nil
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("failed to get roleID: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *server) getRoleID(ctx context.Context, role string) (int64, error) {
	builderSelect := sq.Select("id").
		From("role").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"name": role}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return 0, err
	}

	var roleID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&roleID)
	if err != nil {
		fmt.Printf("failed to get roleID: %v", err)
		return 0, err
	}

	return roleID, nil
}

func (s *server) getUser(ctx context.Context, userID int64) (*desc.User, int64) {
	builderSelect := sq.Select("id", "info_id", "created_at", "updated_at").
		From("\"user\"").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userID}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		fmt.Printf("failed to build query: %v", err)
		return nil, 0
	}

	var user desc.User
	var userInfoID int64
	var createdAt, updatedAt time.Time
	err = s.pool.QueryRow(ctx, query, args...).Scan(&user.Id, &userInfoID, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, 0
	}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, userInfoID
}

func (s *server) getUserInfo(ctx context.Context, userInfoID int64) (*desc.UserInfo, error) {
	builderUserInfoSelect := sq.Select("name", "email", "password", "password_confirm", "role_id").
		From("user_info").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userInfoID}).
		Limit(1)

	query, args, err := builderUserInfoSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var userInfo desc.UserInfo
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userInfo.Name, &userInfo.Email, &userInfo.Password, &userInfo.PasswordConfirm, &userInfo.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_info: %v", err)
	}

	return &userInfo, nil
}

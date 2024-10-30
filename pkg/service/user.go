package service

import (
	"context"
	"fmt"

	"github.com/fastbyt3/diy-mssql-test/pkg/models"
	"github.com/fastbyt3/diy-mssql-test/pkg/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{
		userRepo: &userRepo,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id int32) (*models.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, username, email string) error {
	user := &models.User{
		Username: username,
		Email:    email,
	}

	if len(username) == 0 {
		return fmt.Errorf("Username needs to be present")
	}

	return s.userRepo.CreateUser(ctx, user)
}

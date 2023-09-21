package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/be/internal/domain"
	"webook/be/internal/repository"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

var _ UserService = (*userService)(nil)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) SignUp(ctx context.Context, user domain.User) error {
	// 密码加密存储
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return u.userRepo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先找用户
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码了
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

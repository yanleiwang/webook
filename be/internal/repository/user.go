package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/be/internal/domain"
	"webook/be/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

type userRepository struct {
	dao dao.UserDAO
}

func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	return u.dao.Insert(ctx, u.domainToEntity(user))

}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u.entityToDomain(user), err
}

func (u *userRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Ctime: user.Ctime.UnixMilli(),
	}
}

func (u *userRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Phone:    user.Phone.String,
		Ctime:    time.UnixMilli(user.Ctime),
	}
}

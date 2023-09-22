package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/be/internal/domain"
	"webook/be/internal/repository/cache"
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
	FindById(ctx context.Context, id int64) (domain.User, error)
}

func NewUserRepository(dao dao.UserDAO, userCache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: userCache,
	}
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(user))

}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(user), err
}

func (r *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然是有数据
		return u, nil
	}
	// 没这个数据
	//if err == cache.ErrKeyNotExist {
	// 去数据库里面加载
	//}

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 我这里怎么办？
			// 打日志，做监控
			//return domain.User{}, err
		}
	}()
	return u, err

	// 这里怎么办？ err = io.EOF
	// 要不要去数据库加载？
	// 看起来我不应该加载？
	// 看起来我好像也要加载？

	// 选加载 —— 做好兜底，万一 Redis 真的崩了，你要保护住你的数据库
	// 我数据库限流呀！

	// 选不加载 —— 用户体验差一点

	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了，你也不知道有没有数据
}

func (r *userRepository) domainToEntity(user domain.User) dao.User {
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

func (r *userRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Phone:    user.Phone.String,
		Ctime:    time.UnixMilli(user.Ctime),
	}
}

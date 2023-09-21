package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱或者手机号已存在")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	// SELECT * FROM `users` WHERE `email`=?
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(&user).Error

	// 如何判断mysql 返回错误类型？
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		// mysql唯一键冲突错误码： 1062
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

// User 直接对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(persistent object)
type User struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 全部用户唯一
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值
	// 但是不能有多个 ""
	Phone sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string

	// 往这面加

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}

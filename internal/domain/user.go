package domain

import (
	"context"
	"time"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        int            `gorm:"primarykey;autoIncrement:true"`
	Email     string         `gorm:"type:varchar(100);column:email;unique"`
	Password  string         `gorm:"type:varchar(200);column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

// func (u *User) TableName() string {
// 	return "users"
// }

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10); err != nil {
		return err
	} else {
		u.Password = string(bytes)
	}
	return nil
}

type UserUseCase interface {
	Login(beegoCtx *beegoContext.Context, email, password string) (interface{}, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type UserLogin struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserLoginResponse struct {
	Token     string    `json:"token"`
	ExpiredAt string    `json:"expired_at"`
	User      UserLogin `json:"user"`
}

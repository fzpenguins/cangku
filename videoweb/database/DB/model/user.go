package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uid       int64  `json:"uid" gorm:"primaryKey;autoIncrement:true"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"delete_at"`
}

func (user *User) SetPassword(password string) error {
	//先加密再保存
	tempPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(tempPassword)
	return nil
}

func (user *User) VerifyPassword(password string) bool {
	//先解密再比对
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) UploadAvatar(url string) {
	user.AvatarUrl = url
}

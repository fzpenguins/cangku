package dao

import (
	"context"
	"strconv"
	"time"
	"videoweb/database/DB/model"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) CreateUser(user *model.User) (err error) {
	err = Db.Model(&model.User{}).Create(&user).Error
	return
}

func (dao *UserDao) UploadAvatar(url string, uid int64) (err error) {
	user, _ := dao.FindUserByUid(strconv.FormatInt(uid, 10))
	user.AvatarUrl = url
	err = dao.DB.Save(user).Error
	return
}

func (dao *UserDao) FindUserByUid(uid string) (*model.User, error) {
	var retUser *model.User
	u, _ := strconv.ParseInt(uid, 10, 64)
	err := Db.Model(&model.User{}).Where("uid = ?", u).First(&retUser).Error
	return retUser, err
}

func (dao *UserDao) FindUserByName(name string) (*model.User, error) {
	var retUser *model.User
	err := Db.Model(&model.User{}).Where("username = ?", name).First(&retUser).Error
	return retUser, err
}

func (dao *UserDao) FindUserByNameAndUid(name string, uid int64) (*model.User, error) {
	var retUser *model.User
	err := Db.Model(&model.User{}).Where("uid = ?", uid).Where("username = ?", name).First(&retUser).Error
	return retUser, err
}

func (dao *UserDao) UpdateDate(user *model.User) error {
	user.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	err := Db.Save(user).Error
	return err
}

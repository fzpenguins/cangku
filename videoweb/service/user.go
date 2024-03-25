package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"time"
	"videoweb/biz/model/hertz/user"
	"videoweb/biz/utils"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/response"
)

func Register(c context.Context, userRegister *user.UserRegisterReq) (interface{}, error) {
	userDao := dao.NewUserDao(c)
	if len(userRegister.GetUsername()) == 0 || len(userRegister.GetUsername()) > 10 {
		return response.BadResponse(), errors.New("用户名长度应小于10")
	}
	if len(userRegister.GetPassword()) <= 5 {
		return response.BadResponse(), errors.New("密码长度不小于5")
	}
	_, err := userDao.FindUserByName(userRegister.Username)
	if err == nil {
		return response.BadResponse(), err
	}
	u := model.User{
		Uid:       0,
		Username:  userRegister.Username,
		Password:  "",
		AvatarUrl: "",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: "",
		DeletedAt: "",
	}
	err = u.SetPassword(userRegister.Password)
	if err != nil {
		return response.BadResponse(), err
	}

	err = dao.Db.Model(&model.User{}).Create(&u).Error
	if err != nil {
		return response.BadResponse(), err
	}
	return response.GoodResponse(), nil
}

func Login(c context.Context, userLogin *user.UserLoginReq) (interface{}, error) {
	userDao := dao.NewUserDao(c)
	usr, err := userDao.FindUserByName(userLogin.Username)
	if err != nil {
		return response.BadResponse(), err
	}
	if !usr.VerifyPassword(userLogin.Password) {
		return response.BadResponse(), nil
	}
	err = userDao.UpdateDate(usr)
	if err != nil {
		return response.BadResponse(), err
	}
	accessToken, refreshToken, err := utils.GenerateToken(usr.Uid, usr.Username)
	if err != nil {
		return response.BadResponse(), err
	}

	return response.TokenData{
		Data:         response.BestUserResponse(usr),
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func UploadAvatar(c context.Context, ctx *app.RequestContext, url string) (interface{}, error) {
	userDao := dao.NewUserDao(c)
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	usr, _ := userDao.FindUserByUid(strconv.FormatInt(claim.Uid, 10))
	err = userDao.UpdateDate(usr)
	if err != nil {
		return response.BadResponse(), err
	}
	err = userDao.UploadAvatar(url, usr.Uid)
	if err != nil {
		return response.BadResponse(), err
	}
	usr.AvatarUrl = url
	return response.BestUserResponse(usr), nil
}

func SearchUserInfo(c context.Context, ctx *app.RequestContext, req *user.UserInfoReq) (interface{}, error) {
	userDao := dao.NewUserDao(c)
	usr, err := userDao.FindUserByUid(req.GetUid())
	if err != nil {
		return response.BadResponse(), err
	}
	return response.BestUserResponse(usr), nil
}

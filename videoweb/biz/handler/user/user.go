// Code generated by hertz generator.

package user

import (
	"context"
	"videoweb/response"
	"videoweb/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	user "videoweb/biz/model/hertz/user"
)

// Register .
// @router /user/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserRegisterReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.Register(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}
	//resp := new(user.GoodResponse)

	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /user/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Info .
// @router /user/info [GET]
func Info(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.SearchUserInfo(ctx, c, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}
	c.JSON(consts.StatusOK, resp)
}

// AvatarUpload .
// @router /user/avatar/upload [PUT]
func AvatarUpload(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserUploadReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.UploadAvatar(ctx, c, req.GetAvatarUrl())
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}
	c.JSON(consts.StatusOK, resp)
}

// MFAGet .
// @router /auth/mfa/qrcode [GET]
func MFAGet(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserMFAReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.GetMfa(ctx, c)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// MFABind .
// @router /auth/mfa/bind [POST]
func MFABind(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserBindMFAReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.BindMFA(ctx, c, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

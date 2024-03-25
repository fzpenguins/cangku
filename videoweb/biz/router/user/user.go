// Code generated by hertz generator. DO NOT EDIT.

package user

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	user "videoweb/biz/handler/user"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_auth := root.Group("/auth", _authMw()...)
		{
			_mfa := _auth.Group("/mfa", _mfaMw()...)
			_mfa.POST("/bind", append(_mfabindMw(), user.MFABind)...)
			_mfa.GET("/qrcode", append(_mfagetMw(), user.MFAGet)...)
		}
	}
	{
		_user := root.Group("/user", _userMw()...)
		_user.GET("/info", append(_infoMw(), user.Info)...)
		_user.POST("/login", append(_loginMw(), user.Login)...)
		_user.POST("/register", append(_registerMw(), user.Register)...)
		{
			_avatar := _user.Group("/avatar", _avatarMw()...)
			_avatar.PUT("/upload", append(_avataruploadMw(), user.AvatarUpload)...)
		}
	}
}
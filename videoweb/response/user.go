package response

import (
	"strconv"
	"videoweb/biz/model/hertz/common"
	"videoweb/biz/model/hertz/following"
	"videoweb/biz/model/hertz/user"
	"videoweb/database/DB/model"
	"videoweb/pkg/e"
)

func BestUserResponse(usr *model.User) interface{} {
	return user.GoodUserResponse{
		Base: &common.BaseResponse{
			Code: e.SuccessCode,
			Msg:  e.SuccessMsg,
		},
		Data: &user.DataResponse{
			Uid:       strconv.FormatInt(usr.Uid, 10),
			Username:  usr.Username,
			AvatarUrl: usr.AvatarUrl,
			CreatedAt: usr.CreatedAt,
			UpdatedAt: usr.UpdatedAt,
			DeletedAt: usr.DeletedAt,
		},
	}
}

func UserInfoResponse(users []*model.User, size int64) interface{} {
	var res []*following.ToUserInfoResponse
	for _, v := range users {
		item := &following.ToUserInfoResponse{
			Uid:       strconv.FormatInt(v.Uid, 10),
			Username:  v.Username,
			AvatarUrl: v.AvatarUrl,
		}
		res = append(res, item)
	}
	return following.UserInfoResponse{
		Base: &common.BaseResponse{
			Code: e.SuccessCode,
			Msg:  e.SuccessMsg,
		},
		Data: &following.UserDataResponse{
			Items: res,
			Total: size,
		},
	}
}

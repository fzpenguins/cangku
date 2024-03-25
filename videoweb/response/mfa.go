package response

import (
	"videoweb/biz/model/hertz/common"
	"videoweb/biz/model/hertz/user"
	"videoweb/pkg/e"
)

func MFAResponse(secret string, qrcode string) user.GoodMFAResponse {
	return user.GoodMFAResponse{
		Base: &common.BaseResponse{
			Code: e.SuccessCode,
			Msg:  e.SuccessMsg,
		},
		Data: &user.MFAData{
			Secret: secret,
			Qrcode: qrcode,
		},
	}
}

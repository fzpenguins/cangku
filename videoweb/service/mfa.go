package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/pquerna/otp/totp"
	"log"
	"videoweb/biz/model/hertz/user"
	"videoweb/biz/utils"
	"videoweb/database/DB/dao"
	"videoweb/response"
)

func GetMfa(c context.Context, ctx *app.RequestContext) (interface{}, error) {
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}

	img, err := dao.GenerateQRCode(claim.Uid)
	if err != nil {
		return response.BadResponse(), err
	}

	imgByBase64, err := dao.ImageToBase64(img)
	if err != nil {
		return response.BadResponse(), err
	}
	return response.MFAResponse(dao.Key.Secret(), imgByBase64), nil
}

func BindMFA(c context.Context, ctx *app.RequestContext, MFAReq *user.UserBindMFAReq) (interface{}, error) {
	valid := totp.Validate(MFAReq.GetSecret(), MFAReq.GetCode()) //key设置为全局变量

	if !valid {
		log.Println("Invalid passcode!")
		return response.BadResponse(), errors.New("Invalid MFA secret!")
	}
	//codeBytes, err := base64.StdEncoding.DecodeString(MFAReq.GetCode())
	//if err != nil {
	//	return response.BadResponse(), err
	//}

	return response.GoodResponse(), nil
}

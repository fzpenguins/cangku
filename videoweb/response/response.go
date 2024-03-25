package response

import (
	"videoweb/biz/model/hertz/common"
	"videoweb/biz/model/hertz/video"
	"videoweb/pkg/e"
)

func BadResponse() interface{} {
	return common.BaseResponse{
		Code: e.FailureCode,
		Msg:  e.FailureMsg,
	}
}

func GoodResponse() interface{} {
	return common.BaseResponse{
		Code: e.SuccessCode,
		Msg:  e.SuccessMsg,
	}
}

func TotalResp(items []*video.ItemsResponse, size int) interface{} {
	var res video.VideoTotalResponse

	res.Items = items

	res.Base = &common.BaseResponse{
		Code: e.SuccessCode,
		Msg:  e.SuccessMsg,
	}
	res.Total = int64(size)
	return res
}

func Resp(items []*video.ItemsResponse) interface{} {
	var res video.VideoResponse
	res.Items = items
	res.Base = &common.BaseResponse{
		Code: e.SuccessCode,
		Msg:  e.SuccessMsg,
	}
	return res
}

type TokenData struct {
	Data         interface{} `json:"data"`
	RefreshToken string      `json:"refresh_token"`
	AccessToken  string      `json:"access_token"`
}

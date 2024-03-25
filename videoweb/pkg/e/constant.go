package e

import "time"

//response的相关设置
const (
	SuccessCode = 10000
	SuccessMsg  = "success"

	FailureCode = -1
	FailureMsg  = "密码错误"
)

//设有accesstoken和refreshtoken的有效时间
const (
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

//消息类型
const (
	TypePrivateMessage   = "type1"
	TypeGetHistory       = "type2"
	TypeGetUnreadHistory = "type3"
	TypeGroupMessage     = "type4"
	TypeGetGroupHistory  = "type5"
)

//与websocket相关的
const (
	PageSize = 10

	SingleChat               = 0
	GroupChat                = 1
	GetHistoryFromSingleChat = 2
	GetUnreadFromSingleChat  = 3
	GetHistoryFromGroupChat  = 4

	MaxStore = 4 * 1024
)

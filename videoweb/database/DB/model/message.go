package model

//type Message struct {
//	SendUser     string // 发送人
//	ReceivedUser string // 接收人
//	SendTime     string // 发送时间
//	Content      string // 消息内容
//	IsPublic     bool   // 消息类型是否是公开的 true 公开 false 私信
//	IsReceived   bool   // 接收人是否接收成功 true 接收成功 false 离线还未接收（当接收人离线时，设置为false，当对方上线时，将消息发过去，改为true）
//	IsSend       bool   // 是否是发送消息，用于区分发送消息和上线下线消息（true 发送消息 false 上线/下线消息）
//	IsImg        bool   // 头像是否是图片
//	Pic          string // 头像图片地址
//}

type Message struct {
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
	FromUid   string `json:"from_uid"`
	ToUid     string `json:"to_uid"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	ReadTag   bool   `json:"read_tag"`
}

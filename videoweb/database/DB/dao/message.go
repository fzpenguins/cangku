package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"videoweb/database/DB/model"
	"videoweb/pkg/e"
)

type MsgDao struct {
	*gorm.DB
}

func GetMsgDao(ctx context.Context) *MsgDao {
	return &MsgDao{NewDBClient(ctx)}
}

//分页的详情由server端设定即可，无需client传入

func (dao *MsgDao) GetHistoryFromSingleChat(pageNum int, from, to string) (msgs []*model.Message, err error) {
	err = dao.DB.Model(&model.Message{}).Where("from_uid = ? AND to_uid = ? AND type = ?", from, to, 0).Limit(e.PageSize).
		Offset(pageNum * e.PageSize).Find(&msgs).Error
	return
}

func (dao *MsgDao) GetHistoryFromGroupChat(pageNum int, to string) (msgs []*model.Message, err error) {
	err = dao.DB.Model(&model.Message{}).Where("to_uid = ? AND type = ?", to, 1).Limit(e.PageSize).
		Offset(pageNum * e.PageSize).Find(&msgs).Error
	return
}

func (dao *MsgDao) GetUnreadFromSingleChat(pageNum int, from, to string) (msgs []*model.Message, err error) {
	err = dao.DB.Model(&model.Message{}).Where("from_uid = ? AND to_uid = ? AND read_tag = ?", from, to, false).Limit(e.PageSize).
		Offset(pageNum * e.PageSize).Find(&msgs).Error
	return
}

func (dao *MsgDao) StoreSingleChatMsg(from, to, content string, readTag bool) error {
	msg := &model.Message{
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		DeletedAt: "",
		FromUid:   from,
		ToUid:     to,
		Type:      e.SingleChat,
		Content:   content,
		ReadTag:   readTag,
	}
	return dao.DB.Model(&model.Message{}).Create(&msg).Error
}

func (dao *MsgDao) TurnToRead(from, to string) error {
	return dao.Model(&model.Message{}).Where("from_uid = ? AND to_uid = ?", from, to).Update("read_tag", true).Error
}

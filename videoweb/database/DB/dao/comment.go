package dao

import (
	"context"
	"gorm.io/gorm"
	"videoweb/database/DB/model"
)

type CommentDao struct {
	*gorm.DB
}

func NewComment(ctx context.Context) *CommentDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &CommentDao{NewDBClient(ctx)}
}

func (dao *CommentDao) CreateComment(comment *model.Comment) (err error) {
	err = dao.DB.Model(&model.Comment{}).Create(&comment).Error
	return
}

func (dao *CommentDao) FindCommentByCid(cid int64) (comment model.Comment, err error) {
	err = dao.DB.Model(&model.Comment{}).Where("cid = ?", cid).First(&comment).Error
	return
}

func (dao *CommentDao) FindCommentByVid(vid int64, pageSize int64, pageNum int64) (comment []*model.Comment, err error) {
	err = dao.DB.Model(&model.Comment{}).Where("vid = ?", vid).Limit(int(pageNum)).Offset(int(pageSize * pageNum)).
		Find(&comment).Error
	return comment, err
}

func (dao *CommentDao) FindCommentByCidAndVid(vid, cid int64, pageSize, pageNum int64) (comment []*model.Comment, err error) {
	err = dao.DB.Model(&model.Comment{}).Where("vid = ? AND cid = ?", vid, cid).Limit(int(pageNum)).Offset(int(pageSize * pageNum)).
		Find(&comment).Error
	return comment, err
}

func (dao *CommentDao) FindCommentInCid(cid int64, pageSize int64, pageNum int64) (comment []*model.Comment, err error) {
	err = dao.DB.Model(&model.Comment{}).Where(" cid = ?", cid).Limit(int(pageNum)).Offset(int(pageSize * pageNum)).
		Find(&comment).Error
	return comment, err
}

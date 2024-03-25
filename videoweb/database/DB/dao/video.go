package dao

import (
	"context"
	"gorm.io/gorm"
	"videoweb/database/DB/model"
)

type VideoDao struct {
	*gorm.DB
}

func NewVideo(ctx context.Context) *VideoDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &VideoDao{NewDBClient(ctx)}
}

func (dao *VideoDao) CreateVideo(video *model.Video) (err error) {
	err = dao.DB.Model(&model.Video{}).Create(&video).Error
	return
}

func (dao *VideoDao) FindVideoByVid(vid int64) (video model.Video, err error) {

	err = dao.DB.Model(&model.Video{}).Where("vid = ?", vid).First(&video).Error
	return
}

func (dao *VideoDao) FindVideosByTimeStr(timeStr string) (videoList []*model.Video, err error) {
	err = dao.DB.Model(&model.Video{}).Where("created_at >= ?", timeStr).Limit(50).Order("created_at desc").Find(&videoList).Error
	return
}

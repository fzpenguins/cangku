package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"time"
	"videoweb/biz/model/hertz/video"
	"videoweb/biz/utils"
	"videoweb/config"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/database/cache"
	"videoweb/response"
)

func PublishVideo(c context.Context, ctx *app.RequestContext, videoInfo *video.VideoPublishReq) (interface{}, error) {
	if len(videoInfo.GetTitle()) == 0 {
		return response.BadResponse(), errors.New("标题不能为空")
	}
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	videoName := GenerateVideoName(claim.Uid)
	coverName := GenerateCoverName(claim.Uid)
	videoURL := fmt.Sprintf("https://%s/%s/%s", config.EndPoint, config.BucketName, videoName)
	coverURL := fmt.Sprintf("https://%s/%s/%s", config.EndPoint, config.BucketName, coverName)
	//要把视频文件存储到minio的步骤还没写
	videoKey, err := UploadVideo(videoInfo.Data.VideoUrl)
	if err != nil {
		fmt.Println(err)
		return response.BadResponse(), err
	}
	coverKey, err := UploadCover(videoInfo.Data.CoverUrl)
	if err != nil {
		return response.BadResponse(), err
	}

	v := model.Video{
		Uid:         claim.Uid,
		VideoUrl:    videoURL,
		CoverUrl:    coverURL,
		Title:       videoInfo.GetTitle(),
		Description: videoInfo.GetDescription(),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),

		Key:      videoKey,
		CoverKey: coverKey,
	}
	err = dao.Db.Model(&model.Video{}).Create(&v).Error
	if err != nil {
		return response.BadResponse(), err
	}
	cache.AddVisitCount(c, &v)
	return response.GoodResponse(), nil
}

func ListVideos(c context.Context, ctx *app.RequestContext, videoRequest *video.VideoListReq) (interface{}, error) {
	uid, _ := strconv.ParseInt(videoRequest.GetUid(), 10, 64)
	if videoRequest.GetPageNum() < 0 || videoRequest.GetPageSize() < 0 {
		return response.BadResponse(), errors.New("请重新操作")
	}
	_, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	var videos []model.Video
	var count int64
	err = dao.Db.Model(&model.Video{}).Where("uid = ?", uid).
		Order("updated_at desc").Limit(int(videoRequest.PageSize)).
		Offset(int(videoRequest.PageNum * videoRequest.PageSize)).Find(&videos).Count(&count).Error //分页操作
	if err != nil {
		return response.BadResponse(), err
	}

	var videosInfo []*video.ItemsResponse
	for _, value := range videos {
		value.VideoUrl, _ = GetURL(value.Key)
		value.CoverUrl, _ = GetURL(value.CoverKey)
		value.LikeCount = int(cache.LikeCount(c, value.Vid))
		value.VisitCount = int(cache.VisitCount(c, value.Vid))
		videosInfo = append(videosInfo, response.MakeVideoResponse(&value))
	}
	return response.TotalResp(videosInfo, int(count)), nil
}

func PopularRank(c context.Context, ctx *app.RequestContext, videoRequest *video.VideoPopularReq) (interface{}, error) {

	var videos []model.Video
	if videoRequest.GetPageNum() < 0 || videoRequest.GetPageSize() < 0 {
		return response.BadResponse(), errors.New("请重新操作")
	}
	//err := dao.Db.Model(&model.Video{}).Order("visit_count desc").Limit(int(videoRequest.GetPageSize())).
	//	Offset(int(videoRequest.GetPageNum() * videoRequest.GetPageSize())).Find(&videos).Error //分页操作
	//if err != nil {
	//	return response.BadResponse(), err
	//}
	//count := len(videos)
	//offset := videoRequest.GetPageNum() * videoRequest.GetPageSize()
	//cnt, resp, vidString := cache.FindVideoVidsRank(c)
	vids, offset, cnt := cache.FindVideoVidsRank(c, videoRequest)
	//for i:= offset;i-offset<cnt;i++{
	//	vid,_ := strconv.ParseInt(vidString[i],10,64)
	//	vids = append(vids, vid)
	//}
	err := dao.Db.Model(&model.Video{}).Where("vid IN ?", vids).Find(&videos).Error
	if err != nil {
		return response.BadResponse(), err
	}
	for i := offset; i-offset < videoRequest.GetPageSize() && i < cnt; i++ {
		// cache.AddVisitCount(c, &videos[i])                             //增加点击量
		videos[i].VisitCount = int(cache.VisitCount(c, vids[i])) //获取点击量
		videos[i].LikeCount = int(cache.LikeCount(c, vids[i]))
	}

	var videosInfo []*video.ItemsResponse

	for _, value := range videos {
		value.VideoUrl, _ = GetURL(value.Key)
		value.CoverUrl, _ = GetURL(value.CoverKey)
		videosInfo = append(videosInfo, response.MakeVideoResponse(&value))
	}
	return response.Resp(videosInfo), nil
}

func SearchVideo(c context.Context, ctx *app.RequestContext, videoRequest *video.VideoSearchReq) (interface{}, error) {
	if videoRequest.GetPageNum() < 0 || videoRequest.GetPageSize() < 0 {
		return response.BadResponse(), errors.New("请重新操作")
	}
	var videos []model.Video
	//var query *gorm.DB
	query := dao.Db.Where("description LIKE ? OR title LIKE ?", "%"+videoRequest.Keywords+"%", "%"+videoRequest.Keywords+"%")

	if videoRequest.GetUsername() != "" {
		d := dao.NewUserDao(c)
		user, _ := d.FindUserByName(videoRequest.GetUsername())
		query = query.Where("uid = ?", user.Uid)

	}

	if videoRequest.GetFromDate() > 0 {
		FromDateString := time.Unix(videoRequest.GetFromDate(), 0).Format("2006-01-02 15:04:05")
		query = query.Where("created_at >= ?", FromDateString)
	}
	if videoRequest.GetToDate() > 0 {
		ToDateString := time.Unix(videoRequest.GetToDate(), 0).Format("2006-01-02 15:04:05")
		query = query.Where("created_at >= ?", ToDateString)
	}
	var size int64
	query.Find(&videos).Count(&size)

	offset := videoRequest.PageNum * videoRequest.PageSize
	for i := offset; i-offset < size; i++ {
		// cache.AddVisitCount(c, &videos[i])                             //增加点击量
		videos[i].VisitCount = int(cache.VisitCount(c, videos[i].Vid)) //获取点击量
		videos[i].LikeCount = int(cache.LikeCount(c, videos[i].Vid))
	}
	var videosInfo []*video.ItemsResponse

	for _, value := range videos {
		value.VideoUrl, _ = GetURL(value.Key)
		value.CoverUrl, _ = GetURL(value.CoverKey)
		videosInfo = append(videosInfo, response.MakeVideoResponse(&value))
	}
	return response.TotalResp(videosInfo, int(size)), nil
}

func FeedVideo(c context.Context, ctx *app.RequestContext, videoReq *video.VideoFeedReq) (interface{}, error) {
	var videosInfo []*video.ItemsResponse
	//claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	//if err != nil {
	//	return response.BadResponse(), err
	//}
	rawTimeStamp := videoReq.GetLatestTime()
	//if req.LatestTime == nil {
	//	currentTime := time.Now().UnixMilli()
	//	req.LatestTime = &currentTime
	//} 对于latesttime来说，没有则取现在，否则取选定的值
	intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
	if err != nil {
		return response.BadResponse(), err
	}
	timeDate := time.Unix(intTime, 0).Format("2006-01-02 15:04:05")
	videoDao := dao.NewVideo(c)
	videoList, err := videoDao.FindVideosByTimeStr(timeDate)
	if err != nil {
		return response.BadResponse(), err
	}

	for _, m := range videoList {
		m.VisitCount = int(cache.VisitCount(c, m.Vid))
		m.LikeCount = int(cache.LikeCount(c, m.Vid))
		videosInfo = append(videosInfo, response.MakeVideoResponse(m))
	}
	return response.Resp(videosInfo), nil
}

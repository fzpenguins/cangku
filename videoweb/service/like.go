package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/redis/go-redis/v9"
	"strconv"
	"videoweb/biz/model/hertz/like"
	"videoweb/biz/model/hertz/video"
	"videoweb/biz/utils"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/database/cache"
	"videoweb/response"
)

func ActionLike(c context.Context, ctx *app.RequestContext, likeRequest *like.ActionReq) (interface{}, error) {
	if len(likeRequest.GetVid()) == 0 && len(likeRequest.GetCid()) == 0 {
		return response.BadResponse(), errors.New("没有选中点赞对象")
	}
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	//对视频进行点赞
	if len(likeRequest.GetVid()) != 0 {
		videoDao := dao.NewVideo(c)
		//v := new(model.Video)
		tempNum, err := strconv.ParseInt(likeRequest.GetVid(), 10, 64)
		if err != nil {
			return response.BadResponse(), err
		}
		v, err := videoDao.FindVideoByVid(tempNum)
		if err != nil {
			return response.BadResponse(), nil
		}

		if likeRequest.GetActionType() == "1" {
			cache.AddLikeCount(c, claim.Uid, &v)
		} else if likeRequest.GetActionType() == "2" {
			cache.DecrLikeCount(c, claim.Uid, &v)
		} else {
			return response.BadResponse(), errors.New("操作失败")
		}

		return response.GoodResponse(), nil
	} else {
		commentDao := dao.NewComment(c)
		tempNum, err := strconv.ParseInt(likeRequest.GetCid(), 10, 64)
		if err != nil {
			return response.BadResponse(), err
		}
		com, err := commentDao.FindCommentByCid(tempNum)
		if err != nil {
			return response.BadResponse(), nil
		}

		if likeRequest.GetActionType() == "1" {
			cache.AddCommentLikeCount(c, claim.Uid, &com)
		} else if likeRequest.GetActionType() == "2" {
			cache.DecrCommentLikeCount(c, claim.Uid, &com)
		} else {
			return response.BadResponse(), errors.New("操作失败")
		}

		return response.GoodResponse(), nil
	}
}

func ListLikes(c context.Context, ctx *app.RequestContext, likeRequest *like.ListReq) (interface{}, error) {
	tempNum, err := strconv.ParseInt(likeRequest.GetUid(), 10, 64)
	key := cache.GetVideoLikeFromUser(tempNum)
	var length []string
	pos := 0
	for {

		length = cache.RedisClient.ZRevRangeByScore(c, key, &redis.ZRangeBy{
			Min:    "1",
			Max:    "1",
			Offset: int64(pos),
			Count:  likeRequest.GetPageSize(),
		}).Val()
		if cache.RedisClient.ZCard(c, key).Val() <= likeRequest.GetPageSize() {
			break
		}
		pos += len(length)
		if int64(pos) >= likeRequest.GetPageNum()*likeRequest.GetPageSize() {
			break
		}
	}

	var vids []int64
	for _, v := range length {
		c, _ := strconv.ParseInt(v, 10, 64)
		vids = append(vids, c)
	}
	var res []model.Video
	err = dao.Db.Model(&model.Video{}).Where("vid IN ?", vids).Find(&res).Error
	if err != nil {
		return response.BadResponse(), err
	}
	var videosInfo []*video.ItemsResponse

	for _, value := range res { //还得把redis的点击量和点赞量赋值
		value.VisitCount = int(cache.VisitCount(c, value.Vid))
		value.LikeCount = int(cache.LikeCount(c, value.Vid))
		videosInfo = append(videosInfo, response.MakeVideoResponse(&value))
	}
	return response.Resp(videosInfo), nil

}

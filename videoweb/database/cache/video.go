package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"videoweb/biz/model/hertz/video"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
)

func VisitCount(ctx context.Context, vid int64) int64 {
	countStr, _ := RedisClient.Get(ctx, GetVideoVisitCountKey(vid)).Result()
	count, _ := strconv.ParseInt(countStr, 10, 64)
	return count
}

func LikeCount(ctx context.Context, vid int64) int64 {
	countStr, _ := RedisClient.Get(ctx, GetVideoLikeCountKey(vid)).Result()
	count, _ := strconv.ParseInt(countStr, 10, 64)
	return count
}

func AddVisitCount(ctx context.Context, video *model.Video) {
	RedisClient.Incr(ctx, GetVideoVisitCountKey(video.Vid))
	RedisClient.ZIncrBy(ctx, "visit_count", 1, strconv.Itoa(int(video.Vid)))
}

func AddLikeCount(ctx context.Context, uid int64, video *model.Video) {
	if RedisClient.ZScore(ctx, GetVideoLikeFromUser(uid), strconv.FormatInt(video.Vid, 10)).Val() != 1 {
		RedisClient.Incr(ctx, GetVideoLikeCountKey(video.Vid))
		RedisClient.ZIncrBy(ctx, "like_count", 1, strconv.Itoa(int(video.Vid)))
	}
	RedisClient.ZAdd(ctx, GetVideoLikeFromUser(uid), redis.Z{Member: video.Vid, Score: 1})
	//RedisClient.SetBit(ctx, GetVideoLikeFromUser(uid), video.Vid, 1)
}

func DecrLikeCount(ctx context.Context, uid int64, video *model.Video) {
	if RedisClient.ZScore(ctx, GetVideoLikeFromUser(uid), strconv.FormatInt(video.Vid, 10)).Val() != 0 {
		RedisClient.Decr(ctx, GetVideoLikeCountKey(video.Vid))
		RedisClient.ZIncrBy(ctx, "like_count", -1, strconv.Itoa(int(video.Vid)))
	}
	RedisClient.ZAdd(ctx, GetVideoLikeFromUser(uid), redis.Z{Member: video.Vid, Score: 0})
	//RedisClient.Decr(ctx, GetVideoLikeCountKey(video.Vid))
	//RedisClient.ZIncrBy(ctx, "like_count", -1, strconv.Itoa(int(video.Vid)))
}

func GetVideoLikeCountKey(vid int64) string {
	return fmt.Sprintf("like_count/video:%s", strconv.Itoa(int(vid)))
}

func GetVideoLikeFromUser(uid int64) string {
	return fmt.Sprintf("%s/like_count/video", strconv.Itoa(int(uid))) //传入点赞的人的uid
}

func GetVideoVisitCountKey(vid int64) string {
	return fmt.Sprintf("visit_count/video:%s", strconv.Itoa(int(vid)))
}

func VideoCacheKey(vid int64) string {
	return fmt.Sprintf("video:%s", strconv.Itoa(int(vid)))
}

func QueryVideoKey(vid int64) string {
	return fmt.Sprintf("video_history/user:%s", strconv.Itoa(int(vid)))
}

func UploadVideoInCache(ctx context.Context, video model.Video) error {
	key := VideoCacheKey(video.Vid)
	jsonStr, err := json.Marshal(video)
	if err != nil {
		return err
	}
	RedisClient.Set(ctx, key, string(jsonStr), 300)
	return nil
}

func FindVideoByVid(ctx context.Context, vid int64) (model.Video, error) {
	key := VideoCacheKey(vid)
	retStringCmd := RedisClient.Get(ctx, key)
	var v model.Video
	if retStringCmd.Err() == redis.Nil {
		//不在缓存中
		videoDao := dao.NewVideo(ctx)
		v, err := videoDao.FindVideoByVid(vid)
		if err != nil {
			return v, nil
		}

		jsonStr, _ := json.Marshal(v)
		RedisClient.Set(ctx, key, string(jsonStr), 300)
		return v, nil
	}
	jsonStr := retStringCmd.Val()
	json.Unmarshal([]byte(jsonStr), v)
	return v, nil
}

func SaveQueryVideoResult(ctx context.Context) {

}

func FindVideoVidsRank(ctx context.Context, videoRequest *video.VideoPopularReq) (vids []int64, offset int64, cnt int64) { // ([]int64, []string){//(int64, []video.VideoRankResp, []string) {
	//resp := make([]video.VideoRankResp, 0)
	offset = videoRequest.GetPageNum() * videoRequest.GetPageSize()
	cnt, _ = RedisClient.ZCard(ctx, "visit_count").Result()

	Stringcmds, _ := RedisClient.ZRevRange(ctx, "visit_count", 0, -1).Result()
	if offset >= int64(len(Stringcmds)) {
		return []int64{}, 0, 0
	}
	for i := offset; i-offset < videoRequest.GetPageSize() && i < cnt; i++ {
		vid, _ := strconv.ParseInt(Stringcmds[i], 10, 64)
		vids = append(vids, vid)
	}
	//for _, Stringcmd := range Stringcmds {
	//	v := RedisClient.ZScore(ctx, "visit_count", Stringcmd).Val()
	//	vid, _ := strconv.Atoi(Stringcmd)
	//	resp = append(resp, video.VideoRankResp{
	//		Vid:        strconv.FormatInt(int64(vid), 10),
	//		VisitCount: int64(v),
	//	})
	//}

	return
}

func AddVideoComments(ctx context.Context, vid, cid string) error {
	err := RedisClient.SAdd(ctx, VideoCommentKey(vid), cid).Err()
	return err
}

func DecrVideoComments(ctx context.Context, vid, cid string) error {
	err := RedisClient.SRem(ctx, VideoCommentKey(vid), cid).Err()
	return err
}

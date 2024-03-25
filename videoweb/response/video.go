package response

import (
	"strconv"
	"videoweb/biz/model/hertz/video"
	"videoweb/database/DB/model"
)

func MakeVideoResponse(v *model.Video) *video.ItemsResponse {
	return &video.ItemsResponse{
		Vid:          strconv.FormatInt(v.Vid, 10),
		Uid:          strconv.FormatInt(v.Uid, 10),
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   int64(v.VisitCount),
		LikeCount:    int64(v.LikeCount),
		CommentCount: int64(v.CommentCount),
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
		DeletedAt:    v.DeletedAt,
	}
}

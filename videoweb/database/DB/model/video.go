package model

type Video struct {
	Vid          int64  `json:"vid" gorm:"primaryKey;autoIncrement:true"`
	Uid          int64  `json:"uid"`
	VideoUrl     string `json:"video_url"`
	CoverUrl     string `json:"cover_url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VisitCount   int    `json:"visit_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
	Key          string `json:"key"`
	CoverKey     string `json:"cover_key"`
}

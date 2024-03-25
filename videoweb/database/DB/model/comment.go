package model

type Comment struct {
	Uid        int64  `json:"uid" `
	Vid        int64  `json:"vid"`
	Cid        int64  `json:"cid" gorm:"primaryKey;autoIncrement:true"`
	ParentId   int64  `json:"parent_id"`
	LikeCount  int64  `json:"like_count"`
	ChildCount int64  `json:"child_count"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"delete_at"`
}

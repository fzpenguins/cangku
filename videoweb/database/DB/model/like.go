package model

type Like struct {
	FromUid int64 `json:"from_uid"`
	ToUid   int64 `json:"to_uid"`
}

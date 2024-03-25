package model

type Relation struct {
	FromUid int64 `json:"from_uid"`
	ToUid   int64 `json:"to_uid"`
	Status  int64 `json:"status"`
}

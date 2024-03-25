package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
	"videoweb/biz/model/hertz/friends"
	"videoweb/biz/utils"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/response"
)

func ListFriends(c context.Context, ctx *app.RequestContext, req *friends.FriendsListReq) (interface{}, error) {
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	var myFriends []*model.User
	//互关为好友

	////获得的uid能够在follower中查找到，就加入，否则跳过
	var res []model.Relation
	err = dao.Db.Model(&model.Relation{}).Where("from_uid = ? AND status = ?", claim.Uid, 0).Find(&res).Error
	if err != nil {
		return response.BadResponse(), err
	}
	//var resp []model.Relation
	var wg sync.WaitGroup
	result := make(chan *model.User)
	for i := 0; i < len(res); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var d model.Relation
			var u *model.User
			dao.Db.Model(&model.Relation{}).Where("from_uid = ? AND to_uid = ? AND status = ?", res[i].ToUid, claim.Uid, 0).Find(&d)
			dao.Db.Model(&model.User{}).Where("uid = ?", d.FromUid).Find(&u)

			result <- u
		}(i)
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	for user := range result {
		if user.Uid != 0 {
			myFriends = append(myFriends, user)
		}
	}
	totalPage := (int64(len(myFriends)) + req.GetPageSize() - 1) / req.GetPageSize()
	if req.GetPageNum() < 0 {
		return response.BadResponse(), errors.New("请重新操作")
	} else if req.GetPageNum() >= totalPage && totalPage > 0 { //从0开始计算的
		return response.BadResponse(), errors.New("请重新操作")
	}

	return response.UserInfoResponse(myFriends, int64(len(myFriends))), nil
}

package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"sync"
	"videoweb/biz/model/hertz/following"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/response"
)

func ListFollower(c context.Context, ctx *app.RequestContext, req *following.FollowingListReq) (interface{}, error) {
	toUid, _ := strconv.ParseInt(req.GetUid(), 10, 64)
	var res []model.Relation
	var count int64
	if req.GetPageNum() >= 0 && req.GetPageSize() >= 0 {
		err := dao.Db.Model(&model.Relation{}).Offset(int(req.GetPageNum()*req.GetPageSize())).
			Limit(int(req.GetPageSize())).Where("to_uid = ? AND status = ?", toUid, 0).
			Find(&res).Count(&count).Error
		if err != nil {
			return response.BadResponse(), err
		}
	} else {
		return response.BadResponse(), errors.New("请重新操作")
	}

	var users []*model.User
	// 声明一个 WaitGroup
	var wg sync.WaitGroup

	// 声明一个通道用于接收结果
	resultCh := make(chan *model.User)
	for i := 0; i < len(res); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var item *model.User
			err := dao.Db.Model(&model.User{}).Where("uid = ? ", res[i].FromUid).Find(&item).Error
			if err == nil {
				resultCh <- item
			}
		}(i)
	}

	go func() {
		wg.Wait()       // 等待所有 goroutine 执行完成
		close(resultCh) // 关闭通道
	}()

	for user := range resultCh {
		users = append(users, user)
	}

	return response.UserInfoResponse(users, count), nil
}

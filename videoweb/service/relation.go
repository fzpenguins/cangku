package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"videoweb/biz/model/hertz/relation"
	"videoweb/biz/utils"
	"videoweb/database/DB/dao"
	"videoweb/database/DB/model"
	"videoweb/response"
)

func ActionRelation(c context.Context, ctx *app.RequestContext, req *relation.RelationActionReq) (interface{}, error) {
	toUid, _ := strconv.ParseInt(req.ToUid, 10, 64)
	var temp model.User
	//先检查uid是否正确
	err := dao.Db.Model(&model.User{}).Where("uid = ?", toUid).First(&temp).Error
	if err != nil {
		return response.BadResponse(), err
	}
	if req.GetActionType() != 0 && req.GetActionType() != 1 {
		return response.BadResponse(), err
	}
	claim, err := utils.ParseToken(string(ctx.GetHeader("access_token")))
	if err != nil {
		return response.BadResponse(), err
	}
	if toUid == claim.Uid {
		return response.BadResponse(), errors.New("请重新操作")
	}
	var r model.Relation
	var count int64
	dao.Db.Model(&model.Relation{}).Where("to_uid = ? AND from_uid = ?", toUid, claim.Uid).
		First(&r).Count(&count)
	if count == 0 {
		r = model.Relation{
			FromUid: claim.Uid,
			ToUid:   toUid,
			Status:  req.GetActionType(),
		}
		err = dao.Db.Model(&model.Relation{}).Create(&r).Error
		if err != nil {
			return response.BadResponse(), err
		}
	} else {
		err = dao.Db.Model(&model.Relation{}).Where("to_uid = ? AND from_uid = ?", toUid, claim.Uid).Update("status", req.GetActionType()).Error
		if err != nil {
			return response.BadResponse(), err
		}
	}
	return response.GoodResponse(), nil
}

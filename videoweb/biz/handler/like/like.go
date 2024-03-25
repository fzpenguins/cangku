// Code generated by hertz generator.

package like

import (
	"context"
	"videoweb/response"
	"videoweb/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	like "videoweb/biz/model/hertz/like"
)

// Action .
// @router /like/action [POST]
func Action(ctx context.Context, c *app.RequestContext) {
	var err error
	var req like.ActionReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.ActionLike(ctx, c, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// List .
// @router /like/list [GET]
func List(ctx context.Context, c *app.RequestContext) {
	var err error
	var req like.ListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	resp, err := service.ListLikes(ctx, c, &req)
	if err != nil {
		c.JSON(consts.StatusOK, response.BadResponse())
		return
	}

	c.JSON(consts.StatusOK, resp)
}

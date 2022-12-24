package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 社区相关逻辑
// @Summary 获取社区列表数据
// @Description 返回所有的社区数据
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunityList
// @Failure 400 {object} _ResponseCommunityList
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
// @Summary 获取社区分类详情
// @Description 返回社区分类详情
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path integer false "社区id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunityList
// @Failure 400 {object} _ResponseCommunityList
// @Router /community [get]
func CommunityDetailHandler(c *gin.Context) {
	// 获取url路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

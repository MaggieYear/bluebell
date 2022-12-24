package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子
// @Description 创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.Post false "帖子"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c获取当前登录用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = int64(userID)
	// 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}

	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 查看帖子详情
// @Summary 查看帖子详情
// @Description 查看帖子详情
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path integer false "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Router /post [get]
func GetPostDetailHandler(c *gin.Context) {
	// 获取url路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表
// @Summary 获取帖子列表
// @Description 获取帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param page query string false "分页页码"
// @Param size query string false "每页显示数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")

	var (
		pageNum  int64
		pageSize int64
		err      error
	)

	pageNum, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 0
	}
	pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("logic.GetPostListHandler() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostVoteHandler 投票
// @Summary 帖子投票接口
// @Description 帖子投票接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamVoteData false "投票参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Router /vote [post]
func PostVoteHandler(c *gin.Context) {
	p := new(models.ParamVoteData)
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("PostVoteHandler failed", zap.Error(err))
		// 类型断言
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 翻译并去除错误提示中的结构体标识
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取当前登录用户id
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(int64(currentUserID), p); err != nil {
		zap.L().Error("logic.VoteForPost() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// PostOrderListHandler 根据最新时间/分数排序查询帖子列表
// @Summary 根据最新时间/分数排序查询帖子列表
// @Description 查询帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Failure 400 {object} _ResponsePostList
// @Router /postlist [get]
func PostOrderListHandler(c *gin.Context) {
	// 1.获取参数；2.从redis查询id列表；3.根据id去数据库查询帖子详细信息
	// 初始化参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("PostOrderListHandler ShouldBindQuery failed, err:", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(&p)
	if err != nil {
		zap.L().Error("logic.PostOrderListHandler() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区查询帖子
// @Success 200 {object} _ResponsePostList
func GetCommunityPostListHandler(c *gin.Context) {
	// 初始化参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("GetCommunityPostListHandler ShouldBindQuery failed, err:", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityPostList(&p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostListHandler() failed, err:", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

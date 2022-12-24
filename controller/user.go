package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamSignUp false "注册参数"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	// 1.参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// @Summary 用户登录
// @Description 用户登录
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.ParamLogin false "登录参数"
// @Success 200 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("login with invalid param", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID,
		"user_name": user.Username,
		"token":     user.Token,
	})
}

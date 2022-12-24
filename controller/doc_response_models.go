package controller

import "bluebell/models"

// 专门用来存放接口文档用到的model

type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`
	Message string                  `json:"message"` // 提示消息
	Data    []*models.ApiPostDetail `json:"data"`    // 数据
}

type _ResponseCommunityList struct {
	Code    ResCode             `json:"code"`
	Message string              `json:"message"` // 提示消息
	Data    []*models.Community `json:"data"`    // 数据
}

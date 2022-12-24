package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetCommunityList 获取社区列表
func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id int64) (*models.Community, error) {
	return mysql.GetCommunityById(id)
}

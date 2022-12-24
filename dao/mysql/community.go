package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := "select * from community order by create_time desc"
	err = db.Select(&data, sqlStr)
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no community data")
		return nil, ErrorNoData
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetCommunityDetail 查询社区详情
func GetCommunityById(id int64) (data *models.Community, err error) {
	sqlStr := "select community_id,community_name,introduction,create_time,update_time from community where community_id=?"
	data = new(models.Community)
	err = db.Get(data, sqlStr, id)
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no community data")
		return nil, ErrorNoData
	}
	if err != nil {
		zap.L().Warn("mysql GetCommunityDetail failed, err", zap.Error(err))
		return nil, err
	}
	return data, nil
}

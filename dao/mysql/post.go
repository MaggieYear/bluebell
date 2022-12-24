package mysql

import (
	"bluebell/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into post(post_id, title, content, author_id, community_id, status) values (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID, p.Status)
	return
}

// GetPostById 查询帖子详情
func GetPostById(id int64) (data *models.Post, err error) {
	sqlStr := "select post_id,title,content,author_id,community_id,status,create_time,update_time from post where post_id=?"
	data = new(models.Post)
	err = db.Get(data, sqlStr, id)
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no Post data")
		return nil, ErrorNoData
	}
	if err != nil {
		zap.L().Warn("mysql GetPostDetail failed, err", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// GetPostList 获取帖子列表
func GetPostList(pageNum, pageSize int64) (data []*models.Post, err error) {
	sqlStr := "select post_id,title,content,author_id,community_id,status,create_time,update_time from post limit ?,?"
	// 初始化切片
	data = make([]*models.Post, 0, 2)
	err = db.Select(&data, sqlStr, (pageNum-1)*pageSize, pageSize)
	return
}

// GetPostListByIDs 根据ids获取帖子列表
func GetPostListByIDs(ids []string) (data []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,status,create_time,update_time 
	from post 
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&data, query, args...)
	return
}

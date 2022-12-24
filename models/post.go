package models

import "time"

// Post 帖子
type Post struct {
	ID          int64      `json:"id" db:"post_id"`
	AuthorID    int64      `json:"author_id" db:"author_id"`
	CommunityID int64      `json:"community_id" db:"community_id" binding:"required"`
	Status      int64      `json:"status" db:"status"`
	Title       string     `json:"title" db:"title" binding:"required"`
	Content     string     `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time  `json:"create_time" db:"create_time"`
	UpdateTime  *time.Time `json:"update_time" db:"update_time"`
}

// ApiPostDetail 帖子详情
type ApiPostDetail struct {
	AuthorName string             `json:"author_name"`
	VoteNum    int64              `json:"vote_num"`
	*Post      `json:"post"`      // 嵌入帖子结构体
	*Community `json:"community"` // 嵌入社区结构体
}

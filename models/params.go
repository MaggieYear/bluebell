package models

// ParamSignUp 注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email"`
}

// ParamLogin 登录参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票参数
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required"`                       // 帖子id
	Direction int8   `json:"direction,string" binding:"required,oneof=1 0 -1"` // 赞成票1；反对票-1
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page" example:"1"`       // 页码
	Size        int64  `json:"size" form:"size" example:"10"`      // 每页数据量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)

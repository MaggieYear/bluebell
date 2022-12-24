package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	id, _ := snowflake.GenID2()
	p.ID = int64(id)
	err = mysql.CreatePost(p)

	redis.CreatePost(p.ID, p.CommunityID)
	return err
}

// GetPostById 查询帖子详情
func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子
	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", id), zap.Error(err))
		return
	}
	// 查询作者信息
	user, err := mysql.GetUserByID(uint64(post.AuthorID))
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
		return
	}
	// 查询社区信息
	commnunity, err := mysql.GetCommunityById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityById failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName: user.Username,
		Post:       post,
		Community:  commnunity,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(pageNum, pageSize int64) (postDetails []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}

	postDetails = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {

		// 查询作者信息
		user, err := mysql.GetUserByID(uint64(post.AuthorID))
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		// 查询社区信息
		commnunity, err := mysql.GetCommunityById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityById failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post:       post,
			Community:  commnunity,
		}
		postDetails = append(postDetails, postDetail)
	}
	return
}

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}

// GetPostOrderList 根据排序获取帖子列表
func GetPostOrderList(p *models.ParamPostList) (postDetails []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetPostOrderList", zap.Any("ids", ids))
	// 根据id去数据库查询
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	// 查询每篇帖子的投票数
	votes, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	postDetails = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 查询作者信息
		user, err := mysql.GetUserByID(uint64(post.AuthorID))
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		// 查询社区信息
		commnunity, err := mysql.GetCommunityById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityById failed", zap.Int64("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum:    votes[idx],
			Post:       post,
			Community:  commnunity,
		}
		postDetails = append(postDetails, postDetail)
	}
	return
}

// GetCommunityPostList 根据社区查询帖子
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostIDsInOrder", zap.Any("ids", ids))

	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostListByIDs", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(uint64(post.AuthorID))
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityById failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum:    voteData[idx],
			Post:       post,
			Community:  community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 根据是否存在社区id判断执行逻辑
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		// 查询所有列表
		data, err = GetPostOrderList(p)
	} else {
		// 根据社区id查询列表
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

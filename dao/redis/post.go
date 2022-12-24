package redis

import (
	"bluebell/models"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 简单投票算法：《redis实战》
// 投一票+432分，86400/200-->200张赞成票可以续帖子一天
/*投票的几种情况
direction=1,
>如该用户第一次投票，投赞成票；       差值绝对值 1 +432
>若该用户之前投反对票，本次投赞成票； 差值绝对值 2 +432*2
direction=0,
>之前投赞成票，本次取消投票；         差值绝对值 1 -432
>之前投反对票，本次取消投票；         差值绝对值 1 +432
direction=-1
>如该用户第一次投票，投反对票；        差值绝对值 1 -432
>若该用户之前投赞成票，本次投反对票；   差值绝对值 2 -432*2

投票的限制：
>每个用户每个帖子最多只能投一票；
>每个帖子发布一周后允许用户投票，超过一周停止投票；
>到期之后将redis中保存的赞成票及反对票存储到数据表中；
>到期后删除redis的key
*/

const (
	oneWeekSeconds = 7 * 24 * 3600
	scorePerVote   = 432 // 每一票的分数值
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已截止")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// 发帖
func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScore), redis.Z{
		Score:  0,
		Member: postID,
	})
	// 把帖子的id加到社区set里
	// cKey示例：bluebell:community:2
	cKey := getRedisKey(KeyCommunitySet + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)

	_, err := pipeline.Exec()
	return err
}

// VoteForPost 为帖子投票
func VoteForPost(userID, postID string, value float64) error {
	// 1.判断投票限制
	// redis获取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekSeconds {
		return ErrVoteTimeExpire
	}

	// 步骤2和3需要放到一个pipeline事务中操作
	// 2.更新帖子分数
	// 查询之前的投票记录
	vote := client.ZScore(getRedisKey(KeyPostVotedPrefix+postID), userID).Val()

	// 如果本次投票值跟上次投票一致，则不允许重复投票
	if value == vote {
		return ErrVoteRepeated
	}
	var dir float64
	if value > vote {
		dir = 1
	} else {
		dir = -1
	}
	// 绝对值，计算两次投票的差值
	diff := math.Abs(vote - value)
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScore), dir*diff*scorePerVote, postID)

	// 3.记录用户为该帖子投票的数据
	if value == 0 { // 取消投票
		pipeline.ZRem(getRedisKey(KeyPostVotedPrefix+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}

// GetPostIDsInOrder 根据排序获取帖子id列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 默认按时间倒序
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		// 按分数倒序
		key = getRedisKey(KeyPostScore)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的赞成票数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	keys := make([]string, 0, len(ids))
	pipeline := client.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(key, "1", "1")
		keys = append(keys, key)
	}

	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区查询帖子
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// zinterstore：把社区的帖子set与帖子分数的zset联合生成一个新的zset
	orderKey := getRedisKey(KeyPostTime) // bluebell:post:time
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScore) // orderKey= bluebell:post:score
	}
	// 社区的key，示例：bluebell:community:2
	cKey := getRedisKey(KeyCommunitySet + strconv.Itoa(int(p.CommunityID)))

	// 缓存key，减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID)) // bluebell:post:score2

	// 不存在key，需要计算
	if client.Exists(key).Val() < 1 {
		pipeline := client.Pipeline()
		// ZInterStore计算
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		// 设置超时时间
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// key存在则直接获取
	return getIDsFormKey(key, p.Page, p.Size)
}

package redis

// redis key
// 使用命名空间方式，方便查询和拆分
const (
	Prefix             = "bluebell:"
	KeyPostTime        = "post:time"   // ZSet;帖子及发帖时间
	KeyPostScore       = "post:score"  // ZSet;帖子及投票的分数
	KeyPostVotedPrefix = "post:voted:" // ZSet;记录用户及投票类型
	KeyCommunitySet    = "community:"  //set, 保存每个社区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}

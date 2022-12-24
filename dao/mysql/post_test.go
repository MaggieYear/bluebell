package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

// 测试数据库地址
func initdb() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "123456",
		DbName:       "go_test_db",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := TestInit(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {

	initdb()
	post := new(models.Post)
	post.ID = 10
	post.AuthorID = 123
	post.CommunityID = 1
	post.Title = "test"
	post.Content = "测试内容"
	// post = models.Post{
	// 	ID:          10,
	// 	AuthorID:    123,
	// 	CommunityID: 1,
	// 	Title:       "test",
	// 	Content:     "测试内容",
	// }
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost insert record failed:%v\n", err)
	}
	t.Logf("CreatePost insert record success.")
}

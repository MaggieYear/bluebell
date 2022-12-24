package controller

import (
	"net/http"
	"net/http/httptest"

	"bytes"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 用于模拟发起http请求的的测试用例
func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	url := "/api/v1/post"
	router.POST(url, CreatePostHandler)

	body := `{
		"title": "测试标题18",
		"content": "测试内容3sdwerwer3434s34",
		"community_id": 2
	}`

	// 发送请求
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	// 响应记录对象
	w := httptest.NewRecorder()

	// 将响应结果记录到w
	router.ServeHTTP(w, req)
	// 断言响应返回码是否跟200一致，结果写入到t
	assert.Equal(t, 200, w.Code)

	// 方法一：判断响应内容是否存在指定字符串
	//assert.Contains(t, w.Body.String(), "需要登录")

	// 方法二：将响应内容反序列化到ResponseData，然后判断字段是否符合预期
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}

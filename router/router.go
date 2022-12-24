package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "bluebell/docs"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	// 127.0.0.1:8080/signup
	// {"user":"Tim", "pwd":"123456", "re_password":"123456"}
	v1.POST("/signup", controller.SignUpHandler)

	v1.POST("/login", controller.LoginHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1.GET("/community", controller.CommunityHandler)
	// 根据id获取社区详情
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	// 查看帖子详情
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	// 获取帖子列表
	v1.GET("/posts", controller.GetPostListHandler)
	// 根据最新时间/分数排序查询帖子列表
	v1.GET("/postlist", controller.PostOrderListHandler)

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// 发布帖子
		v1.POST("/post", controller.CreatePostHandler)
		// 投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Go Web开发通用脚手架模板

// @titile bluebell
// @version 1.0
// @description 本项目提供用户注册、用户登录、发帖及查看帖子、帖子投票等功能；
// @termsOfService http://swagger.io/terms/
// @contact.name kiki
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8081
// @BasePath /api/v1
func main() {

	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2.初始化日志
	if err := logger.Init(settings.LogSettings, settings.AppSettings.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	// 3.初始化Mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 4.初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init2(settings.AppSettings.StartTime, settings.AppSettings.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}

	// 5.注册路由
	r := router.SetupRouter()
	r.Run(fmt.Sprintf("%s:%d", settings.AppSettings.Host, settings.AppSettings.Port))

	// 6.启动服务（优雅关机）
	shutdown(r)
}

func shutdown(r http.Handler) {
	//启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", settings.AppSettings.Host, settings.AppSettings.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	// 终端Ctr + C 优雅关机

	quit := make(chan os.Signal, 1)                      // 创建一个接收信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	// 阻塞，当接收到上述两种信号时才会往下执行
	<-quit

	zap.L().Info("Shutdown Server ...")
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
		log.Fatal("Server Shutdown: ", err)
	}
	zap.L().Info("Server exiting")
	log.Println("Server exiting")
}

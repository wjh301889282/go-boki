package main

import (
	"context"
	"errors"
	"exchangeapp/config"
	_ "exchangeapp/docs" // main 文件中导入 docs 包
	"exchangeapp/gorm"
	"exchangeapp/router"
	"exchangeapp/websorket"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var once sync.Once

func main() {

	once.Do(func() {
		// 初始化配置文件
		config.InitConfig()
		// 初始化 GORM 数据库连接
		gorm.InitGORM()
		fmt.Println("加载成功配置环境")

	})
	// 设置路由
	r := router.SetupRouter()

	// 获取应用的端口号，如果没有指定，则默认使用 ":8080"
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}

	// 创建 HTTP 服务器并指定端口和处理器
	srv := &http.Server{
		Addr:    port, // 服务器监听的端口
		Handler: r,    // 路由处理器
	}

	// 启动 WebSocket 服务
	go func() {
		if err := websorket.Initwebsorket(); err != nil {
			log.Fatalf("WebSocket 启动失败: %v", err)
		}
		log.Println("WebSocket 服务已启动")
	}()

	// 启动一个新的 goroutine 来运行服务器
	go func() {
		// 启动 HTTP 服务，监听指定的端口
		// 如果出现非预期错误，则日志输出并退出
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 创建一个接收系统中断信号的通道
	quit := make(chan os.Signal, 1)
	// 监听系统中断信号（如 Ctrl+C）
	signal.Notify(quit, os.Interrupt)
	// 阻塞，等待中断信号
	<-quit
	log.Println("Shutdown Server ...")

	// 设置一个 5 秒的超时上下文，用于优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用服务器的 Shutdown 方法来优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

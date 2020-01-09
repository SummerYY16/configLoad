package g

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	log "configLoad/toolkits/logger"
)

var stopWebChannel = make(chan string)

// WebAPIRestart 优雅重启web接口
func WebAPIRestart() {
	log.Println("web服务开始重启")
	WebAPIStop()
	<-stopWebChannel
	WebAPIStart()
}

// WebAPIStart 启动web接口
func WebAPIStart() {
	go initWebAPI()
}

// WebAPIStop 关闭web接口
func WebAPIStop() {
	stopWebChannel <- "stop"
}

// InitWebAPI 初始化api接口
func initWebAPI() {
	appConfig := AppConfigMgrHandler.Config.Load().(*AppConfig)
	host := appConfig.Apiserver.Host
	port := appConfig.Apiserver.Port

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    http.StatusOK,
			"message": "pong",
		})
	})

	address := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}
	log.Println("正在启动web服务")
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("webapi服务启动报错: %s\n", err)
		}
	}()

	// 接受一个退出信号，退出web服务
	<-stopWebChannel
	log.Println("web服务开始关闭 ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("web服务关闭出错:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("web服务已关闭")
		stopWebChannel <- "done"
	}
}

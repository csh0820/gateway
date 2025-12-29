package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/csh0820/gateway/config"
	"github.com/csh0820/gateway/internal/gateway"
	"github.com/csh0820/gateway/pkg/etcd"

	"github.com/gin-gonic/gin"
)

func main() {
	// get config
	cfg := config.GetConfig()

	ctx := context.Background()

	cli := etcd.NewEtcdClient()
	reg := etcd.New(cli)

	gatewayHandler := gateway.NewGatewayHandler(reg)

	if cfg.GatewayMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// start router
	router := gin.New()

	// 使用中间件
	router.Use(gin.Recovery())

	// 注册路由
	router.Any("/*path", gateway.HandleRequest)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    cfg.GatewayAddress,
		Handler: router,
	}

	log.Println("gateway server start...")
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("gateway server start failed:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// stop watch
	watch.Stop()

	// 设置关闭超时
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalln("gateway server shutdown failed:", err)
	} else {
		log.Println("gateway server shutdown!")
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"demo/internal/config"
	"demo/internal/handler"
	"demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "/Users/123jiaru/Desktop/project/hashlink/mpc_go/etc/demo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 设置优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Println("🔗 BSC监控服务已集成启动")
	fmt.Println("📤 TokenEvent将发送到Mock Kafka")

	// 在独立的goroutine中启动服务器
	go func() {
		server.Start()
	}()

	// 等待退出信号
	<-quit
	fmt.Println("\n🛑 收到退出信号，正在优雅关闭服务...")

	// 停止监控服务
	ctx.StopMonitor()

	fmt.Println("✅ 服务已安全退出")
}

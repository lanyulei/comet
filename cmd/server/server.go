package server

import (
	"context"
	"fmt"
	"github.com/lanyulei/comet/common/router"
	"github.com/lanyulei/comet/pkg/config"
	"github.com/lanyulei/comet/pkg/db"
	"github.com/lanyulei/comet/pkg/logger"
	"github.com/lanyulei/comet/pkg/redis"
	"github.com/lanyulei/comet/pkg/tools"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "start API server",
		Example:      "comet server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
}

func setup() {
	// 加载配置文件
	config.Setup(configYml)

	// 日志配置
	logger.Setup()

	// 数据库配置
	db.Setup()

	// Redis 链接
	//redis.Setup()
}

func run() (err error) {
	if viper.GetString("server.mode") == config.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := gin.Default()
	router.Setup(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")),
		Handler: r,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func(cancel context.CancelFunc) {
		cancel()
		// 关闭 redis 连接
		redis.StopChRedis() <- struct{}{}
	}(cancel)

	go func() {
		// 服务连接
		if viper.GetBool("ssl.enable") {
			if err := srv.ListenAndServeTLS(viper.GetString("ssl.pem"), viper.GetString("ssl.key")); err != nil && err != http.ErrServerClosed {
				logger.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("listen: ", err)
			}
		}
	}()

	fmt.Println("\nServer run at:")
	fmt.Printf("-  Local:   http://localhost:%d/ \r\n", viper.GetInt("server.port"))
	fmt.Printf("-  Network: http://%s:%d/ \r\n", tools.GetLocalHost(), viper.GetInt("server.port"))
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n\n", time.Now().Format("2006-01-02 15:04:05.000"))
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", time.Now().Format("2006-01-02 15:04:05"))

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Info("Server exiting")

	return nil
}

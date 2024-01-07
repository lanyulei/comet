package server

import (
	"comet/common/router"
	"comet/pkg/config"
	"comet/pkg/kubernetes/client"
	"comet/pkg/kubernetes/request"
	"comet/pkg/tools"
	"context"
	"fmt"
	"github.com/lanyulei/toolkit/logger"
	"github.com/lanyulei/toolkit/redis"
	"github.com/lanyulei/toolkit/service"
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
	logger.Setup(
		viper.GetString(`log.level`),
		viper.GetString(`log.path`),
		viper.GetInt(`log.maxsize`),
		viper.GetBool(`log.localtime`),
		viper.GetBool(`log.compress`),
		viper.GetBool(`log.console`),
		map[string]interface{}{},
	)

	// 数据库配置
	//db.Setup(
	//	viper.GetString("db.type"),
	//	viper.GetString("db.dsn"),
	//	viper.GetInt("db.maxIdleConn"),
	//	viper.GetInt("db.maxOpenConn"),
	//	viper.GetInt("db.connMaxLifetime"),
	//)

	// Redis 链接
	redis.Setup(viper.GetString("redis.host"), viper.GetString("redis.password"), viper.GetInt("redis.port"), viper.GetInt("redis.db"))

	// 加载 kubernetes 配置
	client.NewClients()
}

func run() (err error) {
	if viper.GetString("server.mode") == config.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := gin.Default()
	router.Setup(r)

	h := request.BuildHandlerChain(r)
	if h == nil {
		h = r
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")),
		Handler: h,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func(cancel context.CancelFunc) {
		cancel()
		// 关闭 redis 连接
		redis.StopChRedis()
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

	// 服务注册
	service.Register(viper.GetInt("server.port"), prefix)

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

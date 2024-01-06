package migrate

import (
	"fmt"
	"github.com/lanyulei/toolkit/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configYml string
	data      bool
	sync      bool
	clean     string
	list      bool
	generate  bool
	StartCmd  = &cobra.Command{
		Use:          "migrate",
		Short:        "synchronous data structure",
		Example:      "comet migrate -c config/settings.yml",
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
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "specify the profile to start the service")
	StartCmd.PersistentFlags().BoolVarP(&data, "data", "d", false, "synchronous data")
	StartCmd.PersistentFlags().BoolVarP(&sync, "sync", "s", true, "synchronize data structures")
	StartCmd.PersistentFlags().StringVarP(&clean, "clear", "e", "", "clear migration records")
	StartCmd.PersistentFlags().BoolVarP(&list, "list", "l", false, "list of migration records")
	StartCmd.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "generate sql files for synchronized data")
}

func setup() {
	// 加载配置文件
	viper.SetConfigFile(configYml) // 指定配置文件
	err := viper.ReadInConfig()    // 读取配置信息
	if err != nil {                // 读取配置信息失败
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

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
	//db.Setup()
}

func run() (err error) {
	if sync {
		// synchronized data structure
		autoMigrate(data)
	} else if clean != "" {
		// clear migration records
		clearMigrate(clean)
	} else if list {
		// view a migration list
		listMigrate()
	} else if generate {
		// generate migration sql file
		generateSQL()
	} else {
		fmt.Println("Enter the following command to view the help information:")
		fmt.Println("  comet migrate -h")
	}
	return
}

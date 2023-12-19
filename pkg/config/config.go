package config

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
  @Desc :
*/

func Setup(configYml string) {
	// loading configuration files
	viper.SetConfigFile(configYml) // specify the profile
	err := viper.ReadInConfig()    // read configuration information
	if err != nil {                // failed to read configuration information
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

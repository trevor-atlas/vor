package system

import (
	"github.com/spf13/viper"
)

func GetString(envName string) string {
	res := viper.GetString(envName)
	if res == "" {
		res = viper.GetString("global." + envName)
	}
	return res
}

func GetBool(envName string) bool {
	res := viper.GetBool(envName)
	if res == false {
		res = viper.GetBool("global." + envName)
	}
	return res
}

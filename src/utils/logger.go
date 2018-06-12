package logger

import (
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func Debug(message string) {
	isDev := viper.GetBool("VOR_IS_DEVELOPMENT_MODE")
	if isDev {
		color.Cyan("dev: " + message)
	}
}

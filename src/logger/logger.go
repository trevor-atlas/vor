package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func Debug(message ...interface{}) {
	isDev := viper.GetBool("devmode")
	if isDev {
		color.Cyan("DEBUG:")
		fmt.Println(message)
	}
}

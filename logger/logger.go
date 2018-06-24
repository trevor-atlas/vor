package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// Log a centralized utility for logging
func Log(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
}

func Error(message string, rest ...interface{}) {
	color.Red(message, rest...)
}

func Debug(message ...interface{}) {
	isDev := viper.GetBool("devmode")
	if isDev {
		color.Cyan("DEBUG:")
		fmt.Println(message...)
	}
}

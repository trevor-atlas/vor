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

func Debug(format string, message ...interface{}) {
	isDev := viper.GetBool("devmode")
	yellow := color.New(color.FgHiYellow).SprintFunc()
	if isDev {
		fmt.Printf(yellow("DEBUG: ") + format, message...)
		fmt.Println()
	}
}

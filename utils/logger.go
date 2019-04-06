package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"trevoratlas.com/vor/formatters"
)

func Info(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
}

func Debug(format string, rest ...interface{}) {
	devmode := viper.GetBool("devmode")
	if devmode {
		fmt.Println(formatters.YELLOW("DEBUG: " + format, rest...))
	}
}

func Error(format string, rest ...interface{}) {
	fmt.Println(formatters.RED(format, rest...))
}

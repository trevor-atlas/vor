package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/formatters"
)

type StructuredLogger interface {
	Info(format string, rest ...interface{})
	Debug(format string, rest ...interface{})
	Error(format string, rest ...interface{})
}

type _logger struct {}
// Ensure Logger implements the StructuredLogger interface
var logger StructuredLogger = &_logger{}

func Info(message string, rest ...interface{}) { logger.Info(message, rest)}
func (l *_logger) Info(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
}

func Debug(format string, rest ...interface{}) { logger.Debug(format, rest) }
func (l *_logger) Debug(format string, rest ...interface{}) {
	devmode := viper.GetBool("devmode")
	if devmode {
		fmt.Println(formatters.YELLOW("DEBUG: " + format, rest...))
	}
}

func Error(format string, rest ...interface{}) { logger.Error(format, rest) }
func (l *_logger) Error(format string, rest ...interface{}) {
	fmt.Println(formatters.RED(format, rest...))
}

package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/trevor-atlas/vor/system"
	"sync"
)

type StructuredLogger interface {
	Info(format string, rest ...interface{})
	Debug(format string, rest ...interface{})
	Error(format string, rest ...interface{})
}

type Logger struct {
	isDev     bool
	fmtYellow func(a ...interface{}) string
	fmtRed    func(a ...interface{}) string
}

// Ensure Logger implements the StructuredLogger interface
var _ StructuredLogger = (*Logger)(nil)

// Use sync to create and return a singleton of the logger
var once sync.Once
var logger *Logger

func New() *Logger {
	once.Do(func() {
		logger = &Logger{}
		isDev := system.Get.Bool("devmode")
		logger.isDev = isDev
		logger.fmtYellow = color.New(color.FgHiYellow).SprintFunc()
		logger.fmtRed = color.New(color.FgHiRed).SprintFunc()
	})
	return logger
}

func (l *Logger) Info(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
	fmt.Println()
}

func (l *Logger) Debug(format string, rest ...interface{}) {
	if l.isDev {
		fmt.Printf(l.fmtYellow("DEBUG: ")+format, rest...)
		fmt.Println()
	}
}

func (l *Logger) Error(format string, rest ...interface{}) {
	fmt.Printf(format, rest...)
	fmt.Println()
}

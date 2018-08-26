package formatters

import "github.com/fatih/color"

var (
	CYAN    func(format string, a ...interface{}) string
	BLUE    func(format string, a ...interface{}) string
	YELLOW  func(format string, a ...interface{}) string
	MAGENTA func(format string, a ...interface{}) string
	RED     func(format string, a ...interface{}) string
)

type StringFormatter interface {
	Init()
}

type DefaultStringFormatter struct{}

func (f *DefaultStringFormatter) Init() {
	CYAN = color.New(color.FgHiCyan).SprintfFunc()
	BLUE = color.New(color.FgHiBlue).SprintfFunc()
	YELLOW = color.New(color.FgHiYellow).SprintfFunc()
	MAGENTA = color.New(color.FgHiMagenta).SprintfFunc()
	RED = color.New(color.FgHiRed).SprintfFunc()
}

func Init(formatter StringFormatter) {
	formatter.Init()
}

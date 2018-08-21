package formatters

import "github.com/fatih/color"

var (
	CYAN    func(a ...interface{}) string
	BLUE    func(a ...interface{}) string
	YELLOW  func(a ...interface{}) string
	MAGENTA func(a ...interface{}) string
	RED     func(a ...interface{}) string
)

type StringFormatter interface {
	Init()
}

type DefaultStringFormatter struct{}

func (f *DefaultStringFormatter) Init() {
	CYAN = color.New(color.FgHiCyan).SprintFunc()
	BLUE = color.New(color.FgHiBlue).SprintFunc()
	YELLOW = color.New(color.FgHiYellow).SprintFunc()
	MAGENTA = color.New(color.FgHiMagenta).SprintFunc()
	RED = color.New(color.FgHiRed).SprintFunc()
}

func Init(formatter StringFormatter) {
	formatter.Init()
}

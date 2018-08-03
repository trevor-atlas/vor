package system

import (
	"github.com/spf13/viper"
)

type Envloader interface {
	String(string) string
	Bool(string) bool
}

// ensure Get implements the Envloader interface
var _ Envloader = (*Get)(nil)

type Get struct{}

func NewENVGetter() *Get {
	return new(Get)
}

func (g *Get) String(envName string) string {
	res := viper.GetString(envName)
	if res == "" {
		res = viper.GetString("global." + envName)
	}
	return res
}

func (g *Get) Bool(envName string) bool {
	res := viper.GetBool(envName)
	if res == false {
		res = viper.GetBool("global." + envName)
	}
	return res
}

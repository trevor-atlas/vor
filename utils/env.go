package utils

import (
	"github.com/spf13/viper"
)

type Envloader interface {
	String(string) string
	Bool(string) bool
}

type ENV struct{
	get Envloader
}

func (h *ENV) String(env string) (string) {
	res := viper.GetString(env)
	if res == "" {
		res = viper.GetString("global." + env)
	}
	return res
}

func (h *ENV) Bool(env string) (string) {
	res := viper.GetString(env)
	if res == "" {
		res = viper.GetString("global." + env)
	}
	return res
}




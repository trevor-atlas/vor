package env

import (
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/system"
)

var (
	JIRA_USERNAME string
	JIRA_APIKEY   string
	JIRA_ORGNAME  string
)

type EnvironmentLoader interface {
	Init()
}

type DefaultLoader struct{}

func (l *DefaultLoader) Init() {
	JIRA_USERNAME = viper.GetString("jira.username")
	if JIRA_USERNAME == "" {
		system.Exit("jira.username config not found.")
	}
	JIRA_APIKEY = viper.GetString("jira.apikey")
	if JIRA_APIKEY == "" {
		system.Exit("jira.apikey config not found.")
	}
	JIRA_ORGNAME = viper.GetString("jira.orgname")
	if JIRA_ORGNAME == "" {
		system.Exit("jira.orgname config not found.")
	}
}

func Init(loader EnvironmentLoader) {
	loader.Init()
}

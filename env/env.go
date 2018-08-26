package env

import (
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/system"
)

var (
	JIRA_USERNAME     string
	JIRA_APIKEY       string
	JIRA_ORGNAME      string
	PULL_REQUEST_BASE string
	GITHUB_APIKEY     string
	DEV_MODE          bool
)

type EnvironmentLoader interface {
	Init()
}

type DefaultLoader struct{}

// this will need to change when there are generic issue
// providers and VCS hosts
func (l *DefaultLoader) Init() {
	DEV_MODE = viper.GetBool("devmode")
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
	GITHUB_APIKEY = viper.GetString("github.apikey")
	if GITHUB_APIKEY == "" {
		system.Exit("No github API key found in vor config (github.apikey)")
	}
	PULL_REQUEST_BASE = viper.GetString("git.pull-request-base")
	if PULL_REQUEST_BASE == "" {
		PULL_REQUEST_BASE = "master"
	}
}

func Init(loader EnvironmentLoader) {
	loader.Init()
}

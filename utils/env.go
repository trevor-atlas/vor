package utils

import (
	"github.com/spf13/viper"
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

const PLACEHOLDER = "PLACEHOLDER"
var LOAD_ERROR = false

func loadEnvOrFail(name string) string {
	env := viper.GetString(name)
	if env == "" {
		Error("%s config not found.", name)
		LOAD_ERROR = true
	}
	if env == PLACEHOLDER {
		Error("%s is not configured, you need to set it in your config", name)
		LOAD_ERROR = true
	}
	return env
}

func loadEnvOrFallback(name string, fallback string) string {
	env := viper.GetString(name)
	if env == "" {
		return fallback
	}
	if env == PLACEHOLDER {
		return fallback
	}
	return env
}

// this will need to change when there are generic issue
// providers and VCS hosts
func (l *DefaultLoader) Init() {
	DEV_MODE = viper.GetBool("devmode")

	JIRA_USERNAME = loadEnvOrFail("jira.username")
	JIRA_APIKEY = loadEnvOrFail("jira.apikey")
	JIRA_ORGNAME = loadEnvOrFail("jira.orgname")
	GITHUB_APIKEY = loadEnvOrFail("github.apikey")

	PULL_REQUEST_BASE = loadEnvOrFallback("git.pull-request-base", "master")
	if LOAD_ERROR {
		Exit("%s is missing necessary values or invalid. Exiting.", viper.ConfigFileUsed())
	}
}

func Init(loader EnvironmentLoader) {
	loader.Init()
}

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"trevoratlas.com/vor/formatters"

	"trevoratlas.com/vor/utils"
)

var (
	CONFIG_FILE string
	PLACEHOLDER_CONFIG = []byte(`# your name
author: PLACEHOLDER
# you@yourdomain.xyz
email: PLACEHOLDER
jira:
    # <your company name, usually contained in the url or your jira install>
    orgname: PLACEHOLDER
    # <your jira username (sometimes an email)>
    username: PLACEHOLDER
    # <your jira api key from id.atlassian.net>
    apikey: PLACEHOLDER
github:
    # <the owner of the repository>
    owner: PLACEHOLDER
    # <your github api key (get this from github.com/settings/tokens)>
    apikey: PLACEHOLDER`)
)

var rootCmd = &cobra.Command{
	Use:   "vor",
	Short: "Vör – make Github and Jira easy",
	Long: `
                  ___          ___
      ___        /\  \        /\  \
     /\  \      /::\  \      /::\  \
     \:\  \    /:/\:\  \    /:/\:\__\
      \:\  \  /:/  \:\  \  /:/ /:/  /
  ___  \:\__\/:/__/ \:\__\/:/_/:/__/___
 /\  \ |:|  |\:\  \ /:/  /\:\/:::::/  /
 \:\  \|:|  | \:\  /:/  /  \::/~~/~~~~
  \:\__|:|__|  \:\/:/  /    \:\~~\
   \::::/__/    \::/  /      \:\__\
    ~~~~         \/__/        \/__/

 This program comes with ABSOLUTELY NO WARRANTY; This is free software, and you are welcome to redistribute it.
 Vör – A fast and flexible commandline tool for working with Github and Jira`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.Init(&utils.DefaultLoader{})
	},
	Run: func(cmd *cobra.Command, args []string){},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	formatters.Init(&formatters.DefaultStringFormatter{})
	cobra.OnInitialize(initConfig)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//  If Persistent flags are defined here they are global
	rootCmd.PersistentFlags().StringVar(&CONFIG_FILE, "config", "", "config file (default is $HOME/.vor, or the current directory)")
	viper.SetDefault("devmode", true)
	viper.SetDefault("git.branchtemplate", "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}")
	viper.SetDefault("git.path", "/usr/local/bin/git")
	viper.SetDefault("git.pull-request-base", "master")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("vor")

	if CONFIG_FILE != "" {
		viper.SetConfigFile(CONFIG_FILE) // Use config file from the flag.
		loadConfig()
		return
	}

	home, homeErr := homedir.Dir()
	if homeErr != nil {
		fmt.Println(homeErr)
		utils.Exit("vor encountered an error attempting to read from the filesystem")
	}


	configPath, walkErr := utils.WalkUpFS(filepath.Base(home))
	if walkErr != nil {
		fmt.Println(walkErr)
		utils.Exit("vor encountered an error attempting to read from the filesystem")
	}

	if configPath != "" {
		utils.Debug("Located a local config file at '%s'", configPath)
		path, _ := filepath.Split(configPath)
		viper.AddConfigPath(path)
		loadConfig()
		return
	}

	XDGConfig := utils.CreateXDGConfig(PLACEHOLDER_CONFIG)
	if XDGConfig != "" {
		utils.Debug("loaded the XDG config file at '%s'", XDGConfig)
		viper.AddConfigPath(XDGConfig)
		loadConfig()
		return
	}

	utils.Exit("Vor could not find a local config file")
}

func loadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		utils.Exit("Vor could not find a local config file")
	}
}

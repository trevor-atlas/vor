package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vor",
	Short: "Vör – make Github and Jira easy",
	Long:  `
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

 Vör Copyright (C) 2018  Trevor Atlas

 This program comes with ABSOLUTELY NO WARRANTY; for details type 'show l'. This is free software, and you are welcome to redistribute it.
 Vör – A fast and flexible commandline tool for working with Github and Jira`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vor.yaml, or the current directory)")

	// Default environment configs
	viper.SetDefault("devmode", false)

	viper.SetDefault("global.branchtemplate", "{projectname}/{jira-issue-number}/{jira-issue-type}/{jira-issue-title}")
	viper.SetDefault("branchtemplate", "")

	viper.SetDefault("projectname", "")

	viper.SetDefault("global.jira.orgname", "")
	viper.SetDefault("jira.orgname", "")

	viper.SetDefault("global.jira.apikey", "")
	viper.SetDefault("global.jira.username", "")

	viper.SetDefault("jira.username", "")
	viper.SetDefault("jira.apikey", "")

	viper.SetDefault("global.github.apikey", "")
	viper.SetDefault("github.apikey", "")

	viper.SetDefault("github.owner", "")

	viper.SetDefault("global.git.path", "/usr/local/bin/git")
	viper.SetDefault("git.path", "")

	viper.SetDefault("global.git.pull-request-base", "master")
	viper.SetDefault("git.pull-request-base", "")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// Search config in home directory with name ".vor"
	viper.SetConfigType("yaml")
	viper.SetConfigName(".vor")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		color.Red("Vor could not find a local config file, this can cause problems and is not recommended\n")
		fmt.Println(err)
		// os.Exit(1)
	}
}

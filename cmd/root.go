package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/env"
	"github.com/trevor-atlas/vor/formatters"
	"github.com/trevor-atlas/vor/system"
	"os"
	"strings"
)

var (
	CONFIG_FILE string
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
		env.Init(&env.DefaultLoader{})
		formatters.Init(&formatters.DefaultStringFormatter{})
	},
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
	//  If Persistent flags are defined here they are global
	rootCmd.PersistentFlags().StringVar(&CONFIG_FILE, "config", "", "config file (default is $HOME/.vor, or the current directory)")
	viper.SetDefault("devmode", false)
	viper.SetDefault("branchtemplate", "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}")
	viper.SetDefault("git.path", "/usr/local/bin/git")
	viper.SetDefault("git.pull-request-base", "master")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("vor")
	viper.SetConfigFile(".vor")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	if CONFIG_FILE != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CONFIG_FILE)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			system.Exit("vor encountered an error attempting to read from the filesystem")
		}
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		color.Red("Vor could not find a local config file, this can cause problems and is not recommended\n")
		fmt.Println(err)
	}
}

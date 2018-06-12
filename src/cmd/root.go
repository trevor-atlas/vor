package commands

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vor",
	Short: "Vör – make Github and Jira easy",
	Long:  `Vör is a fast and flexible commandline tool for working with Github and Jira`,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vor)")
	rootCmd.PersistentFlags().StringP("author", "", "<EMAIL ADDRESS>", "Author name for commit attribution")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	// Default environment configs
	viper.SetDefault("VOR_IS_DEVELOPMENT_MODE", false)
	viper.SetDefault("VOR_BRANCH_TEMPLATE", "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}")
	viper.SetDefault("VOR_GIT_PATH", "/usr/local/bin/git")
	viper.SetDefault("JIRA_API_KEY", "")
	viper.SetDefault("GITHUB_API_KEY", "")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	// Search config in home directory with name "vor" (without extension).
	viper.SetConfigType("yaml")
	viper.SetConfigName("vor")
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
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
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

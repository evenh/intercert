package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	DefaultIntercertDir = constructDir()
)

var rootCmd = &cobra.Command{
	Use:     "intercert",
	Short:   "intercert - Let's Encrypt on LAN",
	Long:    `Fetches ACME certificates in locked down environments by using DNS validation`,
	Version: `X.X.X`,
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Config file name
	viper.SetConfigName("config")
	// *sigh* YAML is probably best for this app
	viper.SetConfigType("yaml")

	// Search these paths for config file
	viper.AddConfigPath(DefaultIntercertDir)

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		PrintErrorAndExit(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintErrorAndExit(err)
	}
}

func constructDir() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		PrintErrorAndExit(err)
	}

	return home + "/.intercert"
}

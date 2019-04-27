package cmd

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// DefaultIntercertDir is the directory for where to store all data related to intercert
	DefaultIntercertDir = constructDir()
	// Version is the human readable version of intercert
	Version string
	// Commit is the git commit used to build this version
	Commit string
)

var rootCmd = &cobra.Command{
	Use:   "intercert",
	Short: "intercert - Let's Encrypt on LAN",
	Long:  `Fetches ACME certificates in locked down environments by using DNS validation`,
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
		errMsg := fmt.Sprintf("Failed to read configuration: %v", err)
		PrintErrorAndExit(errors.New(errMsg))
	}
}

// Execute the cmd chain
func Execute(version string, commit string) {
	Version = version
	Commit = commit

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

// UserAgent returns a formatted user agent (UA) string
func UserAgent() string {
	return fmt.Sprintf("intercert v%s (%s); %s-%s", Version, Commit, runtime.GOOS, runtime.GOARCH)
}

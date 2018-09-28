package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var DefaultIntercertDir = constructDir()


var rootCmd = &cobra.Command{
	Use:   "intercert",
	Short: "intercert - Let's Encrypt on LAN",
	Long:  `Fetches ACME certificates in locked down environments by using DNS validation`,
	Version: `X.X.X`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`intercert - Let's Encrypt on LAN

See --help for usage instructions`)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Server configuration flags
	viper.SetDefault("server.domains", []string{"www.example.int", "sub.example.int"})
	viper.SetDefault("server.agree", false)
	viper.SetDefault("server.storage", DefaultIntercertDir + "/server-data")
}

func initConfig() {
	// Config file name
	viper.SetConfigName("config")
	// *sigh* YAML is probably best for this app
	viper.SetConfigType("yaml")

	// Search these paths for config file
	viper.AddConfigPath(DefaultIntercertDir)

	// Support dynamic reload
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Configuration file changed:", e.Name)
	})

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read configuration :", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func constructDir() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home + "/.intercert"
}

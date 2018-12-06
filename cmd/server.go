package cmd

import (
	"errors"
	"fmt"
	"github.com/evenh/intercert/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// Flags
	serverCmd.PersistentFlags().IntP("port", "p", 6300, "The port the server will listen on")
	serverCmd.PersistentFlags().Bool("agree", false, "Whether you accept the ACME providers Terms of Service. Must be set to true in order to start server.")
	serverCmd.PersistentFlags().StringP("directory", "D", "https://acme-staging-v02.api.letsencrypt.org/directory", "URL to ACME directory")
	serverCmd.PersistentFlags().StringP("dns-provider", "P", "", "The DNS provider to use - see docs for valid providers.")
	serverCmd.PersistentFlags().StringSliceP("domains", "d", nil, "Domains to whitelist")
	serverCmd.PersistentFlags().StringP("email", "e", "", "The email to register with the ACME provider")
	serverCmd.PersistentFlags().StringP("storage", "s", DefaultIntercertDir+"/server-data", "The place to store certificates and other data")

	// Load serverCmd.PersistentFlags() values from config
	bindPrefixedFlags(serverCmd, "server", "port", "agree", "directory", "dns-provider", "domains", "storage", "email")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "The server doing the actual work",
	Long:  `Start the server component, doing the interaction with the ACME server and connected clients`,
	Run: func(cmd *cobra.Command, args []string) {

		c := ServerConfig{
			Port:        viper.GetInt("server.port"),
			Agree:       viper.GetBool("server.agree"),
			Directory:   viper.GetString("server.directory"),
			DnsProvider: viper.GetString("server.dns-provider"),
			Domains:     viper.GetStringSlice("server.domains"),
			Email:       viper.GetString("server.email"),
			Storage:     viper.GetString("server.storage"),
		}

		if !c.Agree {
			PrintErrorAndExit(errors.New("the ACME ToS must be agreed to"))
		}

		fmt.Printf("Listening on port %v\n", c.Port)
		server.StartServer(c.Port)
	},
}

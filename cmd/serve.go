package cmd

import (
	"errors"

	"github.com/evenh/intercert/config"
	"github.com/evenh/intercert/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xenolf/lego/log"
)

func init() {
	rootCmd.AddCommand(serveCmd)

	// Flags
	serveCmd.PersistentFlags().IntP("port", "p", 6300, "The port the server will listen on")
	serveCmd.PersistentFlags().Bool("agree", false, "Whether you accept the ACME providers Terms of Service. Must be set to true in order to start server.")
	serveCmd.PersistentFlags().StringP("directory", "D", "https://acme-staging-v02.api.letsencrypt.org/directory", "URL to ACME directory")
	serveCmd.PersistentFlags().StringP("dns-provider", "P", "", "The DNS provider to use - see docs for valid providers.")
	serveCmd.PersistentFlags().StringSliceP("domains", "d", nil, "Domains to whitelist")
	serveCmd.PersistentFlags().StringP("email", "e", "", "The email to register with the ACME provider")
	serveCmd.PersistentFlags().StringP("storage", "s", DefaultIntercertDir+"/server-data", "The place to store certificates and other data")

	// Load serveCmd.PersistentFlags() values from config
	bindPrefixedFlags(serveCmd, "server", "port", "agree", "directory", "dns-provider", "domains", "storage", "email")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a server instance",
	Long:  `Start the server component, doing the interaction with the ACME server and connected clients`,
	Run: func(cmd *cobra.Command, args []string) {

		c := config.ServerConfig{
			Port:        viper.GetInt("server.port"),
			Agree:       viper.GetBool("server.agree"),
			Directory:   viper.GetString("server.directory"),
			DNSProvider: viper.GetString("server.dns-provider"),
			Domains:     viper.GetStringSlice("server.domains"),
			Email:       viper.GetString("server.email"),
			Storage:     viper.GetString("server.storage"),
		}

		if !c.Agree {
			PrintErrorAndExit(errors.New("the ACME ToS must be agreed to"))
		}

		log.Infof("Listening on port %v", c.Port)

		server.StartServer(&c)
	},
}

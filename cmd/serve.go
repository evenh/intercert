package cmd

import (
	"errors"
	"time"

	"github.com/go-acme/lego/log"

	"github.com/evenh/intercert/config"
	"github.com/evenh/intercert/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	serveCmd.PersistentFlags().DurationP("expiry", "x", 30*time.Minute, "How often to check for expired certificates")
	serveCmd.PersistentFlags().DurationP("renewalThreshold", "r", (24*time.Hour)*15, "How early before expiry shall certificates be renewed")

	// Load serveCmd.PersistentFlags() values from config
	bindPrefixedFlags(serveCmd, "server", "port", "agree", "directory", "dns-provider", "domains", "storage", "email", "expiry", "renewalThreshold")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a server instance",
	Long:  `Start the server component, doing the interaction with the ACME server and connected clients`,
	Run: func(cmd *cobra.Command, args []string) {

		c := config.ServerConfig{
			Port:             viper.GetInt("server.port"),
			Agree:            viper.GetBool("server.agree"),
			Directory:        viper.GetString("server.directory"),
			DNSProvider:      viper.GetString("server.dns-provider"),
			Domains:          viper.GetStringSlice("server.domains"),
			Email:            viper.GetString("server.email"),
			Storage:          viper.GetString("server.storage"),
			ExpiryCheckAt:    viper.GetDuration("server.expiry"),
			RenewalThreshold: viper.GetDuration("server.renewalThreshold"),
		}

		if !c.Agree {
			PrintErrorAndExit(errors.New("the ACME ToS must be agreed to"))
		}

		if c.RenewalThreshold > (24*time.Hour)*60 {
			PrintErrorAndExit(errors.New("renewal threshold can't exceed 60 days"))
		}

		server.StartServer(&c, UserAgent())
		log.Infof("Listening on port %v", c.Port)
	},
}

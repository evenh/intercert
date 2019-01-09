package cmd

import (
	"github.com/evenh/intercert/client"
	"github.com/evenh/intercert/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func init() {
	rootCmd.AddCommand(clientCmd)

	// Flags
	clientCmd.PersistentFlags().StringP("host", "H", "localhost", "The host (server) to connect to")
	clientCmd.PersistentFlags().IntP("port", "p", 6300, "The port the server will listen on")
	clientCmd.PersistentFlags().StringP("storage", "s", DefaultIntercertDir+"/client-data", "The place to store certificates and other data")
	clientCmd.PersistentFlags().StringArrayP("domains", "D", []string{}, "The domains to request certs for")
	clientCmd.PersistentFlags().DurationP("expiry", "e", 30*time.Minute, "How often to check for expired certificates")
	clientCmd.PersistentFlags().DurationP("renewalThreshold", "r", 24*time.Hour, "How early before expiry shall certificates be renewed")

	// Load clientCmd.PersistentFlags() values from config
	bindPrefixedFlags(clientCmd, "client", "host", "port", "storage", "domains", "expiry", "renewalThreshold")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a client instance",
	Long:  `Start a client instance, connecting to a running server instance`,
	Run: func(cmd *cobra.Command, args []string) {

		c := config.ClientConfig{
			Hostname:         viper.GetString("client.host"),
			Port:             viper.GetInt("client.port"),
			Storage:          viper.GetString("client.storage"),
			Domains:          viper.GetStringSlice("client.domains"),
			ExpiryCheckAt:    viper.GetDuration("client.expiry"),
			RenewalThreshold: viper.GetDuration("client.renewalThreshold"),
		}

		client.StartClient(&c, UserAgent())
	},
}

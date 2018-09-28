package cmd

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	// Flags
	serverCmd.Flags().IntP("port", "p", 6300, "The port the server will listen on")
	serverCmd.Flags().Bool("agree", false, "Whether you accept the ACME providers Terms of Service. Must be set to true in order to start server.")
	serverCmd.Flags().StringP("directory", "D", "https://acme-staging-v02.api.letsencrypt.org/directory", "URL to ACME directory")
	serverCmd.Flags().StringP("dns-provider", "P", "", "The DNS provider to use - see docs for valid providers.")
	serverCmd.Flags().StringSliceP("domains", "d", nil, "Domains to whitelist")
	serverCmd.Flags().StringP("email", "e", "", "The email to register with the ACME provider")
	serverCmd.Flags().StringP("storage", "s", DefaultIntercertDir + "/server-data", "The place to store certificates and other data")

	// Mark some as required
	serverCmd.MarkFlagRequired("directory")
	serverCmd.MarkFlagRequired("dns-provider")
	serverCmd.MarkFlagRequired("domains")
	serverCmd.MarkFlagRequired("email")
	serverCmd.MarkFlagRequired("storage")

	// Load flags values from config
	viper.BindPFlag("server.agree", serverCmd.Flags().Lookup("agree"))
	viper.BindPFlag("server.directory", serverCmd.Flags().Lookup("directory"))
	viper.BindPFlag("server.dns-provider", serverCmd.Flags().Lookup("dns-provider"))
	viper.BindPFlag("server.domains", serverCmd.Flags().Lookup("domains"))
	viper.BindPFlag("server.storage", serverCmd.Flags().Lookup("storage"))

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "The server doing the actual work",
	Long:  `Start the server component, doing the interaction with the ACME server and connected clients`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Yo server")
	},
}

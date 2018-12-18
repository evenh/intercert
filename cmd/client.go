package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xenolf/lego/log"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a client instance",
	Long:  `Start a client instance, connecting to a running server instance`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Client running!")
	},
}

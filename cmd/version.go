package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of intercert",
	Long:  `All software has versions. This is intercerts's`,
	Run: func(cmd *cobra.Command, args []string) {
		versionString := fmt.Sprintf("intercert %s (%s)", Version, Commit)
		fmt.Println(versionString)
	},
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xenolf/lego/log"
)

// PrintErrorAndExit logs a fatal error and exits the process with return code 1
func PrintErrorAndExit(err error) {
	log.Fatal(err)
	os.Exit(1)
}

func bindPrefixedFlag(cmd *cobra.Command, prefix string, key string) {
	err := viper.BindPFlag(prefix+"."+key, cmd.PersistentFlags().Lookup(key))

	if err != nil {
		PrintErrorAndExit(err)
	}
}

func bindPrefixedFlags(cmd *cobra.Command, prefix string, keys ...string) {
	for i := 0; i < len(keys); i++ {
		bindPrefixedFlag(cmd, prefix, keys[i])
	}
}

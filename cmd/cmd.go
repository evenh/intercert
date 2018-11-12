package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func printErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func bindPrefixedFlag(cmd *cobra.Command, prefix string, key string) {
	err := viper.BindPFlag(prefix+"."+key, cmd.PersistentFlags().Lookup(key))

	if err != nil {
		printErrorAndExit(err)
	}
}

func bindPrefixedFlags(cmd *cobra.Command, prefix string, keys ...string) {
	for i := 0; i < len(keys); i++ {
		bindPrefixedFlag(cmd, prefix, keys[i])
	}
}

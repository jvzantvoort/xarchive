package main

import (
	msg "github.com/jvzantvoort/xarchive/messages"
	"github.com/spf13/cobra"
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: msg.GetUsage("meta_root"),
	Long:  msg.GetLong("meta_root"),
}

func init() {
	rootCmd.AddCommand(metaCmd)
}

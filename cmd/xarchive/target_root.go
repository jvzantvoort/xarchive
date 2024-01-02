package main

import (
	msg "github.com/jvzantvoort/xarchive/messages"
	"github.com/spf13/cobra"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: msg.GetUsage("target_root"),
	Long:  msg.GetLong("target_root"),
}

func init() {
	rootCmd.AddCommand(targetCmd)
}

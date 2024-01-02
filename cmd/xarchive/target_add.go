package main

import (
	"github.com/jvzantvoort/xarchive/database"
	"github.com/jvzantvoort/xarchive/display"
	msg "github.com/jvzantvoort/xarchive/messages"
	"github.com/jvzantvoort/xarchive/target"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: msg.GetUsage("target_add"),
	Long:  msg.GetLong("target_add"),
	Run:   handleAddCmd,
}

func handleAddCmd(cmd *cobra.Command, args []string) {

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) == 0 {
		display.FatalIfFailed(ErrNoTargets, "Must provide a target")
	}
	db := database.NewDatabase(
		Username,
		Password,
		Hostname,
		Database,
		Port,
	)

	for _, arg := range args {
		tgt, err := target.NewTarget(arg)
		if err != nil {
			log.Fatal(err)
		}
		db.InsertTarget(tgt)
	}
}

func init() {
	targetCmd.AddCommand(addCmd)
}

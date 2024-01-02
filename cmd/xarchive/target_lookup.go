package main

import (
	"fmt"

	"github.com/jvzantvoort/xarchive/database"
	"github.com/jvzantvoort/xarchive/display"
	msg "github.com/jvzantvoort/xarchive/messages"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: msg.GetUsage("target_lookup"),
	Long:  msg.GetLong("target_lookup"),
	Run:   handleLookupCmd,
}

func handleLookupCmd(cmd *cobra.Command, args []string) {

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) == 0 {
		display.FatalIfFailed(fmt.Errorf("Missing Target"), "Must provide a target")
	}

	db := database.NewDatabase(
		Username,
		Password,
		Hostname,
		Database,
		Port,
	)

	for _, arg := range args {
		records, err := db.LookupTarget(arg)
		if err != nil {
			display.FatalIfFailed(err, "Must provide a target")
		}
		for _, rec := range records {
			path := rec.Path
			if targetExists(path) {
				fmt.Printf("%s\n", rec.GetDescr())
			}
		}
	}
}

func init() {
	targetCmd.AddCommand(lookupCmd)
}

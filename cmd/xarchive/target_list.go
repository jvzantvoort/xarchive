package main

import (
	"fmt"
	"strings"

	"github.com/jvzantvoort/xarchive/database"
	"github.com/jvzantvoort/xarchive/display"
	msg "github.com/jvzantvoort/xarchive/messages"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: msg.GetUsage("target_list"),
	Long:  msg.GetLong("target_list"),
	Run:   handleListCmd,
}

func handleListCmd(cmd *cobra.Command, args []string) {

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
		elements := strings.Split(arg, ":")
		if len(elements) != 2 {
			msg := fmt.Sprintf("target %s is mis formatted", arg)
			display.FatalIfFailed(fmt.Errorf("Mis formatted Target"), msg)
		}
		filename := elements[0]
		checksum := elements[1]
		records, err := db.GetTargets(filename, checksum)
		if err != nil {
			display.FatalIfFailed(err, "Must provide a target")

		}
		for _, rec := range records {
			path := rec.Path
			if targetExists(path) {
				fmt.Printf("%s\n", rec.Path)
			}
		}
	}
}

func init() {
	targetCmd.AddCommand(listCmd)
}

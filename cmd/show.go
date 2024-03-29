package cmd

import (
	pmm "github.com/hexylena/pm/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)

}

var showCmd = &cobra.Command{
	Use:   "show [note id]",
	Short: "show a note",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		partial := pmm.PartialNoteId(args[0])
		gn.BubbleShow(partial)
	},
}

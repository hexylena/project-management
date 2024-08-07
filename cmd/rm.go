package cmd

import (
	"fmt"
	pmm "github.com/hexylena/pm/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)

}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "rm",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		partial := pmm.PartialNoteId(args[0])
		note_id, err := gn.GetIdByPartial(partial)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("REMOVING THE FOLLOWING NOTE")
		gn.BubbleShow(partial)
		ga.DeleteNote(gn, note_id)
	},
}

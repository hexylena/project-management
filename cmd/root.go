package cmd

import (
	"fmt"
	pma "github.com/hexylena/pm/adapter"
	pmc "github.com/hexylena/pm/config"
	pmm "github.com/hexylena/pm/models"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "A project manager",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var gn pmm.GlobalNotes
var ga pma.TaskAdapter
var config pmc.HxpmConfig

func Execute(notes pmm.GlobalNotes, adapter pma.TaskAdapter, _config pmc.HxpmConfig) {
	gn = notes
	ga = adapter
	config = _config
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

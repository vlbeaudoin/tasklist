package cmd

import (
	"log"

	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"github.com/vlbeaudoin/tasklist/data"
)

func ListTasks() {
	tasks, err := data.ListTasks()
	if err != nil {
		log.Fatal(err)
	}

	t := tabby.New()
	t.AddHeader("ID", "LABEL")
	for _, task := range tasks {
		t.AddLine(task.ID, task.Label)
	}
	t.Print()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		data.OpenDatabase()
		data.MigrateDatabase()

		ListTasks()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

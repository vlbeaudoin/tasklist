package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vlbeaudoin/tasklist/data"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new task.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		// data.OpenDatabaseWithSqlite()
		data.OpenDatabase()
		data.MigrateDatabase()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

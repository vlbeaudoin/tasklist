package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vlbeaudoin/tasklist/data"
)

func createNewTaskFromArgs(args []string) {
	label := strings.Join(args, " ")
	fmt.Printf("Inserting task [%s]!\n", label)
	data.InsertTask(label)
}

func declareFlagsForAdd() {
	// general.list_after_add
	addCmd.Flags().BoolP(
		"list-after-add", "l", false,
		"List after adding task (config: 'general.list_after_add')")
	viper.BindPFlag(
		"general.list_after_add",
		addCmd.Flags().Lookup("list-after-add"))
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new task.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			data.OpenDatabase()
			data.MigrateDatabase()
			createNewTaskFromArgs(args)
		} else {
			fmt.Println("Not enough arguments after `add`.")
		}
		if viper.GetBool("general.list_after_add") {
			ListTasks()
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	declareFlagsForAdd()
}
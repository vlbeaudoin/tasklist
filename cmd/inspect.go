package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vlbeaudoin/tasklist/data"
)

func stringsToTaskIDs(strings []string) (taskIDs []uint64) {
	for _, s := range strings {
		if taskID, err := strconv.ParseUint(s, 10, 0); err != nil {
			log.Println(err)
		} else {
			// inspectTaskByID(taskID)
			taskIDs = append(taskIDs, taskID)
		}
	}
	return taskIDs
}

func inspectTaskByID(taskID uint64) error {
	// Task
	task, err := data.FindTaskByID(taskID)
	if err != nil {
		return err
	}

	log.Printf("[%d] %s", task.ID, task.Label)

	// Steps
	steps, err := data.FindStepsByTaskID(taskID)
	if err != nil {
		return err
	}

	for _, step := range steps {
		completed := ' '
		if step.Completed {
			completed = 'x'
		}
		log.Printf("- [%c] %s", completed, step.Description)
	}

	return nil
}

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspects a task and its steps.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Inspect requires at least one task ID.")
		}
		data.OpenDatabase()
		data.MigrateDatabase()

		for _, taskID := range stringsToTaskIDs(args) {
			err := inspectTaskByID(taskID)
			if err != nil {
				log.Println(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

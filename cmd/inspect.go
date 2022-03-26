package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vlbeaudoin/tasklist/data"
)

var all bool

func stringsToTaskIDs(strings []string) (taskIDs []uint64) {
	for _, s := range strings {
		if taskID, err := strconv.ParseUint(s, 10, 0); err != nil {
			log.Println("[E]", err)
		} else {
			taskIDs = append(taskIDs, taskID)
		}
	}
	return taskIDs
}

func showTask(task *data.Task) {
	fmt.Printf("( %d ) %s\n", task.ID, task.Label)
}

func showStep(step *data.Step) {
	completed := ' '
	if step.Completed {
		completed = 'x'
	}
	fmt.Printf("( %d.%d ) [%c] %s\n", step.TaskID, step.ID, completed, step.Description)
}

func showTasksAndSteps() {
	tasks, err := data.ListTasks()
	if err != nil {
		log.Println("[E]", err)
	}

	for _, task := range tasks {
		err = inspectTaskByID(task.ID)
		if err != nil {
			log.Fatal("[E]", err)
		}
	}
}

func inspectTaskByID(taskID uint64) error {
	// Task
	task, err := data.FindTaskByID(taskID)
	if err != nil {
		return err
	}

	showTask(&task)

	// Steps
	steps, err := data.FindStepsByTaskID(taskID)
	if err != nil {
		return err
	}

	for _, step := range steps {
		showStep(&step)
	}

	return nil
}

func declareFlagsForInspect() {
	// all
	inspectCmd.Flags().BoolVarP(
		&all, "all", "a", false,
		"Inspect all tasks and steps regardless of args. (Default: false)")
}

func connectDB() {
	err := data.OpenDatabase()
	if err != nil {
		log.Fatal("[E]", err)
	}

	err = data.MigrateDatabase()
	if err != nil {
		log.Fatal("[E]", err)
	}
}

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspects a task and its steps.",
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			connectDB()
			showTasksAndSteps()
		} else {
			if len(args) == 0 {
				log.Fatal("[E] Inspect requires at least one task ID or the '--all, -a' flag.")
			}
			connectDB()

			// Show tasks from taskIDs in args
			for _, taskID := range stringsToTaskIDs(args) {
				err := inspectTaskByID(taskID)
				if err != nil {
					log.Println("[E]", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	declareFlagsForInspect()
}

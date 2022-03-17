package cmd

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vlbeaudoin/tasklist/data"
)

func declareFlagsForImport() {
	// general.list_after_import
	importCmd.Flags().BoolP(
		"list-after-import", "l", false,
		"List after importing tasks (config: 'general.list_after_import')")
	viper.BindPFlag(
		"general.list_after_import",
		importCmd.Flags().Lookup("list-after-import"))
}

func isExistingFile(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func checkFileFormat(path string) string {
	s := strings.Split(path, ".")
	extension := s[len(s)-1]

	return extension
}

/*
validatePaths sorts provided paths according to their state of existence.

It takes a slice of strings and returns 2 slices of strings,
one for valid paths and one for invalid paths.
*/
func validatePaths(paths []string) (validPaths []string, invalidPaths []string) {
	for _, path := range paths {
		if isExistingFile(path) {
			validPaths = append(validPaths, path)
		} else {
			invalidPaths = append(invalidPaths, path)
		}
	}
	return validPaths, invalidPaths
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports a list of tasks from a file.",
	Run: func(cmd *cobra.Command, args []string) {
		successfulImport := false

		if l := len(args); l == 0 {
			log.Fatal("Not enough arguments after `import`.")
		} else {
			data.OpenDatabase()
			data.MigrateDatabase()

			validPaths, invalidPaths := validatePaths(args)

			if len(invalidPaths) > 0 {
				log.Println("Invalid import path(s):", invalidPaths)
			}

			if len(validPaths) > 0 {
				log.Println("Valid import path(s):", validPaths)
			} else {
				log.Fatal("No valid paths found.")
			}

			for _, validPath := range validPaths {
				switch f := checkFileFormat(validPath); f {
				case "csv":
					log.Printf("Detected type '%s' in file '%s'.\n", "csv", validPath)

					file, err := os.Open(validPath)
					if err != nil {
						log.Println(err)
						break
					}
					defer file.Close()

					tasks := []*data.Task{}

					if err := gocsv.UnmarshalFile(file, &tasks); err != nil {
						log.Println("Error during csv unmarshal:", err)
						break
					}

					log.Printf("Found %d records in %s.\n", len(tasks), validPath)

					ignoredPaths := []*data.Task{}

					for i, task := range tasks {
						if task.Label == "" {
							ignoredPaths = append(ignoredPaths, task)
							tasks = append(tasks[:i], tasks[i+1:]...)
						}
					}

					if len(ignoredPaths) != 0 {
						log.Printf("%d record(s) are missing a `Label` and will be ignored.\n", len(ignoredPaths))
					}

					log.Printf("Inserting %d records.\n", len(tasks))

					err = data.InsertTasks(tasks)
					if err != nil {
						log.Println(err)
					}

					successfulImport = true

				case "json":
					log.Printf("Detected type '%s' in file '%s'.\n", "json", validPath)

					// TODO Verify is file is valid json ()
					// go doc json.valid

					// TODO import file as json
					log.Println("Not yet implemented: json import")
					break
				default:
					log.Printf("Unknown type '%s' for file '%s', skipping.\n", f, validPath)
				}
			}

			if successfulImport && viper.GetBool("general.list_after_import") {
				ListTasks()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	declareFlagsForImport()
}

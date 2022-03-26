package data

import (
	"errors"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	gorm.Model
	ID    uint64 `mapper:"id" json:"id"`
	Label string `csv:"label" json:"label"`
	Steps []Step
}

type Step struct {
	gorm.Model
	ID          uint64 `mapper:"id" json:"id"`
	Description string `csv:"description" json:"description"`
	Completed   bool   `csv:"completed" json:"completed"`
	TaskID      uint   `csv:"taskid" json:"taskid"`
}

func OpenDatabase() error {
	var err error

	var dialector gorm.Dialector

	switch t := viper.GetString("db.type"); t {
	case "sqlite":
		log.Println("Using driver gorm.io/driver/sqlite")

		db_path := viper.GetString("db.path")

		if db_path == "" {
			log.Fatal("No valid database file found in `--db-path` or `db.path`.")
		}

		log.Println("Using database file:", db_path)

		dialector = sqlite.Open(db_path)
	default:
		log.Fatalf("Unrecognized database driver requested (%s).\n", t)
	}

	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func MigrateDatabase() error {
	err := db.AutoMigrate(&Task{}, &Step{})
	return err
}

func InsertTask(label string) {
	db.Create(&Task{
		Label: label,
	})
}

func InsertTaskWithSteps(label string, steps []string) {
	if len(steps) > 0 {
		// Populate steps
		structSteps := []Step{}

		for _, step := range steps {
			structSteps = append(structSteps, Step{Description: step})
		}

		// Insert task and steps
		db.Create(&Task{
			Label: label,
			Steps: structSteps,
		})
	} else {
		InsertTask(label)
	}
}

func ListTasks() ([]Task, error) {
	var tasks []Task

	result := db.Model(&Task{}).Find(&tasks)

	return tasks, result.Error
}

func InsertTasks(tasks []*Task) error {
	if len(tasks) == 0 {
		return errors.New("Cannot insert empty batch of tasks.")
	}

	for _, task := range tasks {
		task.ID = 0
	}

	db.CreateInBatches(&tasks, 500)

	return nil
}

func FindTaskByID(taskID uint64) (task Task, err error) {
	result := db.First(&task, taskID)
	return task, result.Error
}

func FindStepsByTaskID(taskID uint64) (steps []Step, err error) {
	result := db.Where("task_id = ?", taskID).Find(&steps)
	return steps, result.Error
}

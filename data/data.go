package data

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	gorm.Model
	Label string `csv:"label" json:"label"`
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

func MigrateDatabase() {
	db.AutoMigrate(&Task{})
}

func InsertTask(label string) {
	db.Create(&Task{
		Label: label,
	})
}

func ListTasks() ([]Task, error) {
	var tasks []Task

	result := db.Model(&Task{}).Find(&tasks)

	return tasks, result.Error
}

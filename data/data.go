package data

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	gorm.Model
	Label string `json:"label"`
}

func OpenDatabaseWithSqlite() error {
	var err error

	db, err = gorm.Open(sqlite.Open(viper.GetString("db.path")), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func OpenDatabase() error {
	var err error

	var dialector gorm.Dialector

	switch t := viper.GetString("db.type"); t {
	case "sqlite":
		fmt.Println("Using driver gorm.io/driver/sqlite")
		dialector = sqlite.Open(viper.GetString("db.path"))
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

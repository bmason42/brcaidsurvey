package db

import (
	"brcaidsurvey/pkg/model"
	"fmt"
	_ "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var dbPath string

func InitDB() error {
	home := os.Getenv("HOME")
	dbPath = fmt.Sprintf("%s/brcaid-db/db.db", home)
	createDbORM()
	return nil
}

func GetDBConnection() (*gorm.DB, error) {
	db, error := gorm.Open("sqlite3", dbPath)
	return db, error
}
func createDbORM() {
	db, _ := GetDBConnection()
	db.AutoMigrate(&model.RegionInfo{},
		&model.SupportConcern{},
		&model.SurveyContact{},
		&model.SurveyResult{})

	//db.Raw("DROP INDEX event_idx on log_entry")
	//db.Raw("DROP INDEX tracking_idx on log_entry")
}
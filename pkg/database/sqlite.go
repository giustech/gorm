package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"time"
)

var SqliteGorm GormConnection

type (
	sqliteConnection struct {
	}
)

func init() {
	SqliteGorm = sqliteConnection{}
}

func (sqliteConnection) GetInstance(_ string, _ string, _ string, _ int, _ string, _ bool) (*gorm.DB, error) {
	var err error
	var connection gorm.Dialector
	var config gorm.Config
	now := time.Now()
	timestampStr := now.Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("sqlite/%s.db", strings.Replace(strings.ReplaceAll(strings.ReplaceAll(timestampStr, "-", ""), ":", ""), " ", "", 1))
	connection = sqlite.Open(sql)
	gormDb, err := gorm.Open(connection, &config)

	if err != nil {
		if strings.ToUpper(os.Getenv("ENVIRONMENT")) != "PRODUCTION" {
			log.Printf("Error to inject sqlite")
		}
	}

	return gormDb, err
}

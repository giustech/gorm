package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	var config gorm.Config
	now := time.Now()
	timestampStr := now.Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("sqlite/%s.db", strings.Replace(strings.ReplaceAll(strings.ReplaceAll(timestampStr, "-", ""), ":", ""), " ", "", 1))
	if !checkIfFolderExists("./sqlite") {
		fmt.Println("Folder sqlite must be created in test context repository")
	}
	return gorm.Open(sqlite.Open(sql), &config)
}

func checkIfFolderExists(folderPath string) bool {
	_, err := os.ReadDir(folderPath)
	if err != nil {
		return false
	}
	return true
}

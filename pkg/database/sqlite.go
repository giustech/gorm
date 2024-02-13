package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	return gorm.Open(sqlite.Open(sql), &config)
}

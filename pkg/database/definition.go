package database

import "gorm.io/gorm"

type (
	GormConnection interface {
		GetInstance(user string, pass string, host string, port int, dbName string, sslMode bool) (*gorm.DB, error)
	}
)

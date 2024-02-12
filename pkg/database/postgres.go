package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var PostgresGorm GormConnection

type (
	postgresConnection struct {
	}
)

func init() {
	PostgresGorm = postgresConnection{}
}

func (postgresConnection) GetInstance(user string, pass string, host string, port int, dbName string, sslMode bool) (*gorm.DB, error) {
	var err error
	var connection gorm.Dialector
	var config gorm.Config
	connection = postgres.Open(generateDsn(user, pass, host, port, dbName, sslMode))
	config = gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "entities.",
			SingularTable: false,
		},
	}
	gormDb, err := gorm.Open(connection, &config)

	if err != nil {
		panic(err)
	}
	return gormDb, err
}

package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
)

type (
	RepositoryDao struct {
		poolGormDb map[string]*gorm.DB
		once       sync.Once
	}

	RepositoryConfiguration struct {
		ConnectionName string
		DBUser         string
		DBPass         string
		DBHost         string
		DBPort         int
		DBName         string
	}
)

var Dao = &RepositoryDao{}

func generateDsn(user string, pass string, host string, port int, dbName string, sslMode bool) string {
	if sslMode {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=enable", host, port, user, dbName, pass)
	} else {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, pass)
	}
}

func (r *RepositoryDao) Start(configurations []RepositoryConfiguration) {
	r.once.Do(func() {
		r.poolGormDb = make(map[string]*gorm.DB)
		for _, configuration := range configurations {
			var con *gorm.DB
			var err error
			switch strings.ToUpper(os.Getenv("ENVIRONMENT")) {
			case "TEST":
				con, err = SqliteGorm.GetInstance(configuration.DBUser, configuration.DBPass, configuration.DBHost, configuration.DBPort, configuration.DBName, false)
				break
			default:
				con, err = PostgresGorm.GetInstance(configuration.DBUser, configuration.DBPass, configuration.DBHost, configuration.DBPort, configuration.DBName, false)
				break
			}
			if err != nil {
				panic(err)
			}
			r.poolGormDb[configuration.ConnectionName] = con
		}
	})
}

func (r *RepositoryDao) GetInstance(name string) *gorm.DB {
	return r.poolGormDb[name]
}

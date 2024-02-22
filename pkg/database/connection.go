package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"sync"
)

type (
	daoConfiguration struct {
		gormDb *gorm.DB
		config RepositoryConfiguration
	}

	RepositoryDao struct {
		poolGormDb map[string]daoConfiguration
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

	GormData struct {
		Config string
	}
)

var (
	Dao            = &RepositoryDao{}
	configurations []RepositoryConfiguration
)

func (gormData GormData) GetDao() *gorm.DB {
	return Dao.getInstance(gormData.Config)
}

func generateDsn(user string, pass string, host string, port int, dbName string, sslMode bool) string {
	if sslMode {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=enable", host, port, user, dbName, pass)
	} else {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, pass)
	}
}

func (r *RepositoryDao) Start(_configurations []RepositoryConfiguration) {
	r.once.Do(func() {
		configurations = _configurations
		r.initConfiguration()
	})
}

func (r *RepositoryDao) initConfiguration() {
	environment := strings.ToUpper(os.Getenv("ENVIRONMENT"))
	r.poolGormDb = make(map[string]daoConfiguration)
	if environment != "TEST" {
		for _, configuration := range configurations {
			log.Printf("Init instance: %s", configuration.DBName)
			con, err := PostgresGorm.GetInstance(configuration.DBUser, configuration.DBPass, configuration.DBHost, configuration.DBPort, configuration.DBName, false)
			if err != nil {
				log.Println(err.Error())
				log.Fatal(err)
			}
			r.poolGormDb[configuration.ConnectionName] = daoConfiguration{
				gormDb: con,
				config: configuration,
			}
		}
	} else {
		con, err := SqliteGorm.GetInstance("", "", "", 5432, "", false)
		if err != nil {
			log.Printf("Error to inject sqlite")
		} else {
			r.poolGormDb["sqlite"] = daoConfiguration{
				gormDb: con,
				config: RepositoryConfiguration{
					ConnectionName: "",
					DBUser:         "",
					DBPass:         "",
					DBHost:         "",
					DBPort:         0,
					DBName:         "",
				},
			}
		}
	}
}

func (r *RepositoryDao) getInstance(name string) *gorm.DB {
	return r.poolGormDb[name].gormDb
}

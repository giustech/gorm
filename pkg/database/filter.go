package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

func buildQuery(db *gorm.DB, filterValue reflect.Value, filterType reflect.Type) *gorm.DB {
	for i := 0; i < filterValue.NumField(); i++ {
		fieldValue := filterValue.Field(i)
		fieldType := filterType.Field(i)

		kind := fieldType.Type.Kind()
		isUUID := fieldType.Type.PkgPath() == "github.com/google/uuid" && fieldType.Type.Name() == "UUID"

		var valueStr string
		if isUUID {
			if u, ok := fieldValue.Interface().(uuid.UUID); ok && u != uuid.Nil {
				valueStr = u.String()
			}
		} else {
			switch kind {
			case reflect.String:
				valueStr = fieldValue.String()
			case reflect.Bool:
				if fieldValue.Bool() {
					valueStr = "true"
				} else {
					valueStr = "false"
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				valueStr = strconv.FormatInt(fieldValue.Int(), 10)
			case reflect.Float32, reflect.Float64:
				valueStr = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
			}
		}

		if valueStr != "" {
			dbTag := fieldType.Tag.Get("gorm")
			if dbTag != "" {
				tags := strings.Split(dbTag, ";")
				dbColumnName := ""
				for _, tag := range tags {
					if strings.HasPrefix(tag, "column:") {
						dbColumnName = strings.TrimPrefix(tag, "column:")
						break
					}
				}
				if dbColumnName != "" {
					if kind == reflect.Bool {
						db = db.Where(dbColumnName+" = ?", fieldValue.Bool())
					} else {
						db = db.Where(dbColumnName+" = ?", valueStr)
					}
				}
			}
		}
	}
	return db
}

func CountFiltered(db *gorm.DB, filterValue reflect.Value, filterType reflect.Type, count *int64) *gorm.DB {
	_gorm := buildQuery(db, filterValue, filterType).Count(count)
	return _gorm
}

func PaginatedFiltered(db *gorm.DB, filterValue reflect.Value, filterType reflect.Type, pageSize int, pageNum int, result interface{}) *gorm.DB {
	var _gorm *gorm.DB
	if pageSize > 0 && pageNum > 0 {
		_gorm = buildQuery(db, filterValue, filterType).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(result)
	} else {
		_gorm = buildQuery(db, filterValue, filterType).Find(result)
	}
	return _gorm
}

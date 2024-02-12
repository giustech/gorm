package database

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strings"
)

const (
	ErrorToSaveNilObject   = "is not possible to save nil object"
	ErrorToUpdateNilObject = "is not possible to update nil object"
)

func HandlerPosgresErrors(context *gin.Context, entityName string, err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if strings.HasPrefix(pgErr.Message, "duplicate key value violates unique constraint") {
			context.String(http.StatusConflict, fmt.Sprintf("Entity %s already exists", entityName))
			return nil
		} else {
			context.String(http.StatusInternalServerError, err.Error())
			return nil
		}
	} else {
		return err
	}
}

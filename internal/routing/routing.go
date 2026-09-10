package routing

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	return router
}

package routing

import (
	"database/sql"

	db_pkg "code.siru.ink/vdh/wahlen-web/internal/db"
	"github.com/gin-gonic/gin"
)

func InitializeRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.Use(db_pkg.DbMiddleware(db))

	return router
}

package routing

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"

	db_pkg "code.siru.ink/vdh/wahlen-web/internal/db"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*.tmpl
var embeddedHtml embed.FS

func InitializeRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	sub, err := fs.Sub(embeddedHtml, "templates")
	if err != nil {
		log.Fatalf("Failed to load sub filesystem: %s", err)
	}
	router.LoadHTMLFS(http.FS(sub), "*.tmpl")

	router.Use(db_pkg.DbMiddleware(db))

	router.GET("/login", Login)
	router.GET("/signup", Signup)
	// router.POST("/signup", SignupPost)

	return router
}

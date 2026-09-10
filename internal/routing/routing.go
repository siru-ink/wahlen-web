package routing

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	db_pkg "code.siru.ink/vdh/wahlen-web/internal/db"
	"github.com/gin-gonic/gin"
)

//go:embed templates
var embeddedHtml embed.FS

func InitializeRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	sub, _ := fs.Sub(embeddedHtml, "templates")
	router.LoadHTMLFS(http.FS(sub), "*.tmpl")

	router.Use(db_pkg.DbMiddleware(db))

	router.GET("/votecount", VoteCount)

	return router
}

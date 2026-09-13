package routing

import (
	"net/http"

	"code.siru.ink/vdh/wahlen-web/internal/cookie"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	flash_message, _ := cookie.Break(c, cookie.FlashCookie)

	context := gin.H{
		"flash_message": flash_message,
	}

	c.HTML(http.StatusOK, "login.tmpl", context)
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

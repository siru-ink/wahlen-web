package routing

import (
	"net/http"

	"code.siru.ink/vdh/wahlen-web/internal/cookie"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	flash_message, _ := cookie.FlashCookie.Break(c)

	context := gin.H{
		"flash_message": flash_message,
	}

	c.HTML(http.StatusOK, "login.tmpl", context)
}

func Signup(c *gin.Context) {
	flash_message, _ := cookie.FlashCookie.Break(c)

	context := gin.H{
		"flash_message": flash_message,
	}

	c.HTML(http.StatusOK, "signup.tmpl", context)
}

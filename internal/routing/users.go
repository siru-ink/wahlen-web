package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

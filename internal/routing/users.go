package routing

import (
	"database/sql"
	"fmt"
	"net/http"

	"code.siru.ink/vdh/wahlen-web/internal/cookie"
	"code.siru.ink/vdh/wahlen-web/internal/db"
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

func SignupPost(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		context := gin.H{
			"flash_message": "Username must not be empty.",
		}
		c.HTML(http.StatusUnprocessableEntity, "signup.tmpl", context)
		return
	}

	password := c.PostForm("password")
	if checkPasswordGuidelines(password) {
		context := gin.H{
			"flash_message": "Password must be at least 15 characters long and contain at least 1 number.",
		}
		c.HTML(http.StatusUnprocessableEntity, "signup.tmpl", context)
		return
	}

	dbase := c.MustGet("db").(*sql.DB)
	if db.CheckUserExistence(dbase, c.Request.Context(), username) {
		context := gin.H{
			"flash_message": fmt.Sprintf(
				"The username <%s> is already taken. Please choose a different one.",
				username),
		}

		c.HTML(http.StatusUnprocessableEntity, "signup.tmpl", context)
		return
	}

	db.RegisterUser(dbase, c.Request.Context(), username, password)

	cookie.FlashCookie.Seal(
		c,
		fmt.Sprintf(
			"You username <%s> is now registered. You can log in with the new account now.",
			username))

	c.Redirect(http.StatusSeeOther, "/login")
}

func checkPasswordGuidelines(password string) (failed bool) {
	return true // all password attempts fail until implemented
}

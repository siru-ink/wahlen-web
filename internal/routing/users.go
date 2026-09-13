package routing

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"code.siru.ink/vdh/wahlen-web/internal/cookie"
	"code.siru.ink/vdh/wahlen-web/internal/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/unicode/norm"
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
	usernameRaw := c.PostForm("username")
	if usernameRaw == "" {
		context := gin.H{
			"flash_message": "Username must not be empty.",
		}
		c.HTML(http.StatusUnprocessableEntity, "signup.tmpl", context)
		return
	}
	username := norm.NFC.String(usernameRaw)

	passwordRaw := c.PostForm("password")
	password := norm.NFC.String(passwordRaw)
	if passwordFailsGuidelines(password) {
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

func passwordFailsGuidelines(password string) (failed bool) {
	// normalize unicode representation before doing length counting
	if !norm.NFC.IsNormalString(password) {
		password = norm.NFC.String(password)
	}

	containsDigit := strings.ContainsFunc(password, unicode.IsDigit)
	atLeast15Char := utf8.RuneCountInString(password) >= 15
	return containsDigit && atLeast15Char
}

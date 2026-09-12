package routing

import (
	"log"
	"net/http"

	"code.siru.ink/vdh/wahlen-web/internal/crypto"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
	text, _ := crypto.Encrypt("Hello, world!")
	log.Printf("%s", text)
	textNew, err := crypto.Decrypt(text)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("%s", textNew)
}

// func SignupPost(c *gin.Context) {
// 	username := c.PostForm("username")
// 	password := c.PostForm("password")

// 	c.C

// 	hashedPassword, err := passwd.HashPassword(password, passwd.DefaultArgonParams())

// }

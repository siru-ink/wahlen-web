package routing

import (
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"code.siru.ink/vdh/wahlen-web/internal/cookie"
	"code.siru.ink/vdh/wahlen-web/internal/crypto"
	"code.siru.ink/vdh/wahlen-web/internal/db"
	"code.siru.ink/vdh/wahlen-web/internal/passwd"
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
	rawUsername := c.PostForm("username")
	rawPassword := c.PostForm("password")
	dbase := c.MustGet("db").(*sql.DB)

	username, err := createUserHelper(rawUsername, rawPassword, dbase, c)
	if err != nil {
		context := gin.H{
			"flash_message": err.Error(),
		}
		c.HTML(http.StatusUnprocessableEntity, "signup.tmpl", context)
		return
	}

	cookie.FlashCookie.Seal(
		c,
		fmt.Sprintf(
			"You username <%s> is now registered. You can log in with the new account now.",
			username))

	c.Redirect(http.StatusSeeOther, "/login")
}

func createUserHelper(rawUsername, rawPassword string, dbase *sql.DB, c *gin.Context) (username string, err error) {
	// What do I get from the user?
	// - username
	// - password
	//
	// What do I need to check?
	// - username must not already exist (uniqueness guarantee)
	// - password must be sufficiently long/secure
	//
	// What do I need to be able to do with this information?
	// - Identify user
	// - Authenticate user
	// - Have a server-unknown encryption key
	//
	// How do I get those things from what I have?
	// - username -> uniquely identify user
	// - full password -> argon2 hash, stored on server, can be verified on future login
	// - partial password -> argon2 hash, server stores salt, cookie stores hash value, can
	//                       be used for encrypting server storage

	// Checks
	if !norm.NFC.IsNormalString(rawUsername) {
		username = norm.NFC.String(rawUsername)
	} else {
		username = rawUsername
	}
	if username == "" {
		return "", errors.New("username must not be empty")
	}
	if db.CheckUserExistence(dbase, c.Request.Context(), username) {
		return "", errors.New("username is already taken")
	}

	password := norm.NFC.String(rawPassword)
	if !strings.ContainsFunc(password, unicode.IsDigit) {
		return "", errors.New("password must contain a number")
	}
	if !(utf8.RuneCountInString(password) >= 20) {
		return "", errors.New("password must at least be 20 characters long")
	}

	// Setup encryptions
	encryptionKeyHash, err := passwd.HashPassword(password[:10], passwd.DefaultArgonParams())
	if err != nil {
		return "", errors.New("hashing password failed")
	}

	encodedEncryptionKeyParams, encodedEncryptionKey, err := passwd.SplitHash(encryptionKeyHash)
	if err != nil {
		log.Fatalf("Deconstructing a new hash must not fail: %v", err)
		// return "", errors.New("deconstructing new hash failed")
	}

	encryptionKey, err := passwd.GetKey(encodedEncryptionKey)
	if err != nil {
		log.Fatal("Deocding the base64 encryption key must not fail: %v", err)
		// return "", errors.New("decoding encryption key failed")
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	encryptedPrivateKey, err := crypto.EncryptVia(fmt.Sprintf("%x", privateKey), encryptionKey)
	if err != nil {
		return "", errors.New("encrypting newly created private key failed")
	}

	encodedPublicKey := base64.RawStdEncoding.EncodeToString(publicKey)

	// Creation in dbase
	err = db.RegisterUser(dbase, c.Request.Context(), username, password, encodedPublicKey, encryptedPrivateKey, encodedEncryptionKeyParams)
	if err != nil {
		return "", errors.New("database failed to create user; please try again")
	}

	// Save encryption key cookie?
	// - Yes: Get it out of the way now
	// -  No: not setting any of the other cookies, and redirecting to login.
	//        Login will have to reset this cookie anyway just in case
	// encryptionPasswordChunk := password[:10]
	// err = cookie.CryptoCookie.Seal(c, encryptionPasswordChunk)
	// if err != nil {
	// 	return username, errors.New("user created, but encryption cookie not set")
	// }

	return username, nil
}

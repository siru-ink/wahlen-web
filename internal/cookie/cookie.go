package cookie

import (
	"net/http"
	"time"

	"code.siru.ink/vdh/wahlen-web/internal/crypto"
	"github.com/gin-gonic/gin"
)

type CookieType int

const (
	SessionCookie CookieType = iota
	FlashCookie
	CryptoCookie
)

var cookieName = map[CookieType]string{
	SessionCookie: "session-id",
	FlashCookie:   "flash-message",
	CryptoCookie:  "crypto-key",
}

var (
	cookiePath         = "/"
	cookieDomain       = "localhost"
	cookieMaxAge       = 24 * 60 * 60 // one full day
	cookieSecureOnly   = false
	cookieHttpOnly     = true
	cookieSameSiteMode = http.SameSiteLaxMode
)

func (ct CookieType) String() string {
	return cookieName[ct]
}

func (ct CookieType) Seal(c *gin.Context, value string) (err error) {
	encryptedValue, err := crypto.Encrypt(value)
	if err != nil {
		return err
	}

	cookieData := http.Cookie{
		Name:     ct.String(),
		Value:    encryptedValue,
		Path:     cookiePath,
		Domain:   cookieDomain,
		MaxAge:   cookieMaxAge,
		Secure:   cookieSecureOnly,
		HttpOnly: cookieHttpOnly,
		SameSite: cookieSameSiteMode,
	}

	c.SetCookieData(&cookieData)

	return nil
}

func (ct CookieType) Break(c *gin.Context) (value string, err error) {
	cookie, err := c.Request.Cookie(ct.String())
	if err != nil {
		return "", nil
	}

	// Flash messages should only be displayed once, so reading them
	// automatically destroys the cookie
	if ct == FlashCookie {
		FlashCookie.Loose(c)
	}

	decryptedValue, err := crypto.Decrypt(cookie.Value)
	if err != nil {
		return "", err
	}
	return decryptedValue, nil
}

func (ct CookieType) Loose(c *gin.Context) {
	cookieData := http.Cookie{
		Name:     ct.String(),
		Value:    "",
		Path:     cookiePath,
		Domain:   cookieDomain,
		Expires:  time.Now(),
		MaxAge:   -1, // delete immediately
		Secure:   cookieSecureOnly,
		HttpOnly: cookieHttpOnly,
		SameSite: cookieSameSiteMode,
	}

	c.SetCookieData(&cookieData)
}

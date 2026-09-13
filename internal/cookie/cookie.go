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
)

var cookieName = map[CookieType]string{
	SessionCookie: "session-id",
}

func (ct CookieType) String() string {
	return cookieName[ct]
}

func Build(name CookieType, value string) (newCookie http.Cookie, err error) {
	encryptedValue, err := crypto.Encrypt(value)
	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:     name.String(),
		Value:    encryptedValue,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

func Break(c *gin.Context, name CookieType) (value string, err error) {
	cookie, err := c.Request.Cookie(name.String())
	if err != nil {
		return "", nil
	}

	decryptedValue, err := crypto.Decrypt(cookie.Value)
	if err != nil {
		return "", err
	}
	return decryptedValue, nil
}

func Loose(name CookieType) (newCookie http.Cookie) {
	return http.Cookie{
		Name:     name.String(),
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now(),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
	"golang.org/x/oauth2"
)

func SetCookie(w http.ResponseWriter, name string, value string, cc config.Cookie) {
	fmt.Println("cookie:", name, value, time.Hour*time.Duration(cc.LifetimeHours), cc.Path)
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cc.Path,
		Expires:  time.Now().Add(time.Hour * time.Duration(cc.LifetimeHours)),
		HttpOnly: cc.HTTPOnly,
	}
	http.SetCookie(w, cookie)
	return
}

func DeleteCookie(w http.ResponseWriter, name string, cc config.Cookie) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    cc.Path,
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
	return
}

func GetCookie(r *http.Request, key string) (string, error) {
	utils.Debug(false, "look for key ", key)
	cookie, err := r.Cookie(key)
	if err != nil || cookie == nil {
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func GetToken(r *http.Request, cc config.Cookie, isReserve bool) (oauth2.Token, error) {

	var (
		token        oauth2.Token
		expireString string
		err          error
		aKey         = cc.Auth.AccessToken
		tKey         = cc.Auth.TokenType
		rKey         = cc.Auth.RefreshToken
		eKey         = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}
	token.AccessToken, err = GetCookie(r, aKey)
	if err != nil {
		return token, err
	}
	token.TokenType, err = GetCookie(r, tKey)
	if err != nil {
		return token, err
	}
	token.RefreshToken, err = GetCookie(r, rKey)
	if err != nil {
		return token, err
	}
	expireString, err = GetCookie(r, eKey)
	if err != nil {
		return token, err
	}
	token.Expiry, err = time.Parse("2006-01-02 15:04:05", expireString)

	return token, err
}

func SetToken(rw http.ResponseWriter, isReserve bool, token oauth2.Token, cc config.Cookie) {

	var (
		aKey = cc.Auth.AccessToken
		tKey = cc.Auth.TokenType
		rKey = cc.Auth.RefreshToken
		eKey = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}

	SetCookie(rw, aKey, token.AccessToken, cc)
	SetCookie(rw, tKey, token.TokenType, cc)
	SetCookie(rw, rKey, token.RefreshToken, cc)
	SetCookie(rw, eKey, token.Expiry.Format("2006-01-02 15:04:05"), cc)
}

func DeleteToken(rw http.ResponseWriter, isReserve bool, cc config.Cookie) {

	var (
		aKey = cc.Auth.TokenType
		tKey = cc.Auth.TokenType
		rKey = cc.Auth.RefreshToken
		eKey = cc.Auth.Expire
	)
	if isReserve {
		aKey = cc.Auth.ReservePrefix + aKey
		tKey = cc.Auth.ReservePrefix + tKey
		rKey = cc.Auth.ReservePrefix + rKey
		eKey = cc.Auth.ReservePrefix + eKey
	}

	DeleteCookie(rw, aKey, cc)
	DeleteCookie(rw, tKey, cc)
	DeleteCookie(rw, rKey, cc)
	DeleteCookie(rw, eKey, cc)
}

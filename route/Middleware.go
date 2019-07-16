package route

import (
	"errors"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhenorzz/goploy/core"
)

// CheckToken return token is vaild. Besides user/login router
func CheckToken(w http.ResponseWriter, r *http.Request) error {
	if "/user/login" == r.URL.Path {
		return nil
	}
	if "/user/isShowPhrase" == r.URL.Path {
		return nil
	}
	goployTokenCookie, err := r.Cookie("goploy_token")
	if err != nil {
		response := core.Response{Code: 1, Message: "非法请求"}
		response.Json(w)
		return errors.New("非法请求")
	}
	unPraseToken := goployTokenCookie.Value
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(unPraseToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGN_KEY")), nil
	})

	if err == nil && token.Valid {
		core.GolbalUserID = uint32(claims["id"].(float64))
		core.GolbalUserName = claims["name"].(string)
	} else {
		response := core.Response{Code: 401, Message: err.Error()}
		response.Json(w)
		return err
	}
	return nil
}

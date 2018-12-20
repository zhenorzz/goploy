package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhenorzz/goploy/core"
)

func CheckToken(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		token string
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	token, err := jwt.Parse(reqData.token, func(token *jwt.Token) (interface{}, error) {
		return []byte("your key"), nil
	})

	if err == nil && token.Valid {
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
}

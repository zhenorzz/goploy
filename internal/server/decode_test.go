package server

import (
	"fmt"
	"github.com/zhenorzz/goploy/internal/model"
	"net/url"
	"strconv"
	"testing"
)

func Decode(t *testing.T, g *Goploy) {
	// When you want to parse data from the http URL query, please use the `schema` tag to receive the data.
	// If you need to test, please uncomment it and see the TestDecodeUrlQuery function
	type ReqData struct {
		ServerID int64  `schema:"serverId" validate:"gt=0"`
		Dir      string `schema:"nginxDir" validate:"filepath"` // custom verify filepath rule, you can find in internal/validator
		Password string `schema:"passwd" validate:"password"`   // custom verify password rule, you can find in internal/validator
		File     string `schema:"fileName" validate:"required"`
	}

	// When you want to parse data from the http body, please use `json` tags to receive the data.
	// If you need to test, please uncomment it and see the TestDecodeBody function
	//type ReqData struct {
	//	ServerID int64  `json:"serverId" validate:"gt=0"`
	//	Dir      string `json:"nginxDir" validate:"filepath"` // custom verify filepath rule, you can find in internal/validator
	//	Password string `json:"passwd" validate:"password"`   // custom verify password rule, you can find in internal/validator
	//	File     string `json:"fileName" validate:"required"`
	//}

	var reqData ReqData
	if err := g.Decode(&reqData); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", reqData)
}

// Test parsing data from URL query
// Please check the Decode function ReqData is `schema` tag
func TestDecodeUrlQuery(t *testing.T) {
	// set the url query data
	g := Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		Request:        nil,
		ResponseWriter: nil,
		URLQuery: url.Values{
			"serverId": {strconv.Itoa(3)},
			"nginxDir": {"/usr/local/nginx"},
			"passwd":   {"123456Abc"},
			"fileName": {"test1"},
		},
		Body: nil,
	}

	Decode(t, &g)
}

// Test parsing data from body
// Please check the Decode function ReqData is `json` tag
func TestDecodeBody(t *testing.T) {
	// set the body data
	g := Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		Request:        nil,
		ResponseWriter: nil,
		URLQuery:       nil,
		Body: []byte(`{
			"serverId":3,
			"nginxDir":"/etc/init.d/nginx",
			"passwd":"123",
			"fileName":"test888"
		}`),
	}

	Decode(t, &g)
}

// Test parsing data from URL query and body
func TestDecodeBodyAndUrlQuery(t *testing.T) {
	// set the body and URL query data
	g := Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		Request:        nil,
		ResponseWriter: nil,
		URLQuery: url.Values{
			"serverId": {strconv.Itoa(3)},
			"nginxDir": {"/usr/local/nginx"},
			"passwd":   {"123456Abc"},
			"fileName": {"test1"},
		},
		Body: []byte(`{
			"serverId":4,
			"nginxDir":"/etc/init.d/nginx/",
			"passwd":"test888"
		}`),
	}

	Decode(t, &g)
}

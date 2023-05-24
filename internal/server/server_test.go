package server

import (
	"fmt"
	"github.com/zhenorzz/goploy/internal/model"
	"net/url"
	"strconv"
	"testing"
)

func TestDecodeUrlQuery(t *testing.T) {
	g := Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		Request:        nil,
		ResponseWriter: nil,
		URLQuery: url.Values{
			"serverId":    {strconv.Itoa(3)},
			"dir":         {"/usr/local/nginx"},
			"NewName":     {"test"},
			"CurrentName": {"test1"},
		},
		Body: nil,
	}

	Decode(t, &g)
}

func TestDecodeBody(t *testing.T) {
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
			"dir":"/etc/init.d/nginx",
			"currentName":"test666",
			"newName":"test888"
		}`),
	}

	Decode(t, &g)
}

func TestDecodeBodyAndUrlQuery(t *testing.T) {
	g := Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		Request:        nil,
		ResponseWriter: nil,
		URLQuery: url.Values{
			"serverId":    {strconv.Itoa(3)},
			"dir":         {"/usr/local/nginx"},
			"NewName":     {"test"},
			"CurrentName": {"test1"},
		},
		Body: []byte(`{
			"serverId":4,
			"dir":"/etc/init.d/nginx",
			"newName":"test888"
		}`),
	}

	Decode(t, &g)
}

func Decode(t *testing.T, g *Goploy) {
	type ReqData struct {
		ServerID    int64  `json:"serverId" validate:"gt=0"`
		Dir         string `json:"dir" validate:"filepath"`
		NewName     string `json:"newName" validate:"required"`
		CurrentName string `schema:"currentName" validate:"required"` // test schema when decode body
	}

	var reqData ReqData
	if err := g.Decode(&reqData); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", reqData)
}

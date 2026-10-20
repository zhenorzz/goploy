package server

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/zhenorzz/goploy/internal/model"
)

// queryReq is filled from URL query. schema.Decoder reads the schema tags.
type queryReq struct {
	ServerID int64  `schema:"serverId" validate:"gt=0"`
	Dir      string `schema:"nginxDir" validate:"filepath"`
	Password string `schema:"passwd" validate:"password"`
	File     string `schema:"fileName" validate:"required"`
}

// bodyReq is filled from JSON. encoding/json reads the json tags.
type bodyReq struct {
	ServerID int64  `json:"serverId" validate:"gt=0"`
	Dir      string `json:"nginxDir" validate:"filepath"`
	Password string `json:"passwd" validate:"password"`
	File     string `json:"fileName" validate:"required"`
}

// queryAndBodyReq accepts both steps of Goploy.Decode: query values are
// applied through schema tags, then JSON overwrites the same fields through
// json tags. Fields present only in the query are left in place.
type queryAndBodyReq struct {
	ServerID int64  `schema:"serverId" json:"serverId" validate:"gt=0"`
	Dir      string `schema:"nginxDir" json:"nginxDir" validate:"filepath"`
	Password string `schema:"passwd" json:"passwd" validate:"password"`
	File     string `schema:"fileName" json:"fileName" validate:"required"`
}

func decode(t *testing.T, g *Goploy, dst any) {
	t.Helper()
	if err := g.Decode(dst); err != nil {
		t.Fatal(err)
	}
}

func newGoploy(query url.Values, body []byte) *Goploy {
	return &Goploy{
		UserInfo: model.User{},
		Namespace: struct {
			ID            int64
			PermissionIDs map[int64]struct{}
		}{},
		URLQuery: query,
		Body:     body,
	}
}

func TestDecodeUrlQuery(t *testing.T) {
	g := newGoploy(url.Values{
		"serverId": {strconv.Itoa(3)},
		"nginxDir": {"/usr/local/nginx"},
		"passwd":   {"123456Abc"},
		"fileName": {"test1"},
	}, nil)

	var got queryReq
	decode(t, g, &got)
	if got.ServerID != 3 || got.Dir != "/usr/local/nginx" || got.Password != "123456Abc" || got.File != "test1" {
		t.Fatalf("decoded query = %+v", got)
	}
}

func TestDecodeBody(t *testing.T) {
	g := newGoploy(nil, []byte(`{
		"serverId":3,
		"nginxDir":"/etc/init.d/nginx",
		"passwd":"Abcd1234",
		"fileName":"test888"
	}`))

	var got bodyReq
	decode(t, g, &got)
	if got.ServerID != 3 || got.Dir != "/etc/init.d/nginx" || got.Password != "Abcd1234" || got.File != "test888" {
		t.Fatalf("decoded body = %+v", got)
	}
}

func TestDecodeBodyAndUrlQuery(t *testing.T) {
	g := newGoploy(url.Values{
		"serverId": {strconv.Itoa(3)},
		"nginxDir": {"/usr/local/nginx"},
		"passwd":   {"123456Abc"},
		"fileName": {"test1"},
	}, []byte(`{
		"serverId":4,
		"nginxDir":"/etc/init.d/nginx",
		"passwd":"test8888"
	}`))

	var got queryAndBodyReq
	decode(t, g, &got)
	// Body wins on fields it sets. fileName is only in the query.
	if got.ServerID != 4 || got.Dir != "/etc/init.d/nginx" || got.Password != "test8888" || got.File != "test1" {
		t.Fatalf("decoded query+body = %+v", got)
	}
}

package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zhenorzz/goploy/core"
)

// Common struct
type Common Controller

// Upload file
func (common Common) Upload(w http.ResponseWriter, gp *core.Goploy) {
	type RepData struct {
		FileName string `json:"fileName"`
	}
	file, handler, err := gp.Request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(core.GolbalPath+"tmp/"+handler.Filename, fileBytes, 0755)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	response := core.Response{Data: RepData{FileName: handler.Filename}}
	response.JSON(w)
}

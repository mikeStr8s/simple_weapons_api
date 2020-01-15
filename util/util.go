package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/valyala/fasthttp"
)

// SetResponse takes a context pointer and assignes the response context data for the API
func SetResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(200)
}

// ReadJSONFile reads the contents of a JSON file and returns an
// array of bytes of file data to be parsed into data object
func ReadJSONFile(dataset string) []byte {
	file, err := os.Open(path.Join(DIRNAME, "data", dataset+".json"))
	println(DIRNAME)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	return bytes
}

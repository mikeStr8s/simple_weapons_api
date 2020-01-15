package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"

	"github.com/mikeStr8s/simple_weapons_api/data"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Index is the base landing for routing
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// SetResponse takes a context pointer and assignes the response context data for the API
func SetResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(200)
}

// Lookup is a dynamic handler for lookup datasets
func Lookup(ctx *fasthttp.RequestCtx) {
	SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(data.GetData(path.Base(string(ctx.Path())))); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// Create is a dynamic handler for POST requests
func Create(ctx *fasthttp.RequestCtx) {
	SetResponse(ctx)
	if err := json.NewEncoder(ctx).Encode(data.PostData(path.Base(string(ctx.Path())), ctx.PostBody())); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func main() {
	router := router.New()
	router.GET("/", Index)

	api := router.Group("/api")
	api.GET("/abilityscore", Lookup)
	api.GET("/condition", Lookup)
	api.GET("/damage", Lookup)
	api.GET("/language", Lookup)
	api.GET("/movement", Lookup)
	api.GET("/sense", Lookup)
	api.GET("/skill", Lookup)
	api.GET("/movementspeed", Lookup)
	api.GET("/savingthrow", Lookup)
	api.GET("/sensevalue", Lookup)
	api.GET("/skillvalue", Lookup)
	api.GET("/monster", Lookup)

	api.POST("/movementspeed", Create)
	api.POST("/skillvalue", Create)
	api.POST("/savingthrow", Create)
	api.POST("/sensevalue", Create)
	api.POST("/monster", Create)

	log.Fatal(fasthttp.ListenAndServe(":1234", router.Handler))
}

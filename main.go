package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/mikeStr8s/simple_weapons_api/handlers"
	"github.com/mikeStr8s/simple_weapons_api/middleware"
	"github.com/valyala/fasthttp"
)

// Index is the base landing for routing
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

func main() {
	router := router.New()
	router.GET("/", Index)
	router.POST("/register", handlers.Register)

	api := router.Group("/api")
	api.GET("/abilityscore", handlers.Lookup)
	api.GET("/condition", handlers.Lookup)
	api.GET("/damage", handlers.Lookup)
	api.GET("/language", handlers.Lookup)
	api.GET("/movement", handlers.Lookup)
	api.GET("/sense", handlers.Lookup)
	api.GET("/skill", middleware.Auth(handlers.Lookup))
	api.GET("/movementspeed", handlers.Lookup)
	api.GET("/savingthrow", handlers.Lookup)
	api.GET("/sensevalue", handlers.Lookup)
	api.GET("/skillvalue", handlers.Lookup)
	api.GET("/monster", handlers.Lookup)

	api.POST("/movementspeed", handlers.Create)
	api.POST("/skillvalue", handlers.Create)
	api.POST("/savingthrow", handlers.Create)
	api.POST("/sensevalue", handlers.Create)
	api.POST("/monster", handlers.Create)

	log.Fatal(fasthttp.ListenAndServe(":1234", router.Handler))
}

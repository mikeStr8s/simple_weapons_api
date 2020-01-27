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
	router.POST("/login", handlers.Login)

	api := router.Group("/api")
	api.GET("/abilityscore", middleware.Auth(handlers.Lookup))
	api.GET("/condition", middleware.Auth(handlers.Lookup))
	api.GET("/damage", middleware.Auth(handlers.Lookup))
	api.GET("/language", middleware.Auth(handlers.Lookup))
	api.GET("/movement", middleware.Auth(handlers.Lookup))
	api.GET("/sense", middleware.Auth(handlers.Lookup))
	api.GET("/skill", middleware.Auth(handlers.Lookup))
	api.GET("/movementspeed", middleware.Auth(handlers.Lookup))
	api.GET("/savingthrow", middleware.Auth(handlers.Lookup))
	api.GET("/sensevalue", middleware.Auth(handlers.Lookup))
	api.GET("/skillvalue", middleware.Auth(handlers.Lookup))
	api.GET("/monster", middleware.Auth(handlers.Lookup))

	api.POST("/movementspeed", middleware.Auth(handlers.Create))
	api.POST("/skillvalue", middleware.Auth(handlers.Create))
	api.POST("/savingthrow", middleware.Auth(handlers.Create))
	api.POST("/sensevalue", middleware.Auth(handlers.Create))
	api.POST("/monster", middleware.Auth(handlers.Create))

	log.Fatal(fasthttp.ListenAndServe(":1234", router.Handler))
}

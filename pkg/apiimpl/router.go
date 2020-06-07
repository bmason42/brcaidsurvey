/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"context"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Run - configures and starts the web server
func RunServer() error {
	r := newRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}

func newRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/brcaid/v1")
	v1.Handle("GET", "/about", aboutGetUnversioned)

	addOpenApiDefRoutes(router)
	addSwaggerUIRoutes(router)
	addUnversionedRoutes(router)
	addWebUIRoutes(router)
	return router
}
func addOpenApiDefRoutes(router *gin.Engine) {
	router.StaticFile("/brcaid/openapi-1.yaml", "api/openapi-1.yaml")
	router.StaticFile("/brcaid/swagger.yaml", "api/openapi-1.yaml")
}

func addWebUIRoutes(router *gin.Engine) {
	webUI := static.LocalFile("web/", false)
	webHandler := static.Serve("/brcaid", webUI)
	router.Use(webHandler)
}
func addSwaggerUIRoutes(router *gin.Engine) {
	router.Handle("GET", "/brcaid/swaggerui/index.html", swaggerUIGetHandler)
	router.Handle("GET", "/brcaid/swaggerui", swaggerUIGetHandler)
	router.Handle("GET", "/brcaid/swaggerui/", swaggerUIGetHandler)
	swaggerUI := static.LocalFile("third_party/swaggerui/", false)
	webHandler := static.Serve("/brcaid/swaggerui/", swaggerUI)
	router.Use(webHandler)
}
func addUnversionedRoutes(router *gin.Engine) {
	router.Handle("GET", "/brcaid/about", aboutGetUnversioned)
	router.Handle("GET", "/brcaid/healthcheck", healthCheckGetUnversioned)
}

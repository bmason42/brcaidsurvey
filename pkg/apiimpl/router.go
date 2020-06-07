/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"brcaidsurvey/pkg/model"
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var noAuthRoutesSuffix []string

// Run - configures and starts the web server
func RunServer() error {
	noAuthRoutesSuffix = make([]string, 0)
	noAuthRoutesSuffix = append(noAuthRoutesSuffix, "/about")
	noAuthRoutesSuffix = append(noAuthRoutesSuffix, "/login")

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
func authSkip() bool {
	return os.Getenv("SKIP_AUTH") == "true"
}

func RouteAuthorized(c *gin.Context) {

	if authSkip() {
		return
	}
	atLeastOnePassed := false
	for _, suffix := range noAuthRoutesSuffix {
		if strings.HasSuffix(c.Request.URL.Path, suffix) {
			atLeastOnePassed = true
			break
		}
	}
	//path was not a no auth path, so check session
	if !atLeastOnePassed {
		authToken := c.Request.Header.Get("Authorization")
		session := model.LookupSession(authToken)
		fmt.Println(session)
		//todo, check against auth table of paths... And make that auth table
		atLeastOnePassed = true
	}
	if atLeastOnePassed {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, "")
		c.Abort()
	}

	return
}
func newRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/brcaid/v1", RouteAuthorized)
	v1.Handle("GET", "/about", aboutGetUnversioned)
	v1.Handle("POST", "/login", loginHandler)
	v1.Handle("POST", "/logout", logoutHandler)
	v1.Handle("GET", "/users", userGetHandler)
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

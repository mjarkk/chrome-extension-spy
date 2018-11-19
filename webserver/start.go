package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mjarkk/chrome-extension-spy/types"
)

// StartWebServer starts the web serve
func StartWebServer(tmpDir string, forceClose chan struct{}, extenisons map[string]*types.FullAndSmallExt) error {
	extsMap = extenisons
	gin.SetMode("release")
	r := gin.Default()
	newReq := make(chan types.SmallRequest)
	go func() {
		setupSockServer(r, newReq)
	}()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/proxy/:appid/:url", func(c *gin.Context) {
		proxyHandeler(c, "GET", newReq)
	})
	r.POST("/proxy/:appid/:url", func(c *gin.Context) {
		proxyHandeler(c, "POST", newReq)
	})
	r.GET("/lastRequests", lastRequests)
	r.GET("/extensionsInfo", extensionsInfo)
	r.GET("/requestInfo/:id", requestInfo)
	r.GET("/extLogo/:extID", func(c *gin.Context) {
		extLogo(c, tmpDir)
	})
	r.Static("/js/", "./web_static/build/js/")
	r.StaticFile("/", "./web_static/build/index.html")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-forceClose
	fmt.Println("Trying to stop web server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Unable to shutdown the server run CTRL+C to force quit")
	}
	return nil
}

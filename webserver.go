package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func proxyHandeler(c *gin.Context, reqType string) {
	rawURL := c.Param("url")
	parsedURL, err := url.PathUnescape(rawURL)
	if err != nil {
		c.String(http.StatusConflict, "")
	}

	hc := http.Client{}
	req, err := http.NewRequest(reqType, parsedURL, nil)

	if reqType == "POST" {
		req.Body = c.Request.Body
	}

	for key, value := range c.Request.Header {
		req.Header.Add(key, value[0])
	}

	rs, err := hc.Do(req)
	if err != nil {
		c.String(400, "")
		return
	}

	for key, item := range rs.Header {
		c.Header(key, item[0])
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		body = []byte("")
	}

	c.Data(rs.StatusCode, rs.Header.Get("Content-Type"), body)
}

func proxyHandelerPost(c *gin.Context) {
	proxyHandeler(c, "POST")
}

func proxyHandelerGet(c *gin.Context) {
	proxyHandeler(c, "GET")
}

func startWebServer(closeServer chan struct{}) error {
	gin.SetMode("release")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/proxy/:url", proxyHandelerGet)
	r.POST("/proxy/:url", proxyHandelerPost)
	r.Static("/web_static", "./web_static")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-closeServer
	fmt.Println("Stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	return nil
}

package webserver

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mjarkk/chrome-extension-spy/types"
)

var globalData = []types.Request{}

// extsMap has both the full extension and the in a map. This makes it easy to select an extension
var extsMap map[string]*types.FullAndSmallExt

func proxyHandeler(c *gin.Context, reqType string) {

	var dataToSave types.Request
	defer func() {
		globalData = append(globalData, dataToSave)
	}()

	dataToSave.Type = reqType

	appID := c.Param("appid")
	dataToSave.Extension = extsMap[appID].Small

	rawURL := c.Param("url")
	parsedURL, err := url.PathUnescape(rawURL)
	dataToSave.URL = parsedURL
	if err != nil {
		dataToSave.StatusCode = http.StatusConflict
		dataToSave.ResData = ""
		dataToSave.ResData64 = ""
		c.String(http.StatusConflict, "")
		return
	}

	hc := http.Client{}
	req, err := http.NewRequest(reqType, parsedURL, nil)

	if reqType == "POST" {
		req.Body = c.Request.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		dataToSave.PostBody = buf.String()
	}

	for key, value := range c.Request.Header {
		req.Header.Add(key, value[0])
		dataToSave.RequestHeaders[key] = value[0]
	}

	rs, err := hc.Do(req)
	if err != nil {
		sCode := 400
		dataToSave.StatusCode = sCode
		dataToSave.ResData = ""
		dataToSave.ResData64 = ""
		c.String(sCode, "")
		return
	}

	for key, item := range rs.Header {
		c.Header(key, item[0])
		dataToSave.ResponseHeaders[key] = item[0]
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		body = []byte("")
	}

	contentType := rs.Header.Get("Content-Type")
	dataToSave.ResponseHeaders["Content-Type"] = contentType
	dataToSave.ResData = string(body)
	dataToSave.StatusCode = rs.StatusCode

	c.Data(rs.StatusCode, contentType, body)
}

func proxyHandelerPost(c *gin.Context) {
	proxyHandeler(c, "POST")
}

func proxyHandelerGet(c *gin.Context) {
	proxyHandeler(c, "GET")
}

func lastRequests(c *gin.Context) {
	c.JSON(200, globalData)
}

// StartWebServer starts the web serve
func StartWebServer(forceClose chan struct{}, extenisons map[string]*types.FullAndSmallExt) error {
	extsMap = extenisons
	gin.SetMode("release")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/proxy/:appid/:url", proxyHandelerGet)
	r.POST("/proxy/:appid/:url", proxyHandelerPost)
	r.GET("/lastRequests", lastRequests)
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

package webserver

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
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
		hashData := fmt.Sprintf("%v", dataToSave)
		hashData += string(time.Now().Unix())
		dataToSave.Hash = fmt.Sprintf("%x", sha1.Sum([]byte(hashData)))
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

	dataToSave.RequestHeaders = make(map[string]string, len(c.Request.Header))
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

	dataToSave.ResponseHeaders = make(map[string]string, len(rs.Header))
	for key, item := range rs.Header {
		c.Header(key, item[0])
		dataToSave.ResponseHeaders[key] = item[0]
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		body = []byte("")
	}

	dataToSave.ResData = string(body)
	dataToSave.StatusCode = rs.StatusCode

	c.Data(rs.StatusCode, rs.Header.Get("Content-Type"), body)
}

func proxyHandelerPost(c *gin.Context) {
	proxyHandeler(c, "POST")
}

func proxyHandelerGet(c *gin.Context) {
	proxyHandeler(c, "GET")
}

func extensionsInfo(c *gin.Context) {
	c.JSON(200, extsMap)
}

func lastRequests(c *gin.Context) {
	res := make([]types.SmallRequest, len(globalData))
	for i, item := range globalData {
		res[i].Pkg = item.Extension.Pkg
		res[i].Type = item.Type
		res[i].Code = item.StatusCode
		res[i].URL = item.URL
		res[i].Hash = item.Hash
	}
	c.JSON(200, res)
}

func extLogo(c *gin.Context, tmpDir string) {
	data, exists := extsMap[c.Param("extID")]
	if !exists {
		c.Data(400, "image/jpeg", []byte(""))
		return
	}

	icon := ""
	i := data.Full.Icons
	if len(i.Num128) > 3 {
		icon = i.Num128
	} else if len(i.Num64) > 3 {
		icon = i.Num64
	} else if len(i.Num48) > 3 {
		icon = i.Num48
	} else if len(i.Num32) > 3 {
		icon = i.Num32
	} else if len(i.Num16) > 3 {
		icon = i.Num16
	} else {
		c.Data(400, "image/jpeg", []byte(""))
		return
	}

	buff, err := ioutil.ReadFile(path.Join(tmpDir, data.Small.Pkg, icon))

	if err != nil {
		c.Data(400, "image/jpeg", []byte(""))
		return
	}

	c.Data(200, http.DetectContentType(buff), buff)
}

func requestInfo(c *gin.Context) {
	toSearchFor := c.Param("id")
	var res types.Request
	for _, req := range globalData {
		if req.Hash == toSearchFor {
			res = req
		}
	}
	c.JSON(200, res)
}

// StartWebServer starts the web serve
func StartWebServer(tmpDir string, forceClose chan struct{}, extenisons map[string]*types.FullAndSmallExt) error {
	extsMap = extenisons
	gin.SetMode("release")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/proxy/:appid/:url", proxyHandelerGet)
	r.POST("/proxy/:appid/:url", proxyHandelerPost)
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

package webserver

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mjarkk/chrome-extension-spy/types"
)

func proxyHandeler(c *gin.Context, reqType string, newReq chan types.SmallRequest) {
	var dataToSave types.Request
	defer func() {
		hashData := fmt.Sprintf("%v", dataToSave)
		hashData += string(time.Now().Unix())
		dataToSave.Hash = fmt.Sprintf("%x", sha1.Sum([]byte(hashData)))
		globalData = append(globalData, dataToSave)
		newReq <- mkSmallReq(dataToSave)
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
		dataToSave.ResRawData = fmt.Sprintf("%s", hex.Dump([]byte{}))
		dataToSave.ProxyErrors = append(dataToSave.ProxyErrors, "Not beable to parse url")
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
		dataToSave.ResRawData = fmt.Sprintf("%s", hex.Dump([]byte{}))
		c.String(sCode, "")
		dataToSave.ProxyErrors = append(dataToSave.ProxyErrors, "Error while trying to send request")
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
	dataToSave.ResRawData = fmt.Sprintf("%s", hex.Dump(body))
	dataToSave.StatusCode = rs.StatusCode

	c.Data(rs.StatusCode, rs.Header.Get("Content-Type"), body)
}

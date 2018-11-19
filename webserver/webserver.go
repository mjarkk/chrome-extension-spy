package webserver

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mjarkk/chrome-extension-spy/types"
)

var globalData = []types.Request{}

// extsMap has both the full extension and the in a map. This makes it easy to select an extension
var extsMap map[string]*types.FullAndSmallExt

func extensionsInfo(c *gin.Context) {
	c.JSON(200, extsMap)
}

func mkSmallReq(item types.Request) types.SmallRequest {
	var toReturn types.SmallRequest
	toReturn.Pkg = item.Extension.Pkg
	toReturn.Type = item.Type
	toReturn.Code = item.StatusCode
	toReturn.URL = item.URL
	toReturn.Hash = item.Hash
	return toReturn
}

func lastRequests(c *gin.Context) {
	res := make([]types.SmallRequest, len(globalData))
	for i, item := range globalData {
		res[i] = mkSmallReq(item)
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

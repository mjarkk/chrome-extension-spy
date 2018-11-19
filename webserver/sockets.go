package webserver

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mjarkk/chrome-extension-spy/types"

	"gopkg.in/olahol/melody.v1"
)

func setupSockServer(r *gin.Engine, newReq chan types.SmallRequest) {
	m := melody.New()

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	for true {
		data := <-newReq
		toSendData, err := json.Marshal(data)
		if err == nil {
			m.Broadcast(toSendData)
		}
	}
}

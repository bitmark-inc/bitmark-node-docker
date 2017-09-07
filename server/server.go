package server

import (
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/gin-gonic/gin"
)

type ServiceOptionRequest struct {
	Option string `json:"option"`
}

type WebServer struct {
	Bitmarkd services.Service
	Prooferd services.Service
}

func NewWebServer(bitmarkd, prooferd services.Service) *WebServer {
	return &WebServer{
		Bitmarkd: bitmarkd,
		Prooferd: prooferd,
	}
}

func (ws *WebServer) BitmarkdStartStop(c *gin.Context) {
	var req ServiceOptionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = nil
	switch req.Option {
	case "start":
		err = ws.Bitmarkd.Start()
	case "stop":
		err = ws.Bitmarkd.Stop()
	case "status":
		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": ws.Bitmarkd.IsRunning(),
		})
		return
	default:
		c.String(400, "invalid option")
		return
	}

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"ok":  0,
			"msg": err.Error(),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"ok": 1,
		})
	}
}

func (ws *WebServer) ProoferdStartStop(c *gin.Context) {
	var req ServiceOptionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = nil
	switch req.Option {
	case "start":
		err = ws.Prooferd.Start()
	case "stop":
		err = ws.Prooferd.Stop()
	case "status":
		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": ws.Prooferd.IsRunning(),
		})
		return
	default:
		c.String(400, "invalid option")
		return
	}

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"ok":  0,
			"msg": err.Error(),
		})
	} else {
		c.JSON(200, map[string]interface{}{
			"ok": 1,
		})
	}
}

func (ws *WebServer) DiscoveryStartStop(c *gin.Context) {

}

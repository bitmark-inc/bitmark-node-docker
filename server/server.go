package server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/bitmark-inc/bitmark-node/config"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/bitmarkd/rpc"
	"github.com/gin-gonic/gin"
)

const CONFIG_BUCKET_NAME = "config"

var client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type ServiceOptionRequest struct {
	Option string `json:"option"`
}

type WebServer struct {
	nodeConfig *config.BitmarkNodeConfig
	rootPath   string
	Bitmarkd   services.Service
	Recorderd  services.Service
}

func NewWebServer(nc *config.BitmarkNodeConfig, rootPath string, bitmarkd, recorderd services.Service) *WebServer {
	return &WebServer{
		nodeConfig: nc,
		rootPath:   rootPath,
		Bitmarkd:   bitmarkd,
		Recorderd:  recorderd,
	}
}

func (ws *WebServer) GetChain(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()

	c.SetCookie("bitmark-node-network", network, 0, "", "", false, false)
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": network,
	})
	return
}

// func (ws *WebServer) SetChain(c *gin.Context) {
// 	reqBody := map[string]string{}
// 	err := c.BindJSON(&reqBody)
// 	if err != nil {
// 		c.String(400, "can not parse action option")
// 		return
// 	}

// 	network, ok := reqBody["network"]
// 	if !ok {
// 		c.String(400, "missing arguments")
// 		return
// 	}

// 	err = ws.nodeConfig.SetNetwork(network)
// 	if err != nil {
// 		c.String(400, "can not set network. error: %s", err.Error())
// 		return
// 	}

// 	ws.Bitmarkd.SetNetwork(network)
// 	ws.Recorderd.SetNetwork(network)

// 	c.JSON(200, map[string]interface{}{
// 		"ok": 1,
// 	})
// 	return
// }

func (ws *WebServer) GetConfig(c *gin.Context) {
	config, err := ws.nodeConfig.Get()
	if err != nil {
		c.String(500, "can not read bitmark node config. error: %s", err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": config,
	})
	return
}

func (ws *WebServer) UpdateConfig(c *gin.Context) {
	newConfig := map[string]string{
		"btcAddr": "",
		"ltcAddr": "",
	}

	err := c.BindJSON(&newConfig)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = ws.nodeConfig.Set(newConfig)

	if err != nil {
		c.String(500, "can not set bitmark node config. error: %s", err.Error())
		return
	}

	c.String(200, "")
	return
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
			"result": ws.Bitmarkd.Status(),
		})
		return
	case "info":

		resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/info")
		if err != nil {
			c.String(500, "unable to get bitmark info")
			return
		}
		defer resp.Body.Close()
		bb := bytes.Buffer{}
		io.Copy(&bb, resp.Body)

		if resp.StatusCode != 200 {
			c.String(500, "unable to get bitmark info. message: %s", bb.String())
			return
		}

		var reply rpc.InfoReply
		d := json.NewDecoder(&bb)

		if err := d.Decode(&reply); err != nil {
			c.String(500, "fail to read bitmark info response. error: %s\n", err.Error())
			return
		}

		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": reply,
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

func (ws *WebServer) RecorderdStartStop(c *gin.Context) {
	var req ServiceOptionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = nil
	switch req.Option {
	case "start":
		err = ws.Recorderd.Start()
	case "stop":
		err = ws.Recorderd.Stop()
	case "status":
		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": ws.Recorderd.Status(),
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

package server

import (
	"crypto/tls"
	"net/rpc/jsonrpc"

	"github.com/bitmark-inc/bitmark-node/config"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/bitmarkd/rpc"
	"github.com/gin-gonic/gin"
)

const CONFIG_BUCKET_NAME = "config"

type ServiceOptionRequest struct {
	Option string `json:"option"`
}

type WebServer struct {
	nodeConfig *config.BitmarkNodeConfig
	Bitmarkd   services.Service
	Prooferd   services.Service
}

func NewWebServer(nc *config.BitmarkNodeConfig, bitmarkd, prooferd services.Service) *WebServer {
	return &WebServer{
		nodeConfig: nc,
		Bitmarkd:   bitmarkd,
		Prooferd:   prooferd,
	}
}

func (ws *WebServer) SetChain(c *gin.Context) {
	reqBody := map[string]string{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	network, ok := reqBody["network"]
	if !ok {
		c.String(400, "missing arguments")
		return
	}

	err = ws.nodeConfig.SetNetwork(network)
	if err != nil {
		c.String(400, "can not set network. error: %s", err.Error())
		return
	}

	ws.Bitmarkd.SetNetwork(network)
	ws.Prooferd.SetNetwork(network)

	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
	return
}

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
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		conn, err := tls.Dial("tcp", "127.0.0.1:2130", tlsConfig)
		if nil != err {
			c.String(500, "can not parse action option")
			return
		}

		client := jsonrpc.NewClient(conn)
		defer client.Close()

		var reply rpc.InfoReply
		if err := client.Call("Node.Info", rpc.InfoArguments{}, &reply); err != nil {
			c.String(500, "Node.Info error: %s\n", err.Error())
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
			"result": ws.Prooferd.Status(),
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

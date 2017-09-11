package server

import (
	"crypto/tls"
	"log"
	"net/rpc/jsonrpc"

	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/bitmarkd/rpc"
	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
)

const CONFIG_BUCKET_NAME = "config"

type ServiceOptionRequest struct {
	Option string `json:"option"`
}

type WebServer struct {
	db       *bolt.DB
	Bitmarkd services.Service
	Prooferd services.Service
}

func NewWebServer(dbPath string, bitmarkd, prooferd services.Service) *WebServer {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CONFIG_BUCKET_NAME))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return &WebServer{
		db:       db,
		Bitmarkd: bitmarkd,
		Prooferd: prooferd,
	}
}

func (ws *WebServer) GetConfig(c *gin.Context) {
	config := map[string]string{
		"btcAddr": "",
		"ltcAddr": "",
	}

	ws.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
		for key := range config {
			b := bucket.Get([]byte(key))
			if b != nil {
				config[key] = string(b)
			}
		}
		return nil
	})
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": config,
	})
	return
}

func (ws *WebServer) UpdateConfig(c *gin.Context) {
	config := map[string]string{
		"btcAddr": "",
		"ltcAddr": "",
	}

	err := c.BindJSON(&config)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	ws.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(CONFIG_BUCKET_NAME))
		for key, val := range config {
			err := bucket.Put([]byte(key), []byte(val))
			if err != nil {
				return err
			}
		}
		return nil
	})

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

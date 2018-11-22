package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bitmark-inc/bitmark-node/config"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/logger"
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

type AccountInfo struct {
	network       string
	accountNumber string
	seed          string
}

type WebServer struct {
	Mutex             *sync.Mutex
	nodeConfig        *config.BitmarkNodeConfig
	rootPath          string
	log               *logger.L
	peerPortReachable bool
	Bitmarkd          services.Service
	Recorderd         services.Service
	Accounts          []AccountInfo
}

func NewWebServer(nc *config.BitmarkNodeConfig, rootPath string, bitmarkd, recorderd services.Service) *WebServer {
	return &WebServer{
		Mutex:      &sync.Mutex{},
		nodeConfig: nc,
		rootPath:   rootPath,
		log:        logger.New("webserver"),
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

		resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/details")
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

		var reply DetailReply
		d := json.NewDecoder(&bb)

		if err := d.Decode(&reply); err != nil {
			c.String(500, "fail to read bitmark info response. error: %s\n", err.Error())
			return
		}

		t, _ := time.ParseDuration(reply.Uptime)
		reply.Uptime = t.Round(time.Second).String()

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

func (ws *WebServer) GetLog(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	logFile := filepath.Join(ws.rootPath, c.Param("serviceName"), network, "log", c.Param("serviceName")+".log")

	file, err := os.Open(logFile)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lastLine := ""
	c.Header("X-Content-Type-Options", "nosniff")
	for {
		if !scanner.Scan() {
			fmt.Fprintln(c.Writer, lastLine)
			c.Writer.Flush()
			break
		}
		lastLine = scanner.Text()
	}

	reader := bufio.NewReader(file)
	c.Stream(func(w io.Writer) bool {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(1 * time.Second)
			} else {
				fmt.Fprintf(w, "===== log stopped with error: %s", err.Error())
				return false
			}
		}
		fmt.Fprint(w, line)
		return true
	})
}

func (ws *WebServer) GetAccountNumber(network string) (string, error) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network && item.accountNumber != "" {
			return item.accountNumber, nil
		}
	}

	return "", errors.New("No account in " + network + " network")

}

func (ws *WebServer) GetSeed(network string) (string, error) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network && item.seed != "" {
			return item.seed, nil
		}
	}
	return "", errors.New("No seed of AccountInfo in " + network + " network")
}

func (ws *WebServer) SetAccount(accountNumber, seed, network string) error {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network {
			item.accountNumber = accountNumber
			item.seed = seed
			return nil
		}
	}
	ws.log.Warnf("Account SetAccount:number:%s seed:%s network:%s", accountNumber, seed, network)
	ws.Accounts = append(ws.Accounts, AccountInfo{accountNumber: accountNumber, seed: seed, network: network})
	return nil
}

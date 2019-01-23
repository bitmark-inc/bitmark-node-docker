package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const retryDelay = time.Duration(500 * time.Millisecond)
const retryTimes = 3
const checkInterMs = 2000
const dialTimeout = 2 * time.Second
const bitmarkdDetailApi = "bitmarkd/details"
const bitmarkdDetailHost = "127.0.0.1"
const bitmarkdDetailPort = "2131"

//CheckPortReachableRoutine is a Connection Check Routine
func (ws *WebServer) CheckPortReachableRoutine(host, port string) {
	status := make(chan bool)
	for {
		go func(updateStatus chan<- bool) {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					time.Sleep(retryDelay)
				} else {
					retry = retryTimes + 1
				}
			}
			updateStatus <- connected

		}(status)

		ws.peerPortReachable = <-status
		time.Sleep(time.Duration(checkInterMs) * time.Millisecond)
	}
}

func connToPort(host, port string) bool {
	if host == "" {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), dialTimeout)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}

func getPeers(api string) (int, int) {
	resp, err := client.Get(api)
	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)

	if resp.StatusCode != 200 {
		return 0, 0
	}

	var reply DetailReply
	d := json.NewDecoder(&buf)
	err = d.Decode(&reply)
	if err != nil {
		return 0, 0
	}

	return int(reply.Peers.Incoming), int(reply.Peers.Outgoing)
}

//IsPeerPortReachable returns current status if current peer port is reachable
func (ws *WebServer) IsPeerPortReachable() bool {
	return ws.peerPortReachable
}

//ConnectionStatus return connected node number and Peer port reachability
func (ws *WebServer) ConnectionStatus(c *gin.Context) {
	incoming, outgoing := getPeers(ws.GetBitmarkdDetailApi())
	c.JSON(200, map[string]interface{}{
		"incoming": incoming,
		"outgoing": outgoing,
		"port_state": map[string]interface{}{
			"listening": ws.IsPeerPortReachable(),
		},
	})
	return
}

// GetBitmarkdDetailApi get bitmarkd detail api
func (ws *WebServer) GetBitmarkdDetailApi() string {
	httpRPCPort := os.Getenv("HTTP_RPC_PORT")
	if len(httpRPCPort) == 0 {
		httpRPCPort = bitmarkdDetailPort
	}

	return fmt.Sprintf("https://%s:%s/%s", bitmarkdDetailHost, httpRPCPort, bitmarkdDetailApi)
}

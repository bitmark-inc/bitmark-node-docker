package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

const retryDelay = time.Duration(500 * time.Millisecond)

//CheckPortReachableRoutine is a Connection Check Routine
func (ws *WebServer) CheckPortReachableRoutine(host, port string) {
	stop := make(chan bool)
	defer close(stop)
	for {
		connStat := connCheck(host, port, 1000, 3, stop)
		for {
			ws.peerPortReachable = <-connStat
		}

	}
}

func connCheck(host, port string, checkInterMs, retryTimes int, done <-chan bool) <-chan bool {
	status := make(chan bool)
	defer close(status)

	go func() {
		for {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					retry++
					time.Sleep(retryDelay)
				} else {
					retry = retryTimes + 1
				}
			}
			select {
			case <-done:
				return
			case status <- connected:
			}
			time.Sleep(time.Duration(checkInterMs) * time.Millisecond)
		}
	}()
	return status
}

func connToPort(host, port string) bool {
	if host == "" {
		return false
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return false
	} else {
		defer conn.Close()
		return true
	}
}

func getConnectors() int {
	resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/details")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)

	if resp.StatusCode != 200 {
		return 0
	}

	var reply DetailReply
	d := json.NewDecoder(&buf)
	err = d.Decode(&reply)
	if err != nil {
		return 0
	}

	return int(reply.Peers.Local)
}

//IsPeerPortReachable returns current status if current peer port is reachable
func (ws *WebServer) IsPeerPortReachable() bool {
	return ws.peerPortReachable
}

//ConnectionStatus return connected node number and Peer port reachability
func (ws *WebServer) ConnectionStatus(c *gin.Context) {
	//
	c.JSON(200, map[string]interface{}{
		"connections": getConnectors(),
		"port_state": map[string]interface{}{
			"listening": ws.IsPeerPortReachable(),
		},
	})
	return
}

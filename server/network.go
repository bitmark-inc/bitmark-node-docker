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

const retryDelay = time.Duration(200 * time.Millisecond)

var connectionStatus = false

func connCheckRoutine(host, port string, checkDelay, retryTimes int, done <-chan bool) <-chan bool {
	status := make(chan bool)
	go func() {
		for {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connCheck(host, port)
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

			time.Sleep(time.Duration(checkDelay) * time.Millisecond)
		}
	}()

	return status
}

func connCheck(host, port string) bool {
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

func (ws *WebServer) ConnectionStatus(c *gin.Context) {

	publicIP := os.Getenv("PUBLIC_IP")

	c.JSON(200, map[string]interface{}{
		"connections": getConnectors(),
		"port_state": map[string]interface{}{
			"listening": connCheck(publicIP, "2136"),
		},
	})
	return
}

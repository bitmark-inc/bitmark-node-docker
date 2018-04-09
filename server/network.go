package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

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

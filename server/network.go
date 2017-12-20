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
	resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/info")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)

	if resp.StatusCode != 200 {
		return 0
	}

	var result map[string]interface{}
	d := json.NewDecoder(&buf)
	err = d.Decode(&result)
	if err != nil {
		return 0
	}
	if count, ok := result["client_count"]; !ok {
		return 0
	} else {
		return int(count.(float64))
	}
}

func (ws *WebServer) ConnectionStatus(c *gin.Context) {

	publicIP := os.Getenv("PUBLIC_IP")

	c.JSON(200, map[string]interface{}{
		"connections": getConnectors(),
		"port_state": map[string]interface{}{
			"broadcast": connCheck(publicIP, "2135"),
			"listening": connCheck(publicIP, "2136"),
		},
	})
	return
}

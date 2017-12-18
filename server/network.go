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

func (ws *WebServer) ConnectionStatus(c *gin.Context) {
	resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/info/connectors")
	if err != nil {
		c.String(500, "unable to get connectors")
		return
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)

	if resp.StatusCode != 200 {
		c.String(500, "unable to get bitmark info. message: %s", buf.String())
		return
	}

	var result map[string][]map[string]string
	d := json.NewDecoder(&buf)
	err = d.Decode(&result)
	if err != nil {
		c.String(500, "fail to read bitmark connector response. error: %s\n", err.Error())
		return
	}

	publicIP := os.Getenv("PUBLIC_IP")

	c.JSON(200, map[string]interface{}{
		"connections": len(result["clients"]),
		"port_state": map[string]interface{}{
			"broadcast": connCheck(publicIP, "2135"),
			"listening": connCheck(publicIP, "2136"),
		},
	})
	return
}

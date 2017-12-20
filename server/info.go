package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bitmark-inc/bitmarkd/account"
	"github.com/gin-gonic/gin"
)

type DockerTag struct {
	Name string `json:"name"`
}

type DockerTagsResponse struct {
	Results []DockerTag `json:"results"`
}

func latestVersion() string {
	resp, err := client.Get("https://hub.docker.com/v2/repositories/bitmark/bitmark-node/tags/")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result DockerTagsResponse
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&result)
	if err != nil {
		return ""
	}

	var newVersion float64
	for _, tags := range result.Results {
		v, err := strconv.ParseFloat(tags.Name[1:], 64)
		if err != nil {
			continue
		}
		if v > newVersion {
			newVersion = v
		}
	}

	return fmt.Sprintf("v%0.1f", newVersion)
}

func (ws *WebServer) NodeInfo(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		c.String(500, "wrong network configuration")
		return
	}

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.Open(seedFile)
	if err != nil {
		c.String(500, "fail to open seed file from: %s", seedFile)
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		c.String(500, "can not read config file")
		return
	}
	seed := strings.Trim(strings.Split(buf.String(), ":")[1], "\n")
	a, err := account.PrivateKeyFromBase58Seed(seed)
	if err != nil {
		c.String(500, "unable to get your account. error: %s", err.Error())
		return
	}
	c.SetCookie("bitmark-node-network", network, 0, "", "", false, false)
	c.JSON(200, map[string]interface{}{
		"ok": 1,
		"result": map[string]string{
			"version": os.Getenv("VERSION"),
			"network": network,
			"account": a.Account().String(),
		},
	})
}

func (ws *WebServer) LatestVersion(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"latest": latestVersion(),
	})
	return
}

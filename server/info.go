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

var defaultVersion string = "v0.0.1" // do not change this value
type DockerTag struct {
	Name string `json:"name"`
}

type DockerTagsResponse struct {
	Results []DockerTag `json:"results"`
}

func latestVersion() string {
	resp, err := client.Get("https://hub.docker.com/v2/repositories/bitmark/bitmark-node-docker/tags/")
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

// NodeInfo Return NodeInformation
func (ws *WebServer) NodeInfo(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		c.String(500, ErrorNoNetwork.Error())
		return
	}

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.Open(seedFile)
	if err != nil {
		c.String(500, ErrorOpenSeedFile.Error()+seedFile)
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		c.String(500, ErrorReadConfigFile.Error())
		return
	}
	seed := strings.Trim(strings.Split(buf.String(), ":")[1], "\n")
	a, err := account.PrivateKeyFromBase58Seed(seed)
	if err != nil {
		c.String(500, ErrCombind(ErrorSetAccount, err).Error())
		return
	}
	nodeVersion := os.Getenv("VERSION")
	if len(nodeVersion) == 0 {
		nodeVersion = defaultVersion
		ws.log.Warnf("use default  bitmark node version =%s\n", nodeVersion)
	}
	c.SetCookie("bitmark-node-docker-network", network, 0, "", "", false, false)
	c.JSON(200, map[string]interface{}{
		"ok": 1,
		"result": map[string]string{
			"version": nodeVersion,
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

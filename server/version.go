package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DockerTag struct {
	Name string `json:"name"`
}

type DockerTagsResponse struct {
	Results []DockerTag `json:"results"`
}

func lastVersion() string {
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

func (ws *WebServer) Version(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"current": os.Getenv("VERSION"),
		"latest":  lastVersion(),
	})
	return
}

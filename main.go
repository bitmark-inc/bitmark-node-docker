package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bitmark-inc/bitmark-node/server"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/logger"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl"
)

type Configuration struct {
	Port    int                  `hcl:"port"`
	DataDir string               `hcl:"datadir"`
	Logging logger.Configuration `hcl:"logging"`
}

func (c *Configuration) Parse(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	io.Copy(&buf, f)
	return hcl.Unmarshal(buf.Bytes(), c)
}

func main() {
	defer exitwithstatus.Handler()

	var config Configuration
	err := config.Parse("bitmark-node.conf")
	if err != nil {
		exitwithstatus.Message(err.Error())
	}

	err = logger.Initialise(config.Logging)
	if err != nil {
		exitwithstatus.Message(err.Error())
	}
	defer logger.Finalise()
	rootPath := path.Clean(config.DataDir)

	bitmarkdPath := path.Join(rootPath, "bitmarkd1")
	prooferdPath := path.Join(rootPath, "prooferd")

	err = os.MkdirAll(bitmarkdPath, 0755)
	err = os.MkdirAll(prooferdPath, 0755)

	bitmarkdService := services.NewBitmarkd()
	prooferdService := services.NewProoferd()
	bitmarkdService.Initialise(path.Join(bitmarkdPath, "bitmarkd.conf"))
	defer bitmarkdService.Finalise()
	prooferdService.Initialise(path.Join(prooferdPath, "prooferd.conf"))
	defer prooferdService.Finalise()

	webserver := server.NewWebServer(bitmarkdService, prooferdService)

	r := gin.New()
	apiRouter := r.Group("/api")
	apiRouter.POST("/bitmarkd", webserver.BitmarkdStartStop)
	apiRouter.POST("/prooferd", webserver.ProoferdStartStop)
	r.Run(fmt.Sprintf(":%d", config.Port))
}

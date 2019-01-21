package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bitmark-inc/bitmark-node/config"
	"github.com/bitmark-inc/bitmark-node/server"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/logger"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl"
)

var version string = "v0.1" // do not change this value

type MasterConfiguration struct {
	Port       int                  `hcl:"port"`
	DataDir    string               `hcl:"datadir"`
	Logging    logger.Configuration `hcl:"logging"`
	VersionURL string               `hcl:"versionURL"`
}

func (c *MasterConfiguration) Parse(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	io.Copy(&buf, f)
	return hcl.Unmarshal(buf.Bytes(), c)
}

func init() {
	os.Setenv("VERSION", version)
}

func main() {
	defer exitwithstatus.Handler()

	var confFile string
	var containerIP string
	var uiPath string
	flag.StringVar(&confFile, "config-file", "bitmark-node.conf", "configuration for bitmark-node")
	flag.StringVar(&containerIP, "container-ip", "", "ip address for container")
	flag.StringVar(&uiPath, "ui", "ui/public", "path of ui interface")
	flag.Parse()

	var masterConfig MasterConfiguration
	err := masterConfig.Parse(confFile)
	if err != nil {
		exitwithstatus.Message(err.Error())
	}

	err = logger.Initialise(masterConfig.Logging)
	if err != nil {
		exitwithstatus.Message(err.Error())
	}
	defer logger.Finalise()
	var rootPath string
	if filepath.IsAbs(masterConfig.DataDir) {
		rootPath = masterConfig.DataDir
	} else {
		rootPath, err = filepath.Abs(filepath.Join(filepath.Dir(confFile), masterConfig.DataDir))
		if err != nil {
			exitwithstatus.Message(err.Error())
		}
	}

	bitmarkdPath := filepath.Join(rootPath, "bitmarkd")
	recorderdPath := filepath.Join(rootPath, "recorderd")
	dbPath := filepath.Join(rootPath, "db")

	err = os.MkdirAll(bitmarkdPath, 0755)
	err = os.MkdirAll(recorderdPath, 0755)
	err = os.MkdirAll(dbPath, 0755)

	bitmarkdService := services.NewBitmarkd(containerIP)
	recorderdService := services.NewRecorderd()
	bitmarkdService.Initialise(bitmarkdPath)
	defer bitmarkdService.Finalise()
	recorderdService.Initialise(recorderdPath)
	defer recorderdService.Finalise()

	nodeConfig := config.New()
	err = nodeConfig.Initialise(dbPath)
	if err != nil {
		exitwithstatus.Message(err.Error())
	}

	if network := nodeConfig.GetNetwork(); network != "" {
		bitmarkdService.SetNetwork(network)
		recorderdService.SetNetwork(network)
	}

	webserver := server.NewWebServer(
		nodeConfig,
		rootPath,
		bitmarkdService,
		recorderdService,
		masterConfig.VersionURL,
	)
	peerport := os.Getenv("PeerPort")
	if len(peerport) == 0 {
		peerport = "2136"
	}
	go webserver.CheckPortReachableRoutine(os.Getenv("PUBLIC_IP"), peerport)
	go webserver.ClearCmdErrorRoutine(bitmarkdService)

	r := gin.New()

	r.Use(static.Serve("/", static.LocalFile(uiPath, true)))
	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache")
	})
	apiRouter := r.Group("/api")
	apiRouter.GET("/info", webserver.NodeInfo)
	apiRouter.GET("/config", webserver.GetConfig)
	apiRouter.POST("/config", webserver.UpdateConfig)
	apiRouter.GET("/chain", webserver.GetChain)
	apiRouter.POST("/account/", webserver.NewAccount)
	apiRouter.GET("/account/", webserver.GetAccount)
	apiRouter.GET("/account/save", webserver.SaveAccount)
	apiRouter.GET("/account/delete", webserver.DeleteSavedAccount)
	apiRouter.POST("/account/phrase", webserver.SetRecoveryPhrase)
	apiRouter.GET("/account/phrase", webserver.GetRecoveryPhrase)
	apiRouter.GET("/bitmarkd/conn_stat", webserver.ConnectionStatus)
	apiRouter.POST("/bitmarkd", webserver.BitmarkdStartStop)
	apiRouter.GET("/latestVersion", webserver.LatestVersion)
	apiRouter.POST("/recorderd", webserver.RecorderdStartStop)
	apiRouter.GET("/log/:serviceName", webserver.GetLog)
	apiRouter.POST("/snapshot", webserver.DownloadSnapshot)
	apiRouter.GET("/snapshot-info", webserver.GetSnapshotInfo)

	r.Run(fmt.Sprintf(":%d", masterConfig.Port))
}

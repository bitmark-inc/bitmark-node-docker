package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitmark-inc/bitmark-node-docker/config"
	"github.com/bitmark-inc/bitmark-node-docker/server"
	"github.com/bitmark-inc/bitmark-node-docker/services"
	luaconf "github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/logger"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var version string = "v0.1" // do not change this value
var log *logger.L

// MasterConfiguration is a bitmark-node-docker Configuration file
type MasterConfiguration struct {
	Port       int                  `gluamapper:"port" json:"port"`
	DataDir    string               `gluamapper:"datadir" json:"datadir"`
	Logging    logger.Configuration `gluamapper:"logging" json:"logging"`
	VersionURL string               `gluamapper:"versionURL" json:"versionURL"`
}

// Parse take the filepath and parse the configure
func (c *MasterConfiguration) Parse(filepath string) error {
	if err := luaconf.ParseConfigurationFile(filepath, c); err != nil {
		panic(fmt.Sprintf("config file read failed: %s", err))
	}
	return nil
}

func init() {
	if len(os.Getenv("VERSION")) == 0 {
		log.Warn("Version is set to default")
		os.Setenv("VERSION", version)
	}
}

func main() {
	defer exitwithstatus.Handler()

	var confFile string
	var containerIP string
	var uiPath string
	flag.StringVar(&confFile, "config-file", "bitmark-node-docker.conf", "configuration for bitmark-node-docker")
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
	log = logger.New("bitmark-node-docker")
	log.Info(fmt.Sprintf("DataDirectory:%s", masterConfig.DataDir))
	log.Info(fmt.Sprintf("Port:%d", masterConfig.Port))
	log.Info(fmt.Sprintf("VersionURL:%s", masterConfig.VersionURL))

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

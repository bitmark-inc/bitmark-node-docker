package server

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bitmark-inc/bitmark-node/config"
	"github.com/bitmark-inc/bitmark-node/services"
	"github.com/bitmark-inc/logger"
	"github.com/gin-gonic/gin"
)

const (
	CONFIG_BUCKET_NAME = "config"
	SNAPSHOT_FILENAME  = "snapshot.zip"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

type ServiceOptionRequest struct {
	Option string `json:"option"`
}

type AccountInfo struct {
	network       string
	accountNumber string
	seed          string
}

type WebServer struct {
	Mutex             *sync.Mutex
	nodeConfig        *config.BitmarkNodeConfig
	rootPath          string
	log               *logger.L
	peerPortReachable bool
	Bitmarkd          services.Service
	Recorderd         services.Service
	Accounts          []AccountInfo
	versionURL        string
	SnapshotInfo      *snapshotInfo
}

type snapshotInfo struct {
	Date   string `json:"date"`
	Block  int    `json:"block"`
	URL    string `json:"url"`
	client *http.Client
}

func (s *snapshotInfo) get(versionURL string) ([]byte, error) {
	resp, err := s.client.Get(versionURL)

	if nil != err {
		return []byte{}, errors.New("Unable to get snapshot version info")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		return []byte{}, errors.New("Unable to read snapshot version file")
	}

	return body, nil
}

func NewWebServer(nc *config.BitmarkNodeConfig, rootPath string, bitmarkd, recorderd services.Service, versionURL string) *WebServer {
	return &WebServer{
		Mutex:      &sync.Mutex{},
		nodeConfig: nc,
		rootPath:   rootPath,
		log:        logger.New("webserver"),
		Bitmarkd:   bitmarkd,
		Recorderd:  recorderd,
		versionURL: versionURL,
		SnapshotInfo: &snapshotInfo{
			Date:   "",
			Block:  0,
			URL:    "",
			client: client,
		},
	}
}

func (ws *WebServer) GetChain(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()

	c.SetCookie("bitmark-node-network", network, 0, "", "", false, false)
	c.JSON(http.StatusOK, map[string]interface{}{
		"ok":     1,
		"result": network,
	})
	return
}

// func (ws *WebServer) SetChain(c *gin.Context) {
// 	reqBody := map[string]string{}
// 	err := c.BindJSON(&reqBody)
// 	if err != nil {
// 		c.String(400, "can not parse action option")
// 		return
// 	}

// 	network, ok := reqBody["network"]
// 	if !ok {
// 		c.String(400, "missing arguments")
// 		return
// 	}

// 	err = ws.nodeConfig.SetNetwork(network)
// 	if err != nil {
// 		c.String(400, "can not set network. error: %s", err.Error())
// 		return
// 	}

// 	ws.Bitmarkd.SetNetwork(network)
// 	ws.Recorderd.SetNetwork(network)

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"ok": 1,
// 	})
// 	return
// }

func (ws *WebServer) GetConfig(c *gin.Context) {
	config, err := ws.nodeConfig.Get()
	if err != nil {
		c.String(http.StatusInternalServerError, "can not read bitmark node config. error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"ok":     1,
		"result": config,
	})
	return
}

func (ws *WebServer) UpdateConfig(c *gin.Context) {
	newConfig := map[string]string{
		"btcAddr": "",
		"ltcAddr": "",
	}

	err := c.BindJSON(&newConfig)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = ws.nodeConfig.Set(newConfig)

	if err != nil {
		c.String(http.StatusInternalServerError, "can not set bitmark node config. error: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "")
	return
}

func (ws *WebServer) BitmarkdStartStop(c *gin.Context) {
	var req ServiceOptionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = nil
	switch req.Option {
	case "start":
		err = ws.Bitmarkd.Start()
	case "stop":
		err = ws.Bitmarkd.Stop()
	case "status":
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok":     1,
			"result": ws.Bitmarkd.Status(),
		})
		return
	case "info":

		resp, err := client.Get("https://127.0.0.1:2131/bitmarkd/details")
		if err != nil {
			c.String(http.StatusInternalServerError, "unable to get bitmark info")
			return
		}
		defer resp.Body.Close()
		bb := bytes.Buffer{}
		io.Copy(&bb, resp.Body)

		if resp.StatusCode != http.StatusOK {
			c.String(http.StatusInternalServerError, "unable to get bitmark info. message: %s", bb.String())
			return
		}

		var reply DetailReply
		d := json.NewDecoder(&bb)

		if err := d.Decode(&reply); err != nil {
			c.String(http.StatusInternalServerError, "fail to read bitmark info response. error: %s\n", err.Error())
			return
		}

		t, _ := time.ParseDuration(reply.Uptime)
		reply.Uptime = t.Round(time.Second).String()

		c.JSON(http.StatusOK, map[string]interface{}{
			"ok":     1,
			"result": reply,
		})
		return
	default:
		c.String(400, "invalid option")
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"ok":  0,
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": 1,
		})
	}
}

func (ws *WebServer) RecorderdStartStop(c *gin.Context) {
	var req ServiceOptionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.String(400, "can not parse action option")
		return
	}

	err = nil
	switch req.Option {
	case "start":
		err = ws.Recorderd.Start()
	case "stop":
		err = ws.Recorderd.Stop()
	case "status":
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok":     1,
			"result": ws.Recorderd.Status(),
		})
		return
	default:
		c.String(400, "invalid option")
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"ok":  0,
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": 1,
		})
	}
}

func (ws *WebServer) DiscoveryStartStop(c *gin.Context) {

}

func (ws *WebServer) GetLog(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	logFile := filepath.Join(ws.rootPath, c.Param("serviceName"), network, "log", c.Param("serviceName")+".log")

	file, err := os.Open(logFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lastLine := ""
	c.Header("X-Content-Type-Options", "nosniff")
	for {
		if !scanner.Scan() {
			fmt.Fprintln(c.Writer, lastLine)
			c.Writer.Flush()
			break
		}
		lastLine = scanner.Text()
	}

	reader := bufio.NewReader(file)
	c.Stream(func(w io.Writer) bool {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(1 * time.Second)
			} else {
				fmt.Fprintf(w, "===== log stopped with error: %s", err.Error())
				return false
			}
		}
		fmt.Fprint(w, line)
		return true
	})
}

func (ws *WebServer) GetAccountNumber(network string) (string, error) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network && item.accountNumber != "" {
			return item.accountNumber, nil
		}
	}

	return "", errors.New("No account in " + network + " network")

}

func (ws *WebServer) GetSeed(network string) (string, error) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network && item.seed != "" {
			return item.seed, nil
		}
	}
	return "", errors.New("No seed of AccountInfo in " + network + " network")
}

func (ws *WebServer) SetAccount(accountNumber, seed, network string) error {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	for _, item := range ws.Accounts {
		if item.network == network {
			item.accountNumber = accountNumber
			item.seed = seed
			return nil
		}
	}
	ws.Accounts = append(ws.Accounts, AccountInfo{accountNumber: accountNumber, seed: seed, network: network})
	ws.log.Infof("[SetAccount]Append account Item:")
	return nil
}

// get snapshot info
func getSnapshotInfo(versionURL string, s *snapshotInfo) error {
	// no need to get snapshot info if block count is larger than 0
	if s.Block != 0 {
		return nil
	}

	body, err := s.get(versionURL)

	if nil != err {
		return err
	}

	err = json.Unmarshal(body, s)

	if nil != err {
		return errors.New("Unable to decode snapshot version file")
	}

	return nil
}

// retrieve snapshot info
func (ws *WebServer) GetSnapshotInfo(c *gin.Context) {
	const fn = "[GetSnapshotInfo]"
	err := getSnapshotInfo(ws.versionURL, ws.SnapshotInfo)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		ws.log.Errorf("%s: %s", fn, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"info": ws.SnapshotInfo,
	})
	return
}

// download snapshot zip file
func downloadFile(ws *WebServer) error {
	resp, err := client.Get(ws.SnapshotInfo.URL)

	if nil != err {
		return fmt.Errorf("Cannot download snapshot from %s", ws.SnapshotInfo.URL)
	}
	defer resp.Body.Close()

	bitmarkdDataPath := filepath.Join(ws.Bitmarkd.GetPath(), ws.Bitmarkd.GetNetwork())
	// create directory if not exist
	err = os.MkdirAll(bitmarkdDataPath, 0755)
	if nil != err {
		return err
	}

	filePath := filepath.Join(bitmarkdDataPath, SNAPSHOT_FILENAME)

	// create file
	fo, err := os.Create(filePath)
	if nil != err {
		return err
	}
	defer fo.Close()

	io.Copy(fo, resp.Body)
	if nil != err {
		return err
	}

	// check if file successfully saved
	fileInfo, err := os.Stat(filePath)
	if nil != err {
		return err
	}

	if fileInfo.Size() == 0 {
		return errors.New("Cannot save file")
	}

	return nil
}

type processStatus struct {
}

type bitmarkdStatus struct {
	started bool
	running bool
	error   bool
}

// return bitmarkd running status
func isBitmarkdStop(ws *WebServer) bool {
	stat := ws.Bitmarkd.Status()
	return !stat["started"].(bool)
}

// unzip file
func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// get filenames in zip
		fpath := filepath.Join(dest, f.Name)

		// check zip content
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// create foled
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// create file
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			outFile, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// close file
			outFile.Close()

			if err != nil {
				return filenames, err
			}
		}
	}
	return filenames, nil
}

// recover data form zip file
func recoverData(ws *WebServer) error {
	bitmarkdDataPath := filepath.Join(
		ws.Bitmarkd.GetPath(),
		ws.Bitmarkd.GetNetwork(),
	)

	oldDir := filepath.Join(bitmarkdDataPath, "data")

	// remove existing backup directory
	os.RemoveAll(filepath.Join(bitmarkdDataPath, "data-backup"))

	// backup directory if exist
	if _, err := os.Stat(oldDir); !os.IsNotExist(err) {
		err := os.Rename(
			oldDir,
			filepath.Join(bitmarkdDataPath, "data-backup"))

		if nil != err {
			return fmt.Errorf("Cannot move existing directory %s", oldDir)
		}
	}

	// extract file
	files, err := unzip(
		filepath.Join(bitmarkdDataPath, SNAPSHOT_FILENAME),
		bitmarkdDataPath,
	)

	if nil != err {
		return err
	}

	fmt.Println(files)
	return nil
}

// download snapshot
func (ws *WebServer) DownloadSnapshot(c *gin.Context) {
	const fn = "[DownloadSnapshot]"
	err := getSnapshotInfo(ws.versionURL, ws.SnapshotInfo)
	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		ws.log.Errorf("%s: %s", fn, err)
		return
	}

	// download file
	err = downloadFile(ws)
	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		ws.log.Errorf("%s: %s", fn, err)
		return
	}

	// make sure bitmarkd is not running
	stopped := isBitmarkdStop(ws)
	if !stopped {
		err = ws.Bitmarkd.Stop()
		if nil != err {
			c.String(http.StatusInternalServerError, "Cannot stop bitmarkd")
			ws.log.Errorf("%s: %s", fn, err)
			return
		}
	}

	// overwrite file
	err = recoverData(ws)

	// show response
	if nil != err {
		c.String(http.StatusInternalServerError, err.Error())
		ws.log.Errorf("%s: %s", fn, err)
		return
	}

	c.String(http.StatusOK, "")
	return
}

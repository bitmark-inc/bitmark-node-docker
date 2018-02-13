package server

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/gin-gonic/gin"
)

type RecoveryPhraseArguments struct {
	Phrases string
}

func (ws *WebServer) GetRecoveryPhrase(c *gin.Context) {
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

	a, err := sdk.AccountFromSeed(seed)
	if err != nil {
		c.String(500, "get account from seed")
		return
	}

	phrases := a.RecoveryPhrase()

	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": strings.Join(phrases, " "),
	})
}

func (ws *WebServer) SetRecoveryPhrase(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		c.String(500, "wrong network configuration")
		return
	}

	var args RecoveryPhraseArguments
	err := c.BindJSON(&args)
	if err != nil {
		c.String(400, "fail to parse arguments")
		return
	}

	a, err := sdk.AccountFromRecoveryPhrase(args.Phrases)
	if err != nil {
		c.String(400, "fail to recover an account")
		return
	}

	seed := a.Seed()

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		c.String(500, "fail to open seed file from: %s", seedFile)
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		c.String(500, "fail to update seed file")
		return
	}

	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
}

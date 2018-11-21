package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
)

const ACCT_BUCKET_NAME = "account"
const SEED_KEY_NAME = "seed"

type RecoveryPhraseArguments struct {
	Phrase string `json:"phrase"`
}

func returnError(c *gin.Context, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"ok":      0,
		"message": message,
	})
}

//Get seed string from a bitmarkd seed file
func GetSeedFromFile(seedFile string) (string, error) {
	f, err := os.Open(seedFile)
	if err != nil {
		return "", fmt.Errorf("fail to open seed file. error: %s", err.Error())
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return "", fmt.Errorf("fail not read seed file. error: %s", err.Error())
	}
	seed := strings.Trim(strings.Split(buf.String(), ":")[1], "\n")
	return seed, nil
}

//GetSeedFromDB is to get Seed from DB
func (ws *WebServer) GetSeedFromDB(dbPath, network string) (seed string, err error) {
	if network == "" {
		return "", errors.New("no network")
	}

	db, err := bolt.Open(filepath.Join(dbPath, "bitmark-node.db"), 0600, nil)
	defer db.Close()
	if err != nil {
		return "", err
	}

	err = db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte(ACCT_BUCKET_NAME))
		bucket := rootBucket.Bucket([]byte(network))
		seedbytes := bucket.Get([]byte(SEED_KEY_NAME))
		seed = string(seedbytes[:])
		return nil

	})

	if err != nil {
		return "", err
	}
	return
}

//SaveSeedToDB is to get Seed from DB
func (ws *WebServer) SaveSeedToDB(seed, dbPath, network string) error {
	db, err := bolt.Open(filepath.Join(dbPath, "bitmark-node.db"), 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	if network == "" {
		return errors.New("no network")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ACCT_BUCKET_NAME))
		if err != nil {
			return err
		}
		//set current network
		bkt := tx.Bucket([]byte(ACCT_BUCKET_NAME))
		err = bkt.Put([]byte("network"), []byte(network))
		if err != nil {
			return err
		}

		//if subBucket does not exist create it
		subBucket, err := bkt.CreateBucketIfNotExists([]byte(network))
		if err != nil {
			return err
		}
		//write seed in db
		err = subBucket.Put([]byte("seed"), []byte(seed))
		if err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		return err
	}
	return nil
}

//Get the current account which is set in the bitmarkd proofing file
func (ws *WebServer) GetAccount(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ws.log.Errorf("[GetAccount]", "wrong network configuration")
		returnError(c, 500, "wrong network configuration")
		return
	}
	// Return AccountNumber if there is a record in memory
	number, err := ws.GetAccountNumber(network)
	if err == nil { // If there is a record in AccountInfo, return it
		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": number,
		})
	}
	// Return AccountNumber if there is a seedFile
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	seed, err := GetSeedFromFile(seedFile)
	if err != nil {
		ws.log.Errorf("[GetAccount]", "can not get seed from file:", err.Error())
		returnError(c, 404, fmt.Sprintf("can not get seed from file. reason: %s", err.Error()))
		return
	}

	a, err := sdk.AccountFromSeed(seed)
	if err != nil {
		ws.log.Errorf("[GetAccount]", "can not get account from seed:", err.Error())
		returnError(c, 500, fmt.Sprintf("can not get account from seed. reason: %s", err.Error()))
		return
	}

	ws.SetAccount(a.AccountNumber(), seed, network)
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": a.AccountNumber(),
	})

	dbPath := filepath.Join(ws.rootPath, "db")
	seed, err = ws.GetSeedFromDB(dbPath, network)
	if err != nil {
		returnError(c, 500, fmt.Sprintf("can not get account from saved account. reason: %s", err.Error()))
	}
	ws.SetAccount(a.AccountNumber(), seed, network)
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": a.AccountNumber(),
	})

}

func (ws *WebServer) SaveAccount(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ws.log.Errorf("[SaveAccount]", "wrong network configuration")
		returnError(c, 500, "wrong network configuration")
		return
	}

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	seed, err := GetSeedFromFile(seedFile)
	if err != nil {
		ws.log.Errorf("[GetAccount]", "can not get seed from file:", err.Error())
		returnError(c, 404, fmt.Sprintf("can not get seed from file. reason: %s", err.Error()))
		return
	}
	dbPath := filepath.Join(ws.rootPath, "db")
	ws.SaveSeedToDB(seed, dbPath, network)
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})

}

func (ws *WebServer) NewAccount(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ws.log.Errorf("[NewAccount]", "wrong network configuration")
		returnError(c, 500, "wrong network configuration")
		return
	}
	n := sdk.Testnet
	if network == "bitmark" {
		n = sdk.Livenet
	}
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	if _, err := os.Stat(seedFile); err == nil {
		ws.log.Errorf("[NewAccount]", err)
		returnError(c, 500, fmt.Sprintf("seed file is existed: %s", seedFile))
		return
	}

	a, err := sdk.NewAccount(n)
	if err != nil {
		ws.log.Errorf("[NewAccount]", "fail to create a new account:", err)
		returnError(c, 400, "fail to create a new account")
		return
	}
	seed := a.Seed()

	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		ws.log.Errorf("[NewAccount]", "fail to open seed file", err)
		returnError(c, 500, fmt.Sprintf("fail to open seed file from: %s", seedFile))
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		ws.log.Errorf("[NewAccount]", "fail to update seed file")
		returnError(c, 500, "fail to update seed file")
		return
	}
	ws.SetAccount(a.AccountNumber(), seed, network) // Record in AccountInfo in memory
	ws.log.Infof("Create a NewAccount", a.AccountNumber())
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
}

func (ws *WebServer) GetRecoveryPhrase(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ws.log.Errorf("[GetRecoveryPhrase]", "wrong network configuration")
		returnError(c, 500, "wrong network configuration")
		return
	}

	seed, err := ws.GetSeed(network)
	if err != nil { //read from file
		seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
		seed, err = GetSeedFromFile(seedFile)
		if err != nil {
			ws.log.Errorf("[GetRecoveryPhrase]", err.Error())
			returnError(c, 500, fmt.Sprintf("can not get seed from file. reason: %s", err.Error()))
			return
		}
	}

	a, err := sdk.AccountFromSeed(seed)

	if err != nil {
		ws.log.Errorf("[GetRecoveryPhrase]", err)
		returnError(c, 500, "get account from seed")
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
		ws.log.Errorf("[SetRecoveryPhrase]", "wrong network configuration")
		returnError(c, 500, "wrong network configuration")
		return
	}

	var args RecoveryPhraseArguments
	err := c.BindJSON(&args)
	if err != nil {
		ws.log.Errorf("[SetRecoveryPhrase]", err)
		returnError(c, 400, "invalid request arguments")
		return
	}

	a, err := sdk.AccountFromRecoveryPhrase(args.Phrase)
	if err != nil {
		ws.log.Errorf("[SetRecoveryPhrase]", err)
		returnError(c, 400, "fail to recover an account from the phrase")
		return
	}

	seed := a.Seed()

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		returnError(c, 500, fmt.Sprintf("fail to open seed file from: %s", seedFile))
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		returnError(c, 500, "fail to update seed file")
		return
	}
	ws.SetAccount(a.AccountNumber(), seed, network)
	ws.log.Infof("Set Account by Recovery Phrase")
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
}

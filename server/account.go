package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitmark-inc/bitmark-node/config"
	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	bolt "github.com/coreos/bbolt"
	"github.com/gin-gonic/gin"
)

const SEED_KEY_NAME = "seed"

type RecoveryPhraseArguments struct {
	Phrase string `json:"phrase"`
}

// GetSeedFromFile Get seed string from a bitmarkd seed file
func GetSeedFromFile(seedFile string) (string, error) {
	f, err := os.Open(seedFile)
	if err != nil {
		return "", ErrCombind(ErrorOpenSeedFile, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return "", ErrCombind(ErrorToReadSeedFile, err)
	}
	seed := strings.Trim(strings.Split(buf.String(), ":")[1], "\n")
	return seed, nil
}

// GetSeedFromDB is to get Seed from DB
func (ws *WebServer) GetSeedFromDB(network string) (seed string, err error) {
	if network == "" {
		return "", ErrorNoNetwork
	}

	db := ws.nodeConfig.GetDB()

	if db == nil {
		return "", ErrorGetBoltDB
	}

	err = db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte(config.CONFIG_BUCKET_NAME))
		bucket := rootBucket.Bucket([]byte(network))
		seedbytes := bucket.Get([]byte(SEED_KEY_NAME))
		seed = string(seedbytes[:])
		return nil

	})

	if err != nil {
		return "", ErrCombind(ErrorUpdateBoltDB, err)
	}
	return seed, nil
}

// SaveSeedToDB is to get Seed from DB
func (ws *WebServer) SaveSeedToDB(seed, dbPath, network string) error {
	if network == "" {
		ws.log.Errorf("SaveSeedToDB:%s", ErrorNoNetwork)
		return ErrorNoNetwork
	}

	db := ws.nodeConfig.GetDB()
	if nil == db {
		ws.log.Errorf("SaveSeedToDB:%s", ErrorGetBoltDB)
		return ErrorGetBoltDB
	}

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(config.CONFIG_BUCKET_NAME))
		if err != nil {
			return ErrCombind(ErrorCreateBoltDB, err)
		}
		//set current network
		bkt := tx.Bucket([]byte(config.CONFIG_BUCKET_NAME))
		err = bkt.Put([]byte("network"), []byte(network))
		if err != nil {
			return ErrCombind(ErrorBoltDBCreateButcket, err)
		}

		//if subBucket does not exist create it
		subBucket, err := bkt.CreateBucketIfNotExists([]byte(network))
		if err != nil {
			return ErrCombind(ErrorBoltDBCreateSubButcket, err)
		}
		//write seed in db
		err = subBucket.Put([]byte("seed"), []byte(seed))
		if err != nil {
			return ErrCombind(ErrorSaveSeedToDB, err)
		}
		return nil
	})

	if err != nil {
		ws.log.Errorf("SaveSeedToDB:%s", ErrCombind(ErrorSaveSeedToDB, err))
		return ErrCombind(ErrorSaveSeedToDB, err)
	}
	return nil
}

// LoadSavedAcct load saved account to memory and seed file
func (ws *WebServer) LoadSavedAcct(dbPath, network string) (string, error) {
	seed, err := ws.GetSeedFromDB(network)
	if err != nil {
		ws.log.Errorf("LoadSavedAcct:%s", ErrCombind(ErrorBoltDBGetSeed, err))
		return "", ErrCombind(ErrorBoltDBGetSeed, err)
	}
	if seed == "" {
		ws.log.Errorf("LoadSavedAcct:%s", ErrorGetEmptySeed)
		return "", ErrorGetEmptySeed
	}
	a, err := sdk.AccountFromSeed(seed)
	if err != nil {
		ws.log.Errorf("LoadSavedAcct:%s", ErrCombind(ErrorGetAccountFromSeed, err))
		return "", ErrorGetAccountFromSeed
	}
	err = ws.SetAccount(a.AccountNumber(), seed, network)

	if err != nil {
		ws.log.Errorf("LoadSavedAcct:%s", ErrorSetAccount)
		return "", ErrorSetAccount
	}

	//Save to file and load to memory
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		ws.log.Errorf("LoadSavedAcct:%s", ErrCombind(ErrorOpenSeedFile, err))
		return "", ErrCombind(ErrorOpenSeedFile, err)
	}
	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))

	if err != nil {
		return "", err
	}

	return a.AccountNumber(), nil
}

// GetAccount Get the current account which is set in the bitmarkd proofing file		ws.log.Errorf("GetAccount:%s", ErrCombind(ErrorGetAccountFromSeed, err))
func (ws *WebServer) GetAccount(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ReturnError(c, 500, ErrorNoNetwork.Error())
		return
	}
	// get from saved account
	dbPath := filepath.Join(ws.rootPath, "db")
	_, err := ws.LoadSavedAcct(dbPath, network)

	// Return AccountNumber if there is a record in memory
	number, err := ws.GetAccountNumber(network)
	if err == nil { // If there is a record in AccountInfo, return it
		c.JSON(200, map[string]interface{}{
			"ok":     1,
			"result": number,
		})
		return
	}

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	seed, err := GetSeedFromFile(seedFile)

	if err != nil {
		ReturnError(c, 500, ErrorGetSeedFromFile.Error())
		return
	}

	a, err := sdk.AccountFromSeed(seed)
	if err != nil {
		ReturnError(c, 500, ErrCombind(ErrorGetAccountFromSeed, err).Error())
		return
	}

	ws.SetAccount(a.AccountNumber(), seed, network)
	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": a.AccountNumber(),
	})
}

// NewAccount create a new acccount
func (ws *WebServer) NewAccount(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ReturnError(c, 500, ErrorNoNetwork.Error())
		return
	}
	n := sdk.Testnet
	if network == "bitmark" {
		n = sdk.Livenet
	}
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	if _, err := os.Stat(seedFile); err == nil {
		ReturnError(c, 500, ErrorNoSeedFile.Error())
		return
	}

	a, err := sdk.NewAccount(n)
	if err != nil {
		ReturnError(c, 500, ErrorCreateAccount.Error())
		return
	}
	seed := a.Seed()

	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		ReturnError(c, 500, ErrorOpenSeedFile.Error())
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		ReturnError(c, 500, ErrorToWriteSeedFile.Error())
		return
	}
	ws.SetAccount(a.AccountNumber(), seed, network) // Record in AccountInfo in memory
	err = ws.saveAcct()
	if nil != err {
		ReturnError(c, 500, ErrorAutoSaveAccount.Error())
		return

	}
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
}

// GetRecoveryPhrase get recoveryPhrase of current account
func (ws *WebServer) GetRecoveryPhrase(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ReturnError(c, 500, ErrorNoNetwork.Error())
		return
	}

	seed, err := ws.GetSeed(network)
	if err != nil { //read from file
		seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
		seed, err = GetSeedFromFile(seedFile)
		if err != nil {
			ReturnError(c, 500, ErrCombind(ErrorGetSeedFromFile, err).Error())
			return
		}
	}

	a, err := sdk.AccountFromSeed(seed)

	if err != nil {
		ReturnError(c, 500, ErrorGetAccountFromSeed.Error())
		return
	}

	phrases := a.RecoveryPhrase()

	c.JSON(200, map[string]interface{}{
		"ok":     1,
		"result": strings.Join(phrases, " "),
	})
}

// SetRecoveryPhrase set recovery phrase
func (ws *WebServer) SetRecoveryPhrase(c *gin.Context) {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ReturnError(c, 500, ErrorNoNetwork.Error())
		return
	}
	var args RecoveryPhraseArguments
	err := c.BindJSON(&args)
	if err != nil {
		ReturnError(c, 400, ErrorInvalidArgument.Error())
		return
	}

	a, err := sdk.AccountFromRecoveryPhrase(args.Phrase)
	if err != nil {
		ReturnError(c, 400, ErrorRecoveryFromPhrase.Error())
		return
	}

	seed := a.Seed()

	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		ReturnError(c, 500, ErrorOpenSeedFile.Error())
		return
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		ReturnError(c, 500, ErrorToWriteSeedFile.Error())
		return
	}
	ws.SetAccount(a.AccountNumber(), seed, network)
	err = ws.saveAcct()
	if nil != err {
		ReturnError(c, 500, ErrorAccountCreateNotSaved.Error())
		return
	}
	ws.log.Infof("Set Account by Recovery Phrase")
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
}

// SaveAccount save current account to db and file

func (ws *WebServer) SaveAccount(c *gin.Context) {

	if err := ws.saveAcct(); nil != err {
		ReturnError(c, 500, ErrorSaveAccount.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
	return
}

func (ws *WebServer) saveAcct() error {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		return errors.New("wrong network configuration")
	}
	seed, err := ws.GetSeed(network)
	if err != nil {
		return errors.New("fail to get seed from webserver")
	}
	dbPath := filepath.Join(ws.rootPath, "db")
	err = ws.SaveSeedToDB(seed, dbPath, network)
	if err != nil {
		return errors.New("save to db failed")
	}
	//also save to file
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()

	if err != nil {
		return errors.New("fail to save seedfile")
	}

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))

	//verify
	_, err = ws.GetSeedFromDB(network)
	if err != nil {
		return errors.New("get Seed FromDB failed")
	}
	return nil
}

// DeleteSavedAccount Delete saved account in File and Database
func (ws *WebServer) DeleteSavedAccount(c *gin.Context) {
	if err := ws.deleteSavedAcct(); err != nil {
		ReturnError(c, 500, ErrorDeleteSavedAccount.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": 1,
	})
	return
}

func (ws *WebServer) deleteSavedAcct() error {
	network := ws.nodeConfig.GetNetwork()
	if network == "" {
		ws.log.Warnf("deleteSavedAcct:%s", ErrorNoNetwork)
		return ErrorNoNetwork
	}
	// delete database according to network
	db := ws.nodeConfig.GetDB()
	if nil == db {
		ws.log.Warnf("deleteSavedAcct:%s", ErrorGetBoltDB)
		return ErrorGetBoltDB
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte(config.CONFIG_BUCKET_NAME))
		return rootBucket.Bucket([]byte(network)).Delete([]byte("network"))
	}); err != nil {
		return ErrorUpdateBoltDB
	}

	// delete seed file
	seedFile := filepath.Join(ws.rootPath, "bitmarkd", network, "proof.sign")
	err := os.Remove(seedFile)
	if err != nil {
		return ErrorDeleteSeedFile
	}
	return nil
}

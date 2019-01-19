package server

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var ( // Account
	// Seed File Errors
	ErrorOpenSeedFile    = errors.New("Open seed file failed")
	ErrorGetSeedFromFile = errors.New("Get seed from file failed")
	ErrorNoSeedFile      = errors.New("No seed file")
	ErrorToReadSeedFile  = errors.New("Read Seed file failed")
	ErrorToWriteSeedFile = errors.New("Write Seed file failed")
	ErrorDeleteSeedFile  = errors.New("Delete Seed file failed")
	// Network Errors
	ErrorNoNetwork = errors.New("no network")
	// BoltDB Errors
	ErrorGetBoltDB              = errors.New("Get bolt database failed")
	ErrorUpdateBoltDB           = errors.New("Update bolt database failed")
	ErrorCreateBoltDB           = errors.New("Create bolt database failed")
	ErrorBoltDBCreateButcket    = errors.New("Create bolt db bucket failed")
	ErrorBoltDBCreateSubButcket = errors.New("Create bolt db sub-bucket failed")
	ErrorBoltDBWriteSeed        = errors.New("Write seed to bolt db failed")
	ErrorBoltDBGetSeed          = errors.New("Get seed from bolt db failed")
	ErrorSaveSeedToDB           = errors.New("Save seed to db failed")
	ErrorGetEmptySeed           = errors.New("Empty Seed")
	ErrorGetAccountFromSeed     = errors.New("Get Account from seed fail")
	// Functions Error
	ErrorCreateAccount         = errors.New("Create Account fail")
	ErrorSetAccount            = errors.New("Set account fail")
	ErrorGetAccount            = errors.New("Get Account fail")
	ErrorAutoSaveAccount       = errors.New("Auto save Account fail")
	ErrorRecoveryFromPhrase    = errors.New("Recovery From Phrase fail")
	ErrorAccountCreateNotSaved = errors.New("Account has created but not saved")
	ErrorSaveAccount           = errors.New("Save account fail")
	ErrorGetSeed               = errors.New("Get seed fail")
	ErrorDeleteSavedAccount    = errors.New("Delete saved account fail")
	ErrorReadConfigFile        = errors.New("Read Config file fail")
	ErrorSetConfigFile         = errors.New("Set Config file fail")
	// Webserver Errors
	ErrorInvalidArgument = errors.New("Invalid arguments")
	ErrorInvalidOption   = errors.New("Invalid option")
	ErrorParseOption     = errors.New("Parsse option fail")

	// Snapshot
	ErrorGetSnapshotInfo  = errors.New("Get snapshot info fail")
	ErrorReadSnapshot     = errors.New("Read snapshot from version fail")
	ErrorDownloadSnapshot = errors.New("Download snapshot fail")
	ErrorSaveSnapshot     = errors.New("Save snapshot fail")
	ErrorMoveDirFail      = errors.New("Move Directory fail")

	// Service
	ErrorGetBitmarkdInfo = errors.New("Get bitmarkd info fail")
)

// ErrCombind return an error which combined two errors using - seperator
func ErrCombind(cause, detail error) error {
	return fmt.Errorf("%s-%s", cause, detail)
}

// ReturnError return error response for gin server
func ReturnError(c *gin.Context, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"ok":      0,
		"message": message,
	})
}

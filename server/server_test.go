package server

import (
	"os"
	"testing"

	"github.com/bitmark-inc/logger"
	"github.com/gin-gonic/gin"
)

// fake response of bitmarkd is running
type FakeRunningBitmarkd struct{}

func (f *FakeRunningBitmarkd) Status() map[string]interface{} {
	return map[string]interface{}{
		"started": true,
	}
}

func (f *FakeRunningBitmarkd) Initialise(str string) error { return nil }
func (f *FakeRunningBitmarkd) Finalise() error             { return nil }
func (f *FakeRunningBitmarkd) IsRunning() bool             { return true }
func (f *FakeRunningBitmarkd) SetNetwork(string)           {}
func (f *FakeRunningBitmarkd) Start() error                { return nil }
func (f *FakeRunningBitmarkd) Stop() error                 { return nil }
func (f *FakeRunningBitmarkd) GetPath() string             { return "" }
func (f *FakeRunningBitmarkd) GetNetwork() string          { return "" }

// fake response of bitmarkd is not running
type FakeStoppedBitmarkd struct{}

func (f *FakeStoppedBitmarkd) Status() map[string]interface{} {
	return map[string]interface{}{
		"started": false,
	}
}

func (f *FakeStoppedBitmarkd) Initialise(str string) error { return nil }
func (f *FakeStoppedBitmarkd) Finalise() error             { return nil }
func (f *FakeStoppedBitmarkd) IsRunning() bool             { return true }
func (f *FakeStoppedBitmarkd) SetNetwork(string)           {}
func (f *FakeStoppedBitmarkd) Start() error                { return nil }
func (f *FakeStoppedBitmarkd) Stop() error                 { return nil }
func (f *FakeStoppedBitmarkd) GetPath() string             { return "" }
func (f *FakeStoppedBitmarkd) GetNetwork() string          { return "" }

const (
	DIR = "testing"
)

func removeTestFiles() {
	os.RemoveAll(DIR)
}

// setup for logger
func setup() {
	// create log directory
	os.Mkdir(DIR, 0700)

	// logger config
	conf := logger.Configuration{
		Directory: DIR,
		File:      "testing.log",
		Size:      1048576,
		Count:     10,
		Levels: map[string]string{
			logger.DefaultTag: "critical",
		},
		Console: false,
	}
	if err := logger.Initialise(conf); nil != err {
		panic("logger setup failed: " + err.Error())
	}

	gin.SetMode(gin.TestMode)
}

// isBitmarkdStop
func Test_isBitmarkdStop(t *testing.T) {
	setup()
	defer removeTestFiles()

	bitmarkdTrue := &FakeRunningBitmarkd{}
	recorderd := &FakeRunningBitmarkd{}

	ws := NewWebServer(nil, "", bitmarkdTrue, recorderd, "")

	if isBitmarkdStop(ws) {
		t.Error("isBitmarkdStop returns true when bitmarkd is not running")
	}

	bitmarkdFalse := &FakeStoppedBitmarkd{}

	ws = NewWebServer(nil, "", bitmarkdFalse, recorderd, "")

	if !isBitmarkdStop(ws) {
		t.Error("isBitmarkdStop returns false when bitmarkd is running")
	}
}

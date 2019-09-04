// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/bitmark-inc/bitmark-node-docker/fault"
	"github.com/bitmark-inc/bitmark-node-docker/utils"
	"github.com/bitmark-inc/logger"
)

var (
	ErrRecorderdIsNotRunning = InvalidError("Recorderd is not running")
	ErrRecorderdIsRunning    = InvalidError("Recorderd is running")
)

type Recorderd struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
	rootPath    string
	configFile  string
	network     string
	process     *os.Process
	running     bool
	ModeStart   chan bool
}

func NewRecorderd() *Recorderd {
	return &Recorderd{}
}

func (recorderd *Recorderd) GetPath() string {
	return recorderd.rootPath
}

func (recorderd *Recorderd) GetNetwork() string {
	if len(recorderd.network) == 0 {
		return "bitmark"
	}
	return recorderd.network
}

func (recorderd *Recorderd) Initialise(rootPath string) error {
	recorderd.Lock()
	defer recorderd.Unlock()

	if recorderd.initialised {
		return fault.ErrAlreadyInitialised
	}

	recorderd.rootPath = rootPath

	recorderd.log = logger.New("service-recorderd")

	recorderd.running = false
	recorderd.ModeStart = make(chan bool, 1)

	// all data initialised
	recorderd.initialised = true
	return nil
}

func (recorderd *Recorderd) Finalise() error {
	recorderd.Lock()
	defer recorderd.Unlock()

	if !recorderd.initialised {
		return fault.ErrNotInitialised
	}

	recorderd.initialised = false
	return nil
}

func (recorderd *Recorderd) IsRunning() bool {
	return recorderd.running
}

func (recorderd *Recorderd) Status() map[string]interface{} {
	return map[string]interface{}{
		"started": recorderd.running,
		"running": recorderd.running,
	}
}

func (recorderd *Recorderd) Run(args interface{}, shutdown <-chan struct{}) {
loop:
	for {
		select {

		case <-shutdown:
			break loop
		case start := <-recorderd.ModeStart:
			if start {
				recorderd.Start()
			} else {
				recorderd.Stop()
			}
		}

	}
	close(recorderd.ModeStart)
}

func (recorderd *Recorderd) SetNetwork(network string) {
	if recorderd.running {
		recorderd.Stop()
	}
	recorderd.network = network
	switch network {
	case "testing":
		recorderd.configFile = filepath.Join(recorderd.rootPath, "testing/recorderd.conf")
	case "bitmark":
		fallthrough
	default:
		recorderd.configFile = filepath.Join(recorderd.rootPath, "bitmark/recorderd.conf")
	}
}

func (recorderd *Recorderd) Start() error {
	if recorderd.running {
		recorderd.log.Errorf("Start recorderd failed: %v", ErrRecorderdIsRunning)
		return ErrRecorderdIsRunning
	}

	// Check recorderdConfigFile exists
	recorderd.log.Infof("recorderd config file: %s\n", recorderd.configFile)
	if !utils.EnsureFileExists(recorderd.configFile) {
		recorderd.log.Errorf("Start recorderd failed: %v", fault.ErrNotFoundConfigFile)
		return fault.ErrNotFoundConfigFile
	}

	recorderd.running = true
	stopped := make(chan bool, 1)

	bitmarkDPublicIP := os.Getenv("PUBLIC_IP")
	proofPubPortEnv := os.Getenv("PROOF_PUB_PORT")
	proofSubPortEnv := os.Getenv("PROOF_SUB_PORT")

	go func() {
		for recorderd.running {
			// start recorderd as sub process
			cmd := exec.Command("recorderd", "--config-file="+recorderd.configFile)

			cmd.Env = []string{
				fmt.Sprintf("PUBLIC_IP=%s", bitmarkDPublicIP),
				fmt.Sprintf("PROOF_PUB_PORT=%s", proofPubPortEnv),
				fmt.Sprintf("PROOF_SUB_PORT=%s", proofSubPortEnv),
			}

			// start recorderd as sub process
			stderr, err := cmd.StderrPipe()
			if err != nil {
				recorderd.log.Errorf("Error: %v", err)
				continue
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				recorderd.log.Errorf("Error: %v", err)
				continue
			}
			if err := cmd.Start(); nil != err {
				continue
			}

			recorderd.process = cmd.Process
			recorderd.log.Infof("process id: %d", cmd.Process.Pid)
			stdeReader := bufio.NewReader(stderr)
			stdoReader := bufio.NewReader(stdout)

			go func() {
				for {
					stde, err := stdeReader.ReadString('\n')
					recorderd.log.Errorf("recorderd stderr: %q", stde)
					if nil != err {
						recorderd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			go func() {
				for {
					stdo, err := stdoReader.ReadString('\n')
					recorderd.log.Infof("recorderd stdout: %q", stdo)
					if nil != err {
						recorderd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			if err := cmd.Wait(); nil != err {
				if recorderd.running {
					recorderd.log.Errorf("recorderd has terminated unexpectedly. failed: %v", err)
					recorderd.log.Errorf("recorderd will be restarted in 1 second...")
					time.Sleep(time.Second)
				}
				recorderd.process = nil
				stopped <- true
			}
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	return nil

}

func (recorderd *Recorderd) Stop() error {
	if !recorderd.running {
		recorderd.log.Errorf("Stop recorderd failed: %v", ErrRecorderdIsNotRunning)
		return ErrRecorderdIsNotRunning
	}
	recorderd.running = false

	if err := recorderd.process.Signal(os.Kill); nil != err {
		recorderd.log.Errorf("Send kill to recorderd failed: %v", err)
		return err
	}

	recorderd.log.Infof("Stop recorderd. PID: %d", recorderd.process.Pid)
	recorderd.process = nil
	return nil
}

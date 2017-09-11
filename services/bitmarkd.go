// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bufio"

	"github.com/bitmark-inc/bitmark-node/fault"
	"github.com/bitmark-inc/bitmark-node/utils"

	"github.com/bitmark-inc/logger"

	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	ErrBitmarkdIsNotRunning = fault.InvalidError("Bitmarkd is not running")
	ErrBitmarkdIsRunning    = fault.InvalidError("Bitmarkd is running")
)

type Bitmarkd struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
	configFile  string
	process     *os.Process
	running     bool
	ModeStart   chan bool
}

func NewBitmarkd() *Bitmarkd {
	return &Bitmarkd{}
}

func (bitmarkd *Bitmarkd) Initialise(configFile string) error {
	bitmarkd.Lock()
	defer bitmarkd.Unlock()

	if bitmarkd.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitmarkd.configFile = configFile

	bitmarkd.log = logger.New("service-bitmarkd")

	bitmarkd.running = false
	bitmarkd.ModeStart = make(chan bool, 1)

	// all data initialised
	bitmarkd.initialised = true
	return nil
}

func (bitmarkd *Bitmarkd) Finalise() error {
	bitmarkd.Lock()
	defer bitmarkd.Unlock()

	if !bitmarkd.initialised {
		return fault.ErrNotInitialised
	}

	bitmarkd.initialised = false
	bitmarkd.Stop()
	return nil
}

func (bitmarkd *Bitmarkd) IsRunning() bool {
	return bitmarkd.running
}

func (bitmarkd *Bitmarkd) Status() string {
	if bitmarkd.running {
		return "started"
	} else {
		return "stopped"
	}
}

func (bitmarkd *Bitmarkd) Run(args interface{}, shutdown <-chan struct{}) {
loop:
	for {
		select {

		case <-shutdown:
			break loop
		case start := <-bitmarkd.ModeStart:
			if start {
				bitmarkd.Start()
			} else {
				bitmarkd.Stop()
			}
		}

	}
	close(bitmarkd.ModeStart)
}

func (bitmarkd *Bitmarkd) Start() error {
	if bitmarkd.running {
		bitmarkd.log.Errorf("Start bitmarkd failed: %v", ErrBitmarkdIsRunning)
		return ErrBitmarkdIsRunning
	}

	// Check bitmarkConfigFile exists
	bitmarkd.log.Infof("bitmark config file: %s\n", bitmarkd.configFile)
	if !utils.EnsureFileExists(bitmarkd.configFile) {
		bitmarkd.log.Errorf("Start bitmarkd failed: %v", fault.ErrNotFoundConfigFile)
		return fault.ErrNotFoundConfigFile
	}

	bitmarkd.running = true
	stopped := make(chan bool, 1)

	go func() {
		defer func() {
			stopped <- true
		}()
		for bitmarkd.running {
			// start bitmarkd as sub process
			cmd := exec.Command("bitmarkd", "--config-file="+bitmarkd.configFile)
			// start bitmarkd as sub process
			stderr, err := cmd.StderrPipe()
			if err != nil {
				bitmarkd.log.Errorf("Error: %v", err)
				continue
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				bitmarkd.log.Errorf("Error: %v", err)
				continue
			}
			if err := cmd.Start(); nil != err {
				continue
			}

			bitmarkd.process = cmd.Process
			bitmarkd.log.Infof("process id: %d", cmd.Process.Pid)
			stdeReader := bufio.NewReader(stderr)
			stdoReader := bufio.NewReader(stdout)

			go func() {
				for {
					stde, err := stdeReader.ReadString('\n')
					bitmarkd.log.Errorf("bitmarkd stderr: %q", stde)
					if nil != err {
						bitmarkd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			go func() {
				for {
					stdo, err := stdoReader.ReadString('\n')
					bitmarkd.log.Infof("bitmarkd stdout: %q", stdo)
					if nil != err {
						bitmarkd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			if err := cmd.Wait(); nil != err {
				if bitmarkd.running {
					bitmarkd.log.Errorf("bitmarkd has terminated unexpectedly. failed: %v", err)
					bitmarkd.log.Errorf("bitmarkd will be restarted in 5 second...")
					time.Sleep(5 * time.Second)
				}
				bitmarkd.process = nil
			}
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	return nil

}

func (bitmarkd *Bitmarkd) Stop() error {
	if !bitmarkd.running {
		bitmarkd.log.Errorf("Stop bitmarkd failed: %v", ErrBitmarkdIsNotRunning)
		return ErrBitmarkdIsNotRunning
	}
	bitmarkd.running = false

	if bitmarkd.process == nil {
		return nil
	}
	if err := bitmarkd.process.Signal(os.Kill); nil != err {
		bitmarkd.log.Errorf("Send kill to bitmarkd failed: %v", err)
		return err
	}

	bitmarkd.log.Infof("Stop bitmarkd. PID: %d", bitmarkd.process.Pid)
	bitmarkd.process = nil
	return nil
}

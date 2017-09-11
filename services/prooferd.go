// Copyright (c) 2014-2017 Bitmark Inc.
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
	ErrProoferdIsNotRunning = InvalidError("Prooferd is not running")
	ErrProoferdIsRunning    = InvalidError("Prooferd is running")
)

type Prooferd struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
	configFile  string
	process     *os.Process
	running     bool
	ModeStart   chan bool
}

func NewProoferd() *Prooferd {
	return &Prooferd{}
}

func (prooferd *Prooferd) Initialise(configFile string) error {
	prooferd.Lock()
	defer prooferd.Unlock()

	if prooferd.initialised {
		return fault.ErrAlreadyInitialised
	}

	prooferd.configFile = configFile

	prooferd.log = logger.New("service-prooferd")

	prooferd.running = false
	prooferd.ModeStart = make(chan bool, 1)

	// all data initialised
	prooferd.initialised = true
	return nil
}

func (prooferd *Prooferd) Finalise() error {
	prooferd.Lock()
	defer prooferd.Unlock()

	if !prooferd.initialised {
		return fault.ErrNotInitialised
	}

	prooferd.initialised = false
	return nil
}

func (prooferd *Prooferd) IsRunning() bool {
	return prooferd.running
}

func (prooferd *Prooferd) Status() string {
	if prooferd.running {
		return "started"
	} else {
		return "stopped"
	}
}

func (prooferd *Prooferd) Run(args interface{}, shutdown <-chan struct{}) {
loop:
	for {
		select {

		case <-shutdown:
			break loop
		case start := <-prooferd.ModeStart:
			if start {
				prooferd.Start()
			} else {
				prooferd.Stop()
			}
		}

	}
	close(prooferd.ModeStart)
}

func (prooferd *Prooferd) Start() error {
	if prooferd.running {
		prooferd.log.Errorf("Start prooferd failed: %v", ErrProoferdIsRunning)
		return ErrProoferdIsRunning
	}

	// Check prooferdConfigFile exists
	prooferd.log.Infof("prooferd config file: %s\n", prooferd.configFile)
	if !utils.EnsureFileExists(prooferd.configFile) {
		prooferd.log.Errorf("Start prooferd failed: %v", fault.ErrNotFoundConfigFile)
		return fault.ErrNotFoundConfigFile
	}

	prooferd.running = true
	stopped := make(chan bool, 1)

	go func() {
		for prooferd.running {
			// start prooferd as sub process
			cmd := exec.Command("prooferd", "--config-file="+prooferd.configFile)
			// start prooferd as sub process
			stderr, err := cmd.StderrPipe()
			if err != nil {
				prooferd.log.Errorf("Error: %v", err)
				continue
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				prooferd.log.Errorf("Error: %v", err)
				continue
			}
			if err := cmd.Start(); nil != err {
				continue
			}

			prooferd.process = cmd.Process
			prooferd.log.Infof("process id: %d", cmd.Process.Pid)
			stdeReader := bufio.NewReader(stderr)
			stdoReader := bufio.NewReader(stdout)

			go func() {
				for {
					stde, err := stdeReader.ReadString('\n')
					prooferd.log.Errorf("prooferd stderr: %q", stde)
					if nil != err {
						prooferd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			go func() {
				for {
					stdo, err := stdoReader.ReadString('\n')
					prooferd.log.Infof("prooferd stdout: %q", stdo)
					if nil != err {
						prooferd.log.Errorf("Error: %v", err)
						return
					}
				}
			}()

			if err := cmd.Wait(); nil != err {
				if prooferd.running {
					prooferd.log.Errorf("prooferd has terminated unexpectedly. failed: %v", err)
					prooferd.log.Errorf("prooferd will be restarted in 1 second...")
					time.Sleep(time.Second)
				}
				prooferd.process = nil
				stopped <- true
			}
		}
	}()

	// wait for 1 second if cmd has no error then return nil
	time.Sleep(time.Second * 1)
	return nil

}

func (prooferd *Prooferd) Stop() error {
	if !prooferd.running {
		prooferd.log.Errorf("Stop prooferd failed: %v", ErrProoferdIsNotRunning)
		return ErrProoferdIsNotRunning
	}
	prooferd.running = false

	if err := prooferd.process.Signal(os.Kill); nil != err {
		prooferd.log.Errorf("Send kill to prooferd failed: %v", err)
		return err
	}

	prooferd.log.Infof("Stop prooferd. PID: %d", prooferd.process.Pid)
	prooferd.process = nil
	return nil
}

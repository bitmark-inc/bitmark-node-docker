// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bitmark-inc/logger"
)

type GenericError string
type InvalidError GenericError

func (e InvalidError) Error() string { return string(e) }

var (
	ErrInvalidCommandParams = InvalidError("invalid command params")
)

func EnsureFile(filename string) error {
	fileDir := filepath.Dir(filename)

	err := os.MkdirAll(fileDir, 0700)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filename); err != nil {
		_, err := os.OpenFile(filename, os.O_CREATE, 0600)
		return err
	}
	return nil
}

func SimpleCmd(cmdString ...string) (output string, err error) {
	if len(cmdString) == 0 {
		return "", fmt.Errorf("Invaild command strings")
	}
	outBuffer := &bytes.Buffer{}
	cmd := exec.Command(cmdString[0], cmdString[1:]...)
	cmd.Stdout = outBuffer
	cmd.Stderr = outBuffer
	err = cmd.Run()
	return outBuffer.String(), err
}

// logStdOut determines print result from stdout or not
func getCmdOutput(cmd *exec.Cmd, cmdType string, log *logger.L, logStdOut bool) ([]byte, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	if err := cmd.Start(); nil != err {
		return nil, err
	}

	stde, err := ioutil.ReadAll(stderr)
	if nil != err {
		log.Errorf("Error: %v", err)
	}

	stdo, err := ioutil.ReadAll(stdout)
	if nil != err {
		log.Errorf("Error: %v", err)
	}

	log.Errorf("%s %s stderr: %s", cmd.Path, cmdType, stde)
	if logStdOut {
		log.Infof("%s %s stdout: %s", cmd.Path, cmdType, stdo)
	}

	if len(stde) != 0 {
		stderr.Close()
		stdout.Close()
		return nil, errors.New(string(stde))
	}

	if err := cmd.Wait(); nil != err {
		log.Errorf("%s %s failed: %v", cmd.Path, cmdType, err)
		return nil, err
	}
	stderr.Close()
	stdout.Close()

	return stdo, nil
}

func checkRequireStringParameters(params ...string) error {

	for _, param := range params {
		if "" == param || "0" == param {
			return ErrInvalidCommandParams
		}
	}
	return nil
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func hasBracket(address string) bool {
	if strings.Count(address, "[") != 1 {
		return false
	}
	if strings.Count(address, "]") != 1 {
		return false
	}
	return true
}

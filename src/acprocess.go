package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var lines string
var cmd *exec.Cmd

func start() {
	binary := "acServer"
	if runtime.GOOS == "windows" {
		binary = "acServer.exe"
	}

	fpath := filepath.Join(TempFolder, binary)
	if _, err := os.Stat(fpath); errors.Is(err, os.ErrNotExist) {
		log.Print("Could not find executable: ", fpath, err)
	}
	cmd = exec.Command(fpath)

	if runtime.GOOS != "windows" {
		err := os.Chmod(fpath, 0755)

		if err != nil {
			log.Print("Could not chmod executable: ", fpath, err)
		}
	}

	cmd.Dir = TempFolder
	var stdOut, _ = cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		log.Print("Could not start executable: ", fpath, err)
	}

	out := bufio.NewReader(stdOut)

	lines = ""
	go func() {
		for cmd != nil {
			b, _, err := out.ReadLine()

			if err != nil {
				log.Print("Error acServer process stdout: ", err)
			}

			if len(b) > 0 {
				lines += string(b) + "\n"
			}
		}
	}()

}
func getContent() string {
	return lines
}

func isRunning() bool {
	if cmd != nil && cmd.ProcessState != nil {
		return cmd.ProcessState.Exited()
	} else if cmd != nil && cmd.Process != nil {
		return true
	}
	return false
}

func stop() {
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()

		if err != nil {
			log.Print("Error killing acServer process: ", err)
		}
		cmd = nil
		Status.Players = 0
	}
}

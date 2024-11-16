package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var lines string
var cmd *exec.Cmd

func Start() {
	binary := "acServer"
	if runtime.GOOS == "windows" {
		binary = "acServer.exe"
	}

	fpath := filepath.Join(TempFolder, binary)
	cmd = exec.Command(fpath)

	if runtime.GOOS != "windows" {
		os.Chmod(fpath, 0755)
	}

	cmd.Dir = TempFolder
	var stdOut, _ = cmd.StdoutPipe()
	cmd.Start()

	out := bufio.NewReader(stdOut)

	lines = ""
	go func() {
		for cmd != nil {
			b, _, _ := out.ReadLine()

			if len(b) > 0 {
				lines += string(b) + "\n"
			}
		}
	}()

}
func Get_Content() string {
	return lines
}

func Is_Running() bool {
	if cmd != nil && cmd.ProcessState != nil {
		return cmd.ProcessState.Exited()
	} else if cmd != nil && cmd.Process != nil {
		return true
	}
	return false
}

func Stop() {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd = nil
		Status.Players = 0
	}
}

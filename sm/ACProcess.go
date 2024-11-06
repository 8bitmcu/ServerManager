package sm

import (
	"bufio"
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

	cmd = exec.Command(filepath.Join(Dba.Basepath(), "server", binary))
	cmd.Dir = filepath.Join(Dba.Basepath(), "server")
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
	}
}

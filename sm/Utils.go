package sm

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// open opens the specified URL in the default browser of the user.
func Open_URL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func Print_Interface(t interface{}) {
	s, _ := json.MarshalIndent(t, "", "\t")
	fmt.Print(string(s))
}

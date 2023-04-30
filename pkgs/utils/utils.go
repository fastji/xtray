package utils

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func NormalizeSSR(s string) (r string) {
	r = strings.ReplaceAll(s, "-", "+")
	r = strings.ReplaceAll(r, "_", "/")
	return r
}

func DecodeBase64(str string) (res string) {
	count := len(str) % 4
	for i := 0; i < count; i++ {
		str += "="
	}
	if s, err := base64.StdEncoding.DecodeString(str); err == nil {
		res = string(s)
	}
	return
}

func StringToReader(str string) io.Reader {
	return bytes.NewReader([]byte(str))
}

func GetHomeDir() (homeDir string) {
	u, err := user.Current()
	if err != nil {
		fmt.Println("[CurrentUser]", err)
		return
	}
	return u.HomeDir
}

const (
	IsChildEnv     = "XTRAY_IS_CHILD_PROCESS"
	IsChildProcess = "XTRAY_IS_CHILD_PROCESS=true"
)

func DaemonizeInit() {
	isChild := os.Getenv(IsChildEnv)
	if isChild == "" {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		cmd.Env = append(os.Environ(), IsChildProcess)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
}

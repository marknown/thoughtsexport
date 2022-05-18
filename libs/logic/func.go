package logic

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// FindExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather agressive search, which is the same in all
// systems. That may make it a bit slow, but it will only be run when creating a
// new ExecAllocator.
func FindExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"headless_shell",
		"headless-shell",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"/usr/bin/google-chrome",

		// Windows
		"chrome",
		"chrome.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,

		// Mac
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return ""
}

// ReadOutput grabs the websocket address from chrome's output, returning as
// soon as it is found. All read output is forwarded to forward, if non-nil.
// done is used to signal that the asynchronous io.Copy is done, if any.
func ReadOutput(rc io.ReadCloser, forward io.Writer, done func()) (wsURL string, _ error) {
	prefix := []byte("DevTools listening on")
	var accumulated bytes.Buffer
	bufr := bufio.NewReader(rc)
readLoop:
	for {
		line, err := bufr.ReadBytes('\n')
		if err != nil {
			return "", fmt.Errorf("chrome failed to start:\n%s",
				accumulated.Bytes())
		}
		if forward != nil {
			if _, err := forward.Write(line); err != nil {
				return "", err
			}
		}

		if bytes.HasPrefix(line, prefix) {
			line = line[len(prefix):]
			// use TrimSpace, to also remove \r on Windows
			line = bytes.TrimSpace(line)
			wsURL = string(line)
			break readLoop
		}
		accumulated.Write(line)
	}
	if forward == nil {
		// We don't need the process's output anymore.
		rc.Close()
	} else {
		// Copy the rest of the output in a separate goroutine, as we
		// need to return with the websocket URL.
		go func() {
			io.Copy(forward, bufr)
			done()
		}()
	}
	return wsURL, nil
}

// Exit func
func ExitFunc(process *os.Process) {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("退出", s)
				err := process.Kill()
				if nil != err {
					log.Fatal(err)
				}
				os.Exit(0)
			default:
				log.Println("other", s)
			}
		}
	}()
}

// 删除文件夹下所有内容
func ClearDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCurrentDirectory 获取可执行文件的当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// Exec 执行系统命令，并返回输出
func Exec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if nil != err {
		return "", err
	}

	return output.String(), nil
}

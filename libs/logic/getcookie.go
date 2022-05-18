package logic

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// GetLoginCookieString 获取指定链接，指定cookie的登录状态
func GetLoginCookieString(loginURL string, waitCookieKey string) string {
	// 临时目录
	tempDir, err := ioutil.TempDir("", "chromedp-user-data")
	if err != nil {
		log.Fatal(err)
	}

	tempDir2, err := ioutil.TempDir("", "chromedp-disk-cache")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir2)

	procCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	execPath := FindExecPath()
	if "" == execPath {
		log.Fatal(errors.New("chrome path is not found"))
	}
	cmd := exec.CommandContext(procCtx, execPath,
		// TODO: deduplicate these with allocOpts in chromedp_test.go
		// "--incognito",
		// "--headless",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-gpu",
		"--no-sandbox",

		// TODO: perhaps deduplicate this code with ExecAllocator
		"--user-data-dir="+tempDir,
		"--disk-cache-dir="+tempDir2,
		"--remote-debugging-port=9222",
		"--remote-debugging-address=0.0.0.0",
		`--user-agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"`,
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	wsURL, err := ReadOutput(stderr, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(wsURL)

	// 注册退出函数
	ExitFunc(cmd.Process)

	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), wsURL)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	// 首先打开一个实例
	if err := chromedp.Run(taskCtx, chromedp.Navigate(loginURL)); err != nil {
		log.Fatal(errors.New("Navigate: " + err.Error()))
	}

	cookie := WaitLoginReturnCookieString(taskCtx, waitCookieKey)

	return cookie
}

// WaitLoginReturnCookieString 循环监听指定cookie key的状态
func WaitLoginReturnCookieString(ctx context.Context, waitCookieKey string) string {
	cookieStr := ""
	waitCookieKeyExist := false

	// 循环获取Cookie，并查看是否登录
	for {
		var cookies = []*network.Cookie{}

		// 从 chromedp 里获取cookie
		if err := chromedp.Run(ctx,
			func() chromedp.ActionFunc {
				return func(ctx context.Context) (err error) {
					cookies, _ = network.GetAllCookies().Do(ctx)
					return
				}
			}(),
		); err != nil {
			log.Fatal(errors.New("GetAllCookies Error: " + err.Error()))
		}

		// 如果存在 waitCookieKey 说明已经登录
		for _, cookie := range cookies {
			cookieStr += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
			if cookie.Name == waitCookieKey {
				waitCookieKeyExist = true
			}
		}
		if waitCookieKeyExist {
			break
		} else {
			cookieStr = ""
		}

		time.Sleep(1 * time.Second)
	}

	return cookieStr
}

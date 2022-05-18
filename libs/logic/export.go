package logic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"thoughtsexport/libs/request"
	"thoughtsexport/libs/utils"
)

var rootPath = ""

func ExportOne(url string, cookie string, fileType string) {
	parts := strings.Split(url, "/")
	hashSpace := parts[4]

	req := request.NewRequest(cookie, hashSpace)

	workspace, err := req.GetWorkspace(hashSpace)
	if nil != err {
		panic(err)
	}

	// 如果没有导出权限，开启导出权限
	needCloseOutput := false
	if workspace.WorkspaceSecurity.DisableOutput {

		succeed, err := req.EnableOutput(hashSpace, true)
		if nil != err {
			panic(err)
		}

		if !succeed {
			panic("本文档无法下载，开启导出权限失败。请文档所有者在本工具登录后再尝试")
		}

		// 之前没有导出权限，现在临时开启，下载完成要关闭导出权限
		needCloseOutput = true
	}

	prefixPath := fmt.Sprintf("%s/%s", workspace.Organization.Name, workspace.Name)
	SetRootPath(GetCurrentDirectory() + "/" + prefixPath)
	fmt.Printf("所有文件将保存至 %s\n", GetRootPath())

	nodes, err := req.GetAllNodes(hashSpace, "")
	if nil != err {
		panic(err)
	}
	fmt.Printf("分析完成 %s\n", prefixPath)

	total := len(nodes)
	counter := 0
	for _, node := range nodes {
		counter++
		fmt.Printf("当前进度 %d/%d [%.2f%%] 正在下载文档 %s \r", counter, total, float64(counter)*float64(100)/float64(total), node.Path)

		if node.Type == "folder" {
			// log.Println("我是空目录" + node.Path)
			CreateDir(node.Path + "/")
		} else if node.Type == "document" {
			// 下载 docx
			if fileType == "all" || fileType == "docx" {
				downloadInfo, err := req.GetDownloadUrl(node.ID, node.Path, "docx")
				if nil != err {
					LogDownloadFailedInfo(node, err)
					continue
				}

				_, err = DownloadFile(downloadInfo.DownURL, GetRootPath()+downloadInfo.FullPath)
				if nil != err {
					LogDownloadFailedInfo(node, err)
					continue
				}
			}

			// 下载 html
			if fileType == "all" || fileType == "html" {
				downloadInfo, err := req.GetDownloadUrl(node.ID, node.Path, "html")
				if nil != err {
					LogDownloadFailedInfo(node, err)
					continue
				}

				_, err = DownloadFile(downloadInfo.DownURL, GetRootPath()+downloadInfo.FullPath)
				if nil != err {
					LogDownloadFailedInfo(node, err)
					continue
				}
			}
		} else {
			downloadInfo, err := req.GetDownloadUrlByDetail(node.ID, node.Path)
			if nil != err {
				LogDownloadFailedInfo(node, err)
				continue
			}

			_, err = DownloadFile(downloadInfo.DownURL, GetRootPath()+downloadInfo.FullPath)
			if nil != err {
				LogDownloadFailedInfo(node, err)
				continue
			}
		}

		fmt.Printf("[已完成] %s\n", node.Path)
	}

	fmt.Printf("所有文件已保存至 %s\n", GetRootPath())

	// 关闭导出权限
	if needCloseOutput {
		_, err := req.EnableOutput(hashSpace, false)
		if nil != err {
			panic(err)
		}
	}
}

func CreateDir(fullPath string) error {
	dirPath := path.Dir(fullPath)

	_, err := os.Stat(dirPath)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}

		// err = os.Chmod(dirPath, os.ModeDir)
		// if err != nil {
		// 	return err
		// }
	}

	return err
}

func DownloadFile(url string, filepath string) (int64, error) {
	if utils.FileExist(filepath) {
		return 0, nil
	}

	CreateDir(filepath)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(file, resp.Body)

	return n, err
}

func LogFailedInfo(info string) {
	path := GetRootPath() + "/下载失败的文件清单.txt"
	utils.FileAppend(path, info)
}

func LogDownloadFailedInfo(node *request.Node, err error) {
	info := fmt.Sprintf("%s %s\n", node.Path, err.Error())
	LogFailedInfo(info)
	fmt.Println(info)
}

func SetRootPath(path string) {
	rootPath = path
}

func GetRootPath() string {
	return rootPath
}

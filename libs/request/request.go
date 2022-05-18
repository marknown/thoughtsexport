package request

import (
	"errors"
	"fmt"
	"thoughtsexport/libs/utils"
	"time"

	"github.com/marknown/ohttp"
)

type request struct {
	cookie    string
	firstHash string
}

// NodesResponse 响应
type NodesResponse struct {
	NextPageToken string  `json:"nextPageToken"`
	Nodes         []*Node `json:"result"`
}

// Node 响应结构体里的 Node
type Node struct {
	ID          string `json:"_id"`
	ParentId    string `json:"_parentId"`
	WorkspaceId string `json:"_workspaceId"`
	Created     time.Time
	Title       string
	Type        string
	WithChild   bool
	Path        string
	Info        struct {
		DownloadUrl string
		FileType    string
		FileName    string
	}
}

// Workspaces 响应结构体里的 Workspaces
type Workspaces struct {
	ID           string `json:"_id"`
	Created      time.Time
	Name         string
	Organization struct {
		ID   string `json:"_id"`
		Name string
	}
	WorkspaceSecurity WorkspaceSecurity
}

type WorkspaceSecurity struct {
	DisableShare      bool
	DisableMove       bool
	DisableOutput     bool
	DisableSharespace bool
}

type WorkspaceSecurityResponse struct {
	ID                string `json:"_id"`
	WorkspaceSecurity WorkspaceSecurity
}

// NodeDownload 文档的下载信息
type NodeDownload struct {
	FileType string
	FullPath string
	DownURL  string
}

// IDResponse
type IDResponse struct {
	ID string `json:"id"`
}

// DownInfoResponse
type DownInfoResponse struct {
	ConvertProcess int
	Message        struct {
		DownloadUrl string
		Error       string
	}
}

func NewRequest(cookie string, firstHash string) *request {
	return &request{
		cookie:    cookie,
		firstHash: firstHash,
	}
}

func (r *request) GetWorkspace(hash string) (*Workspaces, error) {
	req := ohttp.InitSetttings()
	req.Timeout = 10 * time.Second
	req.IsAajx = true
	req.Referer = "https://thoughts.teambition.com"
	req.Cookies = r.cookie

	url := fmt.Sprintf("https://thoughts.teambition.com/api/workspaces/%s?pageSize=1000&_=%d", hash, utils.UnixTimstampMillisecond())
	// fmt.Println(url)
	content, _, err := req.Get(url)

	//fmt.Println(content, err)

	if nil != err {
		return nil, err
	}

	var result Workspaces
	err = utils.JSONToStruct(content, &result)

	if nil != err {
		return nil, err
	}

	return &result, nil
}

func (r *request) EnableOutput(hash string, enable bool) (bool, error) {
	settings := ohttp.InitSetttings()
	settings.Timeout = 10 * time.Second
	settings.IsAajx = true
	settings.Referer = "https://thoughts.teambition.com"
	settings.Cookies = r.cookie
	settings.ContentType = "application/json; charset=utf-8"

	// disableOutput 和 enable 刚好相反
	disableStr := "true"
	if enable {
		disableStr = "false"
	}

	params := fmt.Sprintf(`{"optTarget":"disableOutput","optVal":%s}`, disableStr)
	url := fmt.Sprintf("https://thoughts.teambition.com/api/workspaces/%s/workspaceSecurity", hash)

	req, err := settings.NewRequest("PUT", url, params)

	if nil != err {
		return false, err
	}

	resp, err := settings.Do(req)

	if nil != err {
		return false, err
	}

	content, err := resp.ContentString()

	if nil != err {
		return false, err
	}

	var result WorkspaceSecurityResponse
	err = utils.JSONToStruct(content, &result)

	if nil != err {
		return false, err
	}

	return enable != result.WorkspaceSecurity.DisableOutput, nil
}

func (r *request) GetAllNodes(hashSpace string, prefixPath string) ([]*Node, error) {
	var allNodes []*Node
	nodes, err := r.GetNodesByHash(hashSpace, prefixPath)
	fmt.Printf("正在分析 %s\n", prefixPath)

	if nil != err {
		fmt.Println(err)
	}
	for _, node := range nodes {
		if node.WithChild {
			nodes1, err := r.GetAllNodes(node.ID, node.Path)
			if nil != err {
				fmt.Println(err)
			}
			allNodes = append(allNodes, nodes1...)
		}

		allNodes = append(allNodes, node)
	}

	return allNodes, nil
}

func (r *request) GetNodesByHash(hash string, prefixPath string) ([]*Node, error) {
	req := ohttp.InitSetttings()
	req.Timeout = 10 * time.Second
	req.IsAajx = true
	req.Referer = "https://thoughts.teambition.com"
	req.Cookies = r.cookie

	var parentHash = ""
	if r.firstHash != hash {
		parentHash = fmt.Sprintf("&_parentId=%s", hash)
	}

	url := fmt.Sprintf("https://thoughts.teambition.com/api/workspaces/%s/nodes?pageSize=1000%s&_=%d", r.firstHash, parentHash, utils.UnixTimstampMillisecond())
	// fmt.Println(url)
	content, _, err := req.Get(url)

	// fmt.Println(content, err)

	if nil != err {
		return nil, err
	}

	var result NodesResponse
	err = utils.JSONToStruct(content, &result)

	if nil != err {
		return nil, err
	}

	for _, node := range result.Nodes {
		node.Path = fmt.Sprintf("%s/%s", prefixPath, node.Title)
	}

	return result.Nodes, nil
}

func (r *request) GetDownloadUrl(hash string, prefixPath string, fileType string) (*NodeDownload, error) {
	req := ohttp.InitSetttings()
	req.Timeout = 10 * time.Second
	req.IsAajx = true
	req.Referer = "https://thoughts.teambition.com"
	req.Cookies = r.cookie

	dict := map[string]string{
		"docx": "docx",
		"html": "zip",
	}

	if _, ok := dict[fileType]; !ok {
		return nil, errors.New("fileType不合法")
	}

	url := fmt.Sprintf("https://thoughts.teambition.com/convert/api/nodes/%s/export:%s?pageSize=1000&_=%d", hash, fileType, utils.UnixTimstampMillisecond())

	// fmt.Println(url)
	content, _, err := req.Get(url)

	// fmt.Println(content, err)

	if nil != err {
		return nil, err
	}

	var result IDResponse
	err = utils.JSONToStruct(content, &result)

	if nil != err {
		return nil, err
	}

	if result.ID == "" {
		return nil, fmt.Errorf("没有获取到 %s 文档的下载ID", prefixPath)
	}

	url = fmt.Sprintf("https://thoughts.teambition.com/convert/api/exportDocx:polling?pageSize=1000&id=%s&_=%d", result.ID, utils.UnixTimstampMillisecond())
	downInfoResponse, err := r.PollingDownloadUrl(url)

	if nil != err {
		return nil, err
	}

	// 扩展名
	ext := dict[fileType]
	nodeDownload := &NodeDownload{
		FileType: "docx",
		FullPath: prefixPath + "." + ext,
		DownURL:  downInfoResponse.Message.DownloadUrl,
	}

	return nodeDownload, nil
}

func (r *request) PollingDownloadUrl(url string) (*DownInfoResponse, error) {
	req := ohttp.InitSetttings()
	req.Timeout = 10 * time.Second
	req.IsAajx = true
	req.Referer = "https://thoughts.teambition.com"
	req.Cookies = r.cookie

	// fmt.Println(url)
	var err error
	var content string
	var result DownInfoResponse
	var maxRetry = 5
	var retryCounter = 0
	for {
		retryCounter++
		time.Sleep(1 * time.Second)
		content, _, err = req.Get(url)

		// fmt.Println(content, err)

		if nil != err {
			break
		}

		err = utils.JSONToStruct(content, &result)

		if nil != err {
			break
		}

		// 成功获取到下载链接
		if result.ConvertProcess == 1 && result.Message.DownloadUrl != "" {
			break
		}

		// 无法获取下载链接
		if result.ConvertProcess == -1 {
			err = fmt.Errorf("获取下载链接失败[%s]", result.Message.Error)
			break
		}

		// 避免死循环
		if retryCounter >= maxRetry {
			err = fmt.Errorf("获取下载链接失败[已尝试%d次]", maxRetry)
			break
		}
	}

	return &result, err
}

func (r *request) GetDownloadUrlByDetail(hash string, prefixPath string) (*NodeDownload, error) {
	req := ohttp.InitSetttings()
	req.Timeout = 10 * time.Second
	req.IsAajx = true
	req.Referer = "https://thoughts.teambition.com"
	req.Cookies = r.cookie

	url := fmt.Sprintf("https://thoughts.teambition.com/api/workspaces/%s/nodes/%s?pageSize=1000&_=%d", r.firstHash, hash, utils.UnixTimstampMillisecond())

	// fmt.Println(url)
	content, _, err := req.Get(url)

	// fmt.Println(content, err)

	if nil != err {
		return nil, err
	}

	var result Node
	err = utils.JSONToStruct(content, &result)

	if nil != err {
		return nil, err
	}

	if result.Info.DownloadUrl == "" {
		return nil, fmt.Errorf("没有获取到 %s 文档的下载ID", prefixPath)
	}

	nodeDownload := &NodeDownload{
		FileType: result.Info.FileType,
		FullPath: prefixPath,
		DownURL:  result.Info.DownloadUrl,
	}

	return nodeDownload, nil
}

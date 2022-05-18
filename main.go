package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"thoughtsexport/libs/logic"
	"thoughtsexport/libs/utils"

	"github.com/marknown/util"
)

var isDownloading = false

func main() {
	// 打开本地浏览器
	errOpen := util.Open("http://127.0.0.1:43821/submit/url")
	if errOpen != nil {
		fmt.Println(errOpen.Error())
	}

	http.HandleFunc("/submit/url", func(w http.ResponseWriter, r *http.Request) {
		if isDownloading {
			fmt.Fprintf(w, "%s", "下载任务正在进行中，本页面可以关闭，具体看黑色窗口的输出")
			return
		}

		const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>请填写要下载的知识库地址</title>
	<style>
		body{margin:0}
		#main{width:800px;margin:auto;}
		#main div {line-height:20px;}
		#tips{color:red;}
		.url-input{width:538px;height:20px;font-size:12px;}
	</style>
	<script type="text/javascript">
	var Ajax = {
		post: function(url, data, fn) { // data 应为'a=a1&b=b1'这种字符串格式
			var xhr = new XMLHttpRequest();
			xhr.open("POST", url, true);
			// 添加http头，发送信息至服务器时内容编码类型
			xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
			xhr.onreadystatechange = function() {
				if(xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
					fn.call(this, xhr.responseText);
				}
			};
			xhr.send(data);
		}
	}

	function submitURL() {
		let url = document.getElementById('url').value
		let type = document.getElementById('type').value
		let parts = url.split("/")
		if (parts[2] != "thoughts.teambition.com") {
			showTips("链接必须是 https://thoughts.teambition.com/workspaces/xxxxxx/overview 格式");
			return;
		}
		if (parts[3] != "workspaces") {
			showTips("链接必须是 https://thoughts.teambition.com/workspaces/xxxxxx/overview 格式");
			return;
		}

		Ajax.post("/receive/url", "url="+url+"&type="+type, function (response) {
			try {
				var obj = JSON.parse(response)
				if (!obj.success) {
					showTips(obj.message)
				} else {
					hideAction()
					showTips(obj.message)
				}
			} catch(e) {
				// alert(e + "\norigin：" + response)
				showTips(e + "\norigin：" + response)
			}
		});
	}

	function showTips(msg) {
		let tipsObj = document.getElementById('tips')
		tipsObj.innerText = msg
	}
	
	function hideAction() {
		let action = document.getElementById('action')
		action.style.display = "none"
	}
</script>
</head>
<body>
	<div id="main">
		<div>请在下面的输入框中输入要下载的知识库地址。注意事项：1. 请安装 Chrome 谷歌浏览器，下一步会用到。</div>
		<div style="font-size:12px;">格式：https://thoughts.teambition.com/workspaces/xxxxxx/overview</div>
		<div id="action">
			<input id="url" class="url-input"/>
			格式<select id="type">
				<option value="docx">仅word文档</option>
				<option value="html">仅html格式</option>
				<option value="all">以上格式一起导出</option>
			</select>
			<input type="button" value="开始导出" onclick="submitURL()"/>
		</div>
		<div id="tips"></div>
	</div>
</body>
</html>`

		tmpl := template.Must(template.New("image").Parse(tpl))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/receive/url", func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		res := response{
			Success: true,
			Message: "下载任务正在进行中，本页面可以关闭。请在接下来弹出的窗口里登录所思账号，只有登录了才能开始导出文档。",
		}

		url := r.FormValue("url")
		fileType := r.FormValue("type")

		parts := strings.Split(url, "/")

		if isDownloading {
			res.Success = true
			res.Message = "下载任务正在进行中，请不要重复操作，本页面可以关闭"
		} else if !utils.SliceContainStr([]string{"docx", "html", "all"}, fileType) {
			res.Success = false
			res.Message = fmt.Sprintf("%s类型不支持导出", fileType)
		} else if len(parts) < 5 {
			res.Success = false
			res.Message = "链接必须是 https://thoughts.teambition.com/workspaces/xxxxxx/overview 格式"
		} else if parts[2] != "thoughts.teambition.com" {
			res.Success = false
			res.Message = "链接必须是 https://thoughts.teambition.com/workspaces/xxxxxx/overview 格式"
		} else if parts[3] != "workspaces" {
			res.Success = false
			res.Message = "链接必须是 https://thoughts.teambition.com/workspaces/xxxxxx/overview 格式"
		}

		json, err := utils.ToJson(res)
		if err != nil {
			res.Success = false
			res.Message = "json格式化失败：" + err.Error()
		}

		// fmt.Println(url)
		// fmt.Println(fileType)

		// 如果检查成功后，就开始导出文档
		if res.Success {
			// 异步导出文档
			go func(url string) {
				// 标记为下载中
				isDownloading = true

				var loginURL = url
				// 调用 chrome 来登录，并获取登录后的 cookie ，请保证本机安装了 chome 浏览器
				cookie := logic.GetLoginCookieString(loginURL, "TB_ACCESS_TOKEN")
				// 开始导出一个知识库
				logic.ExportOne(loginURL, cookie, fileType)

				// 标记为非下载中
				isDownloading = false

				// 下载完成后，退出程序
				os.Exit(0)
			}(url)
		}

		fmt.Fprintf(w, "%s", json)
	})

	err := http.ListenAndServe(":43821", nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}

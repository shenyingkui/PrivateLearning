package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func UiInit(fileName string, r *gin.Engine) {

	// 创建应用程序实例
	a := app.New()
	w := a.NewWindow("私人定制播放器")
	w.Resize(fyne.Size{800, 600})

	str := binding.NewString()
	str.Set("我是初始化值")

	share := widget.NewLabel("分享链接")
	text := widget.NewLabelWithData(str)
	w.SetContent(text)

	// 创建一个按钮并设置文本

	button := widget.NewButton("选择文件", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			// 这里可以添加处理文件的代码，例如保存文件路径
			dialog.ShowInformation("文件上传成功", "已选择文件: "+reader.URI().Name(), w)
			fmt.Println("选择文件", reader.URI().Path())
			//调用ffmpeg 方法

			// 获取文件信息
			fileInfo, err := os.Stat(reader.URI().Path())
			// 获取文件权限
			fileMode := fileInfo.Mode()

			fmt.Printf("User has read permission for the file: %s\n", fileMode)
			//获取配置
			viper.SetConfigName("conf")
			viper.SetConfigType("yml")
			viper.AddConfigPath(".")

			//读取配置文件
			if err := viper.ReadInConfig(); err != nil {
				fmt.Println("error reading config file, %s", err)
				return
			}
			port := viper.GetInt("server.port")
			fmt.Println("server port", port)
			filePath := viper.GetString("datapath.userpath")
			//ffmpeg -i test.mp4 -force_key_frames "expr:gte(t,n_forced*2)" -strict -2 -c:a aac -c:v libx264 -hls_time 2 -f hls index.m3u8
			separator := filepath.Separator
			tempPath := strings.ReplaceAll(reader.URI().Path(), "/", "\\\\")
			fmt.Println(tempPath)
			// 构建ffmpeg命令
			go runffmpeg(tempPath, filePath, separator, fileName, err)

			//拼接视频播放地址
			url := "http://"
			url = url + "localhost" + ":" + strconv.Itoa(port) + "/index/" + fileName
			fmt.Println("分享视频地址:" + url)
			str.Set(url)
		}, w)
		//fd.SetFilter(storage.NewExtensionFileFilter([]string{".MP4", ".avi"}))
		fd.Show()
	})
	con := container.New(layout.NewVBoxLayout())
	//link := widget.NewHyperlink("云播放视频", &url.URL{Scheme: "http", Host: "www.baid1u.com"})

	copyButton := widget.NewButton("复制分享链接", func() {
		// 获取按钮的文本内容
		text, _ := str.Get()
		// 将文本转换为字符串
		textStr := string(text)
		fmt.Println(textStr)
		// 将字符串复制到系统剪贴板
		clipboard.Write(0, []byte(textStr))

	})

	explanation := widget.NewRichTextWithText("1. 配置ffmpeg 环境变量\n2. 配置外网地址映射 https://www.cpolar.com/docs \n3.选择播放文件\n4.复制分享链接,发送给好友 ")

	con.Add(button)
	con.Add(copyButton)
	con.Add(share)
	con.Add(text)

	//con.Add(buttonShare)
	//con.Add(link)
	con.Add(explanation)

	w.CenterOnScreen() //窗口居中

	w.SetContent(con)
	// 显示并运行应用程序
	w.ShowAndRun()
}

func runffmpeg(tempPath string, filePath string, separator int32, fileName string, err error) {
	cmd := exec.Command("ffmpeg.exe", "-i", tempPath,
		"-codec", "copy",
		"-level", "3.0",
		"-start_number", "0",
		"-hls_list_size", "0",
		"-c:v", "libx264", "-hls_time", "20", "-f", "hls", filePath+string(separator)+fileName+".m3u8")

	fmt.Println(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 运行命令
	err1 := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err1)
	}

	// 打印输出结果
	fmt.Println("Output:", cmd)
}

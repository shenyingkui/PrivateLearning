package view

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

//go:embed play.tmpl
var templates embed.FS

type Page struct {
	Filename string `uri:"filename" binding:"required"`
}

func Init(httpurl string, r *gin.Engine) {
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
	tmpl, err := template.ParseFS(templates, "play.tmpl")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)
	r.StaticFS("/static", http.Dir(filePath))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/index/:filename", func(c *gin.Context) {
		var page Page
		if c.ShouldBindUri(&page) == nil {
			log.Println("====== Only Bind By Query String ======")

		}
		file := &page.Filename
		log.Println(file)
		c.HTML(http.StatusOK, "play.tmpl", gin.H{
			"playUrl": "/static/" + *file + ".m3u8",
		})
	})

	serverAdress := fmt.Sprintf(":%d", port)
	r.Run(serverAdress)
}

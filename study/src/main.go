package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"study/src/ui"
	"study/src/view"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	fileName := generateRandomString(10)

	r := gin.Default()
	go view.Init(fileName, r)
	//r.LoadHTMLFiles("templates/play.tmpl")
	ui.UiInit(fileName, r)

}

package main

import (
	"fmt"
	//"os"
	"log"
	"net/http"
	//"database/sql"

	"github.com/gin-gonic/gin"
	//"github.com/kaptinlin/jsonrepair"

	"main/sm"
)

func main() {

	dba := sm.Open("hello.db")
	dba.Apply_Schema()
	cnt := dba.Table_Exists("user_config")

	fmt.Print(cnt)

	router := gin.Default()
	router.LoadHTMLGlob("htm/*")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.htm", gin.H{
			"page":       "index",
			"is_running": false,
		})
	})

	router.GET("/config", func(c *gin.Context) {

		c.HTML(http.StatusOK, "config.htm", gin.H{
			"page": "config",
			"form": dba.Select_Config(),
		})
	})

	router.Run(":3030")
}

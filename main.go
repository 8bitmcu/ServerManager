package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"main/sm"
	"strconv"
)

var dba sm.Dbaccess

func main() {

	dba = sm.Open("hello.db")
	dba.Apply_Schema()
	cnt := dba.Table_Exists("user_config")

	fmt.Print(cnt)

	router := gin.Default()
	router.Static("/static", "./static")

	router.SetFuncMap(template.FuncMap{
		"derefStr": func(t *string) string {
			if t == nil {
				return ""
			}
			return *t
		},
		"derefInt": func(t *int) string {
			if t == nil {
				return ""
			}
			return strconv.Itoa(*t)
		},
		"toInt": func(t *int) int {
			if t == nil {
				return 0
			}
			return *t
		},
		"toJsBool": func(t *int) bool {
			if t == nil {
				return true
			}
			if *t > 0 {
				return true
			}
			return false
		},
	})

	router.LoadHTMLGlob("htm/*")

	//sm.Parse_Weathers(dba)
	sm.Parse_Tracks(dba)
	//sm.Parse_Cars(dba)

	sm.Dba = dba

	router.GET("/", sm.Route_Index)

	router.GET("/config", sm.Route_Config)
	router.POST("/config", sm.Route_Config)

	router.GET("/content", sm.Route_Content)
	router.POST("/content", sm.Route_Content)

	router.GET("/difficulty", sm.Route_Difficulty)
	router.POST("/difficulty", sm.Route_Difficulty)
	router.GET("/difficulty/:id", sm.Route_Difficulty)
	router.POST("/difficulty/:id", sm.Route_Difficulty)

	router.GET("/class", sm.Route_Class)
	router.POST("/class", sm.Route_Class)
	router.GET("/class/:id", sm.Route_Class)
	router.POST("/class/:id", sm.Route_Class)

	router.GET("/session", sm.Route_Session)
	router.POST("/session", sm.Route_Session)
	router.GET("/session/:id", sm.Route_Session)
	router.POST("/session/:id", sm.Route_Session)

	router.GET("/time", sm.Route_Time)
	router.POST("/time", sm.Route_Time)
	router.GET("/time/:id", sm.Route_Time)
	router.POST("/time/:id", sm.Route_Time)

	router.GET("/event", sm.Route_Event)
	router.POST("/event", sm.Route_Event)
	router.GET("/event/:id", sm.Route_Event)
	router.POST("/event/:id", sm.Route_Event)

	router.GET("/api/car/image/:car/:skin", sm.API_Car_Image)

	router.GET("/api/track/preview/:track/:config", sm.API_Track_Preview_Image)
	router.GET("/api/track/preview/:track", sm.API_Track_Preview_Image)

	router.GET("/api/track/outline/:track/:config", sm.API_Track_Outline_Image)
	router.GET("/api/track/outline/:track", sm.API_Track_Outline_Image)

	router.GET("/api/difficulty/:id", sm.API_Difficulty)
	router.GET("/api/session/:id", sm.API_Session)
	router.GET("/api/class/:id", sm.API_Class)
	router.GET("/api/time/:id", sm.API_Time)

	router.Run(":3030")
}

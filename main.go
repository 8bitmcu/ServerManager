package main

import (
	"html/template"
	"io"
	"log"
	"main/sm"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
)

func loadTemplate(t *template.Template) error {
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".htm") {
			continue
		}
		h, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return err
		}
	}
	return nil
}

func findFile(filePath string) *assets.File {
	for _, file := range Assets.Files {
		if file.Path == filePath {
			return file
		}
	}
	return nil
}

func main() {

	dba := sm.Open("smdata.db")


	dba.Apply_Schema(findFile("/schema.sql"))
	sm.Dba = dba

	router := gin.Default()

	funcMap := template.FuncMap{
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
				return false
			}
			if *t > 0 {
				return true
			}
			return false
		},
	}

	//router.SetFuncMap(funcMap)
	//router.LoadHTMLGlob("*")

	t := template.New("")
	t.Funcs(funcMap)
	err := loadTemplate(t)
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	router.StaticFS("/static", Assets)

	router.GET("/", sm.Route_Index)

	router.GET("/config", sm.Route_Config)
	router.POST("/config", sm.Route_Config)

	router.GET("/content", sm.Route_Content)
	router.POST("/content", sm.Route_Content)

	router.GET("/difficulty", sm.Route_Difficulty)
	router.POST("/difficulty", sm.Route_Difficulty)
	router.GET("/difficulty/:id", sm.Route_Difficulty)
	router.POST("/difficulty/:id", sm.Route_Difficulty)
	router.GET("/difficulty/delete/:id", sm.Route_Delete_Difficulty)

	router.GET("/class", sm.Route_Class)
	router.POST("/class", sm.Route_Class)
	router.GET("/class/:id", sm.Route_Class)
	router.POST("/class/:id", sm.Route_Class)
	router.GET("/class/delete/:id", sm.Route_Delete_Class)

	router.GET("/session", sm.Route_Session)
	router.POST("/session", sm.Route_Session)
	router.GET("/session/:id", sm.Route_Session)
	router.POST("/session/:id", sm.Route_Session)
	router.GET("/session/delete/:id", sm.Route_Delete_Session)

	router.GET("/time", sm.Route_Time)
	router.POST("/time", sm.Route_Time)
	router.GET("/time/:id", sm.Route_Time)
	router.POST("/time/:id", sm.Route_Time)
	router.GET("/time/delete/:id", sm.Route_Delete_Time)

	router.GET("/event", sm.Route_Event)
	router.POST("/event", sm.Route_Event)
	router.GET("/event/:id", sm.Route_Event)
	router.POST("/event/:id", sm.Route_Event)
	router.GET("/event/delete/:id", sm.Route_Delete_Event)

	router.GET("/api/car/:key", sm.API_Car)
	router.GET("/api/car/image/:car/:skin", sm.API_Car_Image)

	router.GET("/api/track/preview/:track/:config", sm.API_Track_Preview_Image)
	router.GET("/api/track/preview/:track", sm.API_Track_Preview_Image)
	router.GET("/api/track/outline/:track/:config", sm.API_Track_Outline_Image)
	router.GET("/api/track/outline/:track", sm.API_Track_Outline_Image)

	router.GET("/api/difficulty/:id", sm.API_Difficulty)
	router.GET("/api/session/:id", sm.API_Session)
	router.GET("/api/class/:id", sm.API_Class)
	router.GET("/api/time/:id", sm.API_Time)

	router.GET("/api/car/recache", sm.API_Recache_Cars)
	router.GET("/api/track/recache", sm.API_Recache_Tracks)
	router.GET("/api/weather/recache", sm.API_Recache_Weathers)

	router.POST("/api/validate/installpath", sm.API_Validate_Installpath)

	router.GET("/api/server/start", sm.API_Console_Start)
	router.GET("/api/server/stop", sm.API_Console_Stop)
	router.GET("/api/server/status", sm.API_Console_Status)

	router.GET("/api/server/entry_list.ini", sm.API_Entry_List)
	router.GET("/api/server/server_cfg.ini", sm.API_Server_Cfg)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "/htm/404.htm", gin.H{})
	})

	router.Run(":3030")



	log.Print(Assets.LocalPath)

}

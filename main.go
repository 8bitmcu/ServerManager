package main

import (
	"html/template"
	"io"
	"log"
	"main/sm"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

var DEBUG bool = true
var dba sm.Dbaccess

func ConfigCompleted() gin.HandlerFunc {
	return func(c *gin.Context) {
		config_filled := dba.Select_Config_Filled()
		if !config_filled && c.Request.URL.Path != "/config" && c.Request.URL.Path != "/content" {
			c.Redirect(http.StatusFound, "/config")
			return
		}
		c.Next()
	}
}

func main() {
	sm.Assets = Assets

	config_folder := os.Getenv("XDG_CONFIG_HOME")
	if runtime.GOOS == "windows" {
		config_folder = os.Getenv("APPDATA")
	}
	sm_path := filepath.Join(config_folder, "servermanager")
	if _, err := os.Stat(sm_path); os.IsNotExist(err) {
		err := os.Mkdir(sm_path, os.ModePerm)
		if err != nil {
			log.Print(err)
		}
	}

	db_path := filepath.Join(sm_path, "smdata.db")
	log.Print("Opening database file located at: " + db_path)
	dba = sm.Open(db_path)
	dba.Apply_Schema(sm.FindFile("/schema.sql"))

	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
	}

	sm.Dba = dba

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	t := template.New("")
	t.Funcs(funcMap)
	err := sm.LoadTemplate(t, ".htm")
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	router.StaticFS("/static", Assets)

	app := router.Group("/")

	app.Use(ConfigCompleted())
	{
		app.GET("/", sm.Route_Index)

		app.GET("/config", sm.Route_Config)
		app.POST("/config", sm.Route_Config)

		app.GET("/content", sm.Route_Content)
		app.POST("/content", sm.Route_Content)

		app.GET("/difficulty", sm.Route_Difficulty)
		app.POST("/difficulty", sm.Route_Difficulty)
		app.GET("/difficulty/:id", sm.Route_Difficulty)
		app.POST("/difficulty/:id", sm.Route_Difficulty)
		app.GET("/difficulty/delete/:id", sm.Route_Delete_Difficulty)
		app.POST("/difficulty/delete/:id", sm.Route_Delete_Difficulty)

		app.GET("/class", sm.Route_Class)
		app.POST("/class", sm.Route_Class)
		app.GET("/class/:id", sm.Route_Class)
		app.POST("/class/:id", sm.Route_Class)
		app.GET("/class/delete/:id", sm.Route_Delete_Class)
		app.POST("/class/delete/:id", sm.Route_Delete_Class)

		app.GET("/session", sm.Route_Session)
		app.POST("/session", sm.Route_Session)
		app.GET("/session/:id", sm.Route_Session)
		app.POST("/session/:id", sm.Route_Session)
		app.GET("/session/delete/:id", sm.Route_Delete_Session)
		app.POST("/session/delete/:id", sm.Route_Delete_Session)

		app.GET("/time", sm.Route_Time)
		app.POST("/time", sm.Route_Time)
		app.GET("/time/:id", sm.Route_Time)
		app.POST("/time/:id", sm.Route_Time)
		app.GET("/time/delete/:id", sm.Route_Delete_Time)
		app.POST("/time/delete/:id", sm.Route_Delete_Time)

		app.GET("/event", sm.Route_Event)
		app.POST("/event", sm.Route_Event)
		app.GET("/event/:id", sm.Route_Event)
		app.POST("/event/:id", sm.Route_Event)
		app.GET("/event/delete/:id", sm.Route_Delete_Event)
		app.POST("/event/delete/:id", sm.Route_Delete_Event)
	}

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
	router.GET("/api/content/recache", sm.API_Recache_Content)

	router.POST("/api/validate/installpath", sm.API_Validate_Installpath)

	router.GET("/api/server/start", sm.API_Console_Start)
	router.GET("/api/server/stop", sm.API_Console_Stop)
	router.GET("/api/server/status", sm.API_Console_Status)

	router.GET("/api/server/entry_list.ini", sm.API_Entry_List)
	router.GET("/api/server/server_cfg.ini", sm.API_Server_Cfg)

	router.NoRoute(sm.NoRoute)

	log.Print("Server up and running on http://localhost:3030")

	if !DEBUG {
		sm.Open_URL("http://localhost:3030")
	}

	go func() {
		res, err := http.Get("https://api.ipify.org")
		if err != nil {
			log.Print(err)
		}
		ip, err := io.ReadAll(res.Body)
		if err != nil {
			log.Print(err)
		}
		sm.Stats.Public_Ip = string(ip)
	}()

	router.Run(":3030")

}

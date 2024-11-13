package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"main/sm"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var dba sm.Dbaccess
var DEBUG bool = true

func ConfigCompletedMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config_filled := dba.Select_Config_Filled()
		if !config_filled && c.Request.URL.Path != "/config" && c.Request.URL.Path != "/content" {
			c.Redirect(http.StatusFound, "/config")
			return
		}

		AuthenticateMiddleware(c)
	}
}

func AuthenticateMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sm.SecretKey, nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	user, err := token.Claims.GetSubject()

	if err != nil {
		log.Print(err)
	}

	c.Set("user", user)
	c.Next()
}

func main() {
	sm.Assets = Assets

	var config_folder string

	flag.StringVar(&config_folder, "p", "", "Configuration path")
	flag.Parse()

	if config_folder == "" {
		config_folder = os.Getenv("XDG_CONFIG_HOME")
		if runtime.GOOS == "windows" {
			config_folder = os.Getenv("APPDATA")
		}
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
	dba.Update_Event_SetComplete()

	gin.SetMode(gin.ReleaseMode)

	sm.Dba = dba
	sm.Cr = sm.ConfigRenderer{}
	sm.Status = sm.Server_Status{}
	sm.Udp = sm.UdpListen()

	go func() {
		for true {
			sm.Udp.Receive()
		}
	}()

	router := gin.New()
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
		"toTime": func(tm int) string {
			tdiff := time.Duration(tm) * time.Second
			return tdiff.Round(time.Second).String()
		},
		"inc": func(i int) int {
			return i + 1
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

	router.GET("/login", sm.Route_Login)
	router.POST("/login", sm.Route_Login)
	router.GET("/logout", sm.Route_Logout)

	app := router.Group("/")
	app.Use(ConfigCompletedMiddlware())
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
		app.GET("/event_edit", sm.Route_Event_Edit)
		app.POST("/event_edit", sm.Route_Event_Edit)
		app.GET("/event_edit/:id", sm.Route_Event_Edit)
		app.POST("/event_edit/:id", sm.Route_Event_Edit)
		app.GET("/event/delete/:id", sm.Route_Delete_Event)
		app.POST("/event/delete/:id", sm.Route_Delete_Event)

		app.GET("/user", sm.Route_User)
		app.POST("/user", sm.Route_User)

		app.GET("/admin", sm.Route_Admin)

		app.GET("/console", sm.Route_Console)
	}

	api := router.Group("/api")
	api.Use(AuthenticateMiddleware)
	{
		api.GET("/car/:key", sm.API_Car)
		api.GET("/car/image/:car/:skin", sm.API_Car_Image)

		api.GET("/track/preview/:track/:config", sm.API_Track_Preview_Image)
		api.GET("/track/preview/:track", sm.API_Track_Preview_Image)
		api.GET("/track/outline/:track/:config", sm.API_Track_Outline_Image)
		api.GET("/track/outline/:track", sm.API_Track_Outline_Image)

		api.GET("/difficulty/:id", sm.API_Difficulty)
		api.GET("/session/:id", sm.API_Session)
		api.GET("/class/:id", sm.API_Class)
		api.GET("/time/:id", sm.API_Time)

		api.GET("/car/recache", sm.API_Recache_Cars)
		api.GET("/track/recache", sm.API_Recache_Tracks)
		api.GET("/weather/recache", sm.API_Recache_Weathers)
		api.GET("/content/recache", sm.API_Recache_Content)

		api.POST("/validate/installpath", sm.API_Validate_Installpath)

		api.GET("/server/start", sm.API_Server_Start)
		api.GET("/server/stop", sm.API_Server_Stop)
		api.GET("/server/status", sm.API_Server_Status)

		api.GET("/server/entry_list.ini", sm.API_Entry_List)
		api.GET("/server/server_cfg.ini", sm.API_Server_Cfg)

	}

	router.NoRoute(sm.NoRoute)

	log.Print("Server up and running on http://localhost:3030")

	if !DEBUG {
		sm.Open_URL("http://localhost:3030")
	}

	sm.Status.Update_Public_Ip()

	srv := &http.Server{
		Addr:    ":3030",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down server...")

	dba.Db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	select {
	case <-ctx.Done():
	}

}

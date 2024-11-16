package main

import (
	"context"
	"flag"
	"html/template"
	"log"
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

var Dba Dbaccess
var Cr ConfigRenderer
var Status Server_Status
var Udp UDPPlugin
var ConfigFolder string
var TempFolder string
var Zf ZipFile

// TODO: checksuming is failing when CSP is enabled
// TODO: generate random SecretKey and storein db
var SecretKey = []byte("XBLn0dUoXPVk742lkRVILa82hbRXz6Tx")
var DEBUG bool = true

func ConfigCompletedMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config_filled := Dba.Select_Config_Filled()
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
		return SecretKey, nil
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

	flag.StringVar(&ConfigFolder, "p", "", "Configuration path")
	flag.Parse()

	if ConfigFolder == "" {
		ConfigFolder = os.Getenv("XDG_CONFIG_HOME")
		if runtime.GOOS == "windows" {
			ConfigFolder = os.Getenv("APPDATA")
		}
	}

	ConfigFolder = filepath.Join(ConfigFolder, "servermanager")
	if _, err := os.Stat(ConfigFolder); os.IsNotExist(err) {
		err := os.Mkdir(ConfigFolder, os.ModePerm)
		if err != nil {
			log.Print(err)
		}
	}

	TempFolder = "/tmp"
	if runtime.GOOS == "windows" {
		TempFolder = os.Getenv("TEMP")
	}
	TempFolder = filepath.Join(TempFolder, "servermanager")
	if _, err := os.Stat(TempFolder); os.IsNotExist(err) {
		err := os.Mkdir(TempFolder, os.ModePerm)
		if err != nil {
			log.Print(err)
		}
	}

	db_path := filepath.Join(ConfigFolder, "smdata.db")
	log.Print("Opening database file located at: " + db_path)
	Dba = Open(db_path)
	Dba.Apply_Schema(FindFile("/schema.sql"))
	Dba.Update_Event_SetComplete()

	gin.SetMode(gin.ReleaseMode)

	Cr = ConfigRenderer{}
	Status = Server_Status{}
	Zf = ZipFile{}
	Udp = UdpListen()

	go func() {
		for true {
			Udp.Receive()
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
	err := LoadTemplate(t, ".htm")
	if err != nil {
		log.Print(err)
	}
	router.SetHTMLTemplate(t)

	router.StaticFS("/static", Assets)

	router.GET("/login", Route_Login)
	router.POST("/login", Route_Login)
	router.GET("/logout", Route_Logout)

	app := router.Group("/")
	app.Use(ConfigCompletedMiddlware())
	{
		app.GET("/", Route_Index)

		app.GET("/config", Route_Config)
		app.POST("/config", Route_Config)

		app.GET("/content", Route_Content)
		app.POST("/content", Route_Content)

		app.GET("/difficulty", Route_Difficulty)
		app.POST("/difficulty", Route_Difficulty)
		app.GET("/difficulty/:id", Route_Difficulty)
		app.POST("/difficulty/:id", Route_Difficulty)
		app.GET("/difficulty/delete/:id", Route_Delete_Difficulty)
		app.POST("/difficulty/delete/:id", Route_Delete_Difficulty)

		app.GET("/class", Route_Class)
		app.POST("/class", Route_Class)
		app.GET("/class/:id", Route_Class)
		app.POST("/class/:id", Route_Class)
		app.GET("/class/delete/:id", Route_Delete_Class)
		app.POST("/class/delete/:id", Route_Delete_Class)

		app.GET("/session", Route_Session)
		app.POST("/session", Route_Session)
		app.GET("/session/:id", Route_Session)
		app.POST("/session/:id", Route_Session)
		app.GET("/session/delete/:id", Route_Delete_Session)
		app.POST("/session/delete/:id", Route_Delete_Session)

		app.GET("/time", Route_Time)
		app.POST("/time", Route_Time)
		app.GET("/time/:id", Route_Time)
		app.POST("/time/:id", Route_Time)
		app.GET("/time/delete/:id", Route_Delete_Time)
		app.POST("/time/delete/:id", Route_Delete_Time)

		app.GET("/event_cat", Route_Event_Category)
		app.POST("/event_cat", Route_Event_Category)
		app.GET("/event_cat/:id", Route_Event_Category)
		app.POST("/event_cat/:id", Route_Event_Category)
		app.GET("/event_cat/delete/:id", Route_Delete_Event_Category)
		app.POST("/event_cat/delete/:id", Route_Delete_Event_Category)
		app.GET("/event/:category_id", Route_Event)
		app.POST("/event/:category_id", Route_Event)
		app.GET("/event/:category_id/:id", Route_Event)
		app.POST("/event/:category_id/:id", Route_Event)
		app.GET("/event/delete/:id", Route_Delete_Event)
		app.POST("/event/delete/:id", Route_Delete_Event)

		app.GET("/user", Route_User)
		app.POST("/user", Route_User)

		app.GET("/admin", Route_Admin)

		app.GET("/console", Route_Console)
	}

	api := router.Group("/api")
	api.Use(AuthenticateMiddleware)
	{
		api.GET("/car/:key", API_Car)
		api.GET("/car/image/:car/:skin", API_Car_Image)

		api.GET("/track/preview/:track/:config", API_Track_Preview_Image)
		api.GET("/track/preview/:track", API_Track_Preview_Image)
		api.GET("/track/outline/:track/:config", API_Track_Outline_Image)
		api.GET("/track/outline/:track", API_Track_Outline_Image)

		api.GET("/difficulty/:id", API_Difficulty)
		api.GET("/session/:id", API_Session)
		api.GET("/class/:id", API_Class)
		api.GET("/time/:id", API_Time)

		api.GET("/content/recache", API_Recache_Content)

		api.POST("/validate/installpath", API_Validate_Installpath)

		api.GET("/server/start", API_Server_Start)
		api.GET("/server/stop", API_Server_Stop)
		api.GET("/server/status", API_Server_Status)

		api.GET("/server/entry_list.ini", API_Entry_List)
		api.GET("/server/server_cfg.ini", API_Server_Cfg)

	}

	router.NoRoute(NoRoute)

	if !DEBUG {
		Open_URL("http://localhost:3030")
	}

	Status.Update_Public_Ip()

	main := &http.Server{
		Addr:    ":3030",
		Handler: router.Handler(),
	}

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	utilities := &http.Server{
		Addr:    ":4040",
		Handler: r.Handler(),
	}

	go func() {
		log.Print("Server up and running on http://localhost:3030")
		if err := main.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print(err)
		}
	}()
	go func() {
		log.Print("Utilities up and running on http://localhost:4040")
		if err := utilities.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down server...")

	Dba.Db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := main.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	select {
	case <-ctx.Done():
	}

}

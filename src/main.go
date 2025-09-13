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
var Status ServerStatus
var Udp UdpPlugin
var ConfigFolder string
var TempFolder string
var SecretKey []byte
var Zf ZipFile

// TODO: checksuming is failing when CSP is enabled
var debug bool = true

func ConfigCompletedMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		configfilled := Dba.selectConfigFilled()
		if !configfilled && c.Request.URL.Path != "/config" && c.Request.URL.Path != "/content" {
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

	dbpath := filepath.Join(ConfigFolder, "smdata.db")
	log.Print("Opening database file located at: " + dbpath)
	Dba = open(dbpath)
	Dba.applySchema(FindFile("/schema.sql"))
	SecretKey = []byte(*Dba.selectConfig().SecretKey)

	gin.SetMode(gin.ReleaseMode)

	Cr = ConfigRenderer{}
	Status = ServerStatus{}
	Zf = ZipFile{}
	Udp = udpListen()

	go func() {
		for {
			Udp.Receive()
		}
	}()

	var router *gin.Engine
	if debug {
		router = gin.Default()
	} else {
		router = gin.New()
	}

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
		"derefInt64": func(t *int64) string {
			if t == nil {
				return ""
			}
			return strconv.FormatInt(*t, 10)
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

	router.GET("/login", routeLogin)
	router.POST("/login", routeLogin)
	router.GET("/logout", routeLogout)

	app := router.Group("/")
	app.Use(ConfigCompletedMiddlware())
	{
		app.GET("/", routeIndex)

		app.GET("/config", routeConfig)
		app.POST("/config", routeConfig)

		app.GET("/content", routeContent)
		app.POST("/content", routeContent)

		app.GET("/difficulty", routeDifficulty)
		app.POST("/difficulty", routeDifficulty)
		app.GET("/difficulty/:id", routeDifficulty)
		app.POST("/difficulty/:id", routeDifficulty)
		app.GET("/difficulty/delete/:id", routeDeleteDifficulty)
		app.POST("/difficulty/delete/:id", routeDeleteDifficulty)

		app.GET("/class", routeClass)
		app.POST("/class", routeClass)
		app.GET("/class/:id", routeClass)
		app.POST("/class/:id", routeClass)
		app.GET("/class/delete/:id", routeDeleteClass)
		app.POST("/class/delete/:id", routeDeleteClass)

		app.GET("/session", routeSession)
		app.POST("/session", routeSession)
		app.GET("/session/:id", routeSession)
		app.POST("/session/:id", routeSession)
		app.GET("/session/delete/:id", routeDeleteSession)
		app.POST("/session/delete/:id", routeDeleteSession)

		app.GET("/time", routeTime)
		app.POST("/time", routeTime)
		app.GET("/time/:id", routeTime)
		app.POST("/time/:id", routeTime)
		app.GET("/time/delete/:id", routeDeleteTime)
		app.POST("/time/delete/:id", routeDeleteTime)

		app.GET("/event_cat", routeEventCategory)
		app.POST("/event_cat", routeEventCategory)
		app.GET("/event_cat/:id", routeEventCategory)
		app.POST("/event_cat/:id", routeEventCategory)
		app.GET("/event_cat/delete/:id", routeDeleteEventCategory)
		app.POST("/event_cat/delete/:id", routeDeleteEventCategory)

		app.GET("/event/:category_id", routeEvent)
		app.POST("/event/:category_id", routeEvent)
		app.GET("/event/:category_id/:id", routeEvent)
		app.POST("/event/:category_id/:id", routeEvent)
		app.GET("/event/delete/:id", routeDeleteEvent)
		app.POST("/event/delete/:id", routeDeleteEvent)

		app.GET("/queue", routeQueue)
		app.POST("/queue", routeQueue)
		app.GET("/queue/delete/:id", routeDeleteQueue)
		app.POST("/queue/delete/:id", routeDeleteQueue)

		app.GET("/user", routeUser)
		app.POST("/user", routeUser)

		app.GET("/admin", routeAdmin)

		app.GET("/server", routeServer)
	}

	api := router.Group("/api")
	api.Use(AuthenticateMiddleware)
	{
		api.GET("/car/:key", apiCar)
		api.GET("/car/image/:car/:skin", apiCarImage)

		api.GET("/track/preview/:track/:config", apiTrackPreviewImage)
		api.GET("/track/preview/:track", apiTrackPreviewImage)
		api.GET("/track/outline/:track/:config", apiTrackOutlineImage)
		api.GET("/track/outline/:track", apiTrackOutlineImage)

		api.GET("/difficulty/:id", apiDifficulty)
		api.GET("/session/:id", apiSession)
		api.GET("/class/:id", apiClass)
		api.GET("/time/:id", apiTime)

		api.GET("/content/recache", apiRecacheContent)

		api.POST("/validate/installpath", apiValidateInstallpath)

		api.GET("/server/start", apiServerStart)
		api.GET("/server/stop", apiServerStop)
		api.GET("/server/status", apiServerStatus)

		api.GET("/server/entry_list.ini", apiEntryList)
		api.GET("/server/server_cfg.ini", apiServerCfg)

		api.GET("/queue/moveup/:id", apiQueueMoveUp)
		api.GET("/queue/movedown/:id", apiQueueMoveDown)
		api.GET("/queue/skipevent", apiQueueSkipEvent)
		api.GET("/queue/clearcompleted", apiQueueClearCompleted)
	}

	router.NoRoute(noRoute)

	if !debug {
		OpenURL("http://localhost:3030")
	}

	Status.updatePublicIp()

	main := &http.Server{
		Addr:    ":3030",
		Handler: router.Handler(),
	}

	go func() {
		log.Print("Server up and running on http://localhost:3030")
		if err := main.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down server...")

	Dba.db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := main.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	select {
	case <-ctx.Done():
	}

}

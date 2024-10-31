package main

import (
	"fmt"
	//"os"
	"log"
	"net/http"
	//"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html/template"
	//"github.com/kaptinlin/jsonrepair"
	"main/sm"
	"strconv"
)

var dba sm.Dbaccess

func handle_config(c *gin.Context) {
	var form sm.User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		dba.Update_Config(form)
	}

	c.HTML(http.StatusOK, "config.htm", gin.H{
		"page": "config",
		"form": dba.Select_Config(),
	})
}

func handle_content(c *gin.Context) {
	var form sm.User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		dba.Update_Content(form)
	}

	c.HTML(http.StatusOK, "content.htm", gin.H{
		"page": "content",
		"form": dba.Select_Config(),
	})
}

func handle_index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.htm", gin.H{
		"page":       "index",
		"is_running": false,
	})
}

func handle_difficulty(c *gin.Context) {
	form := sm.User_Difficulty{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = dba.Select_Difficulty(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("difficulty_name") != "" {
			id := dba.Insert_Difficulty(c.PostForm("difficulty_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/difficulty/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			dba.Update_Difficulty(form)
		}
	}
	c.HTML(http.StatusOK, "difficulty.htm", gin.H{
		"page": "difficulty",
		"list": dba.Select_DifficultyList(false),
		"form": form,
	})
}

func handle_class(c *gin.Context) {
	form := sm.User_Class{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = dba.Select_Class(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("class_name") != "" {
			dba.Insert_Class(c.PostForm("class_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/class/", id))
		} else if c.ShouldBind(&form) == nil {
			dba.Update_Class(form)
		}
	}
	c.HTML(http.StatusOK, "class.htm", gin.H{
		"page":     "class",
		"list":     dba.Select_ClassList(false),
		"car_data": []sm.Cache_Car{},
		"form":     form,
	})
}

func handle_session(c *gin.Context) {
	form := sm.User_Session{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = dba.Select_Session(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("session_name") != "" {
			dba.Insert_Session(c.PostForm("session_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/session/", id))
		} else if c.ShouldBind(&form) == nil {
			log.Print("Updating session")
			dba.Update_Session(form)
		}
	}
	c.HTML(http.StatusOK, "session.htm", gin.H{
		"page": "session",
		"list": dba.Select_SessionList(false),
		"form": form,
	})
}

func handle_time(c *gin.Context) {
	form := sm.User_Time{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = dba.Select_Time_Weather(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("time_name") != "" {
			id := dba.Insert_Time(c.PostForm("time_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/time/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("weather")), &form.Weathers)
			dba.Update_Time(form)
		}
	}

	if len(form.Weathers) == 0 {
		form.Weathers = append(form.Weathers, sm.User_Time_Weather{})
	}

	c.HTML(http.StatusOK, "time.htm", gin.H{
		"page":        "time",
		"list":        dba.Select_TimeList(false),
		"weatherlist": dba.Select_Cache_Weathers(),
		"form":        form,
	})
}

func handle_event(c *gin.Context) {
	c.HTML(http.StatusOK, "time.htm", gin.H{
		"page": "event",
		"list": dba.Select_Events(),
	})

}

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
		"toJsBool": func(t *int) bool {
			if t == nil {
				return true
			}
			if *t > 0 {
				return true
			}
			return false
		},
		"toJson": func(t interface{}) string {
			b, _ := json.Marshal(t)
			return string(b)
		},
	})

	router.LoadHTMLGlob("htm/*")

	//sm.Parse_Weathers(dba)
	//sm.Parse_Tracks(dba)
	//sm.Parse_Cars(dba)

	router.GET("/", handle_index)

	router.GET("/config", handle_config)
	router.POST("/config", handle_config)

	router.GET("/content", handle_content)
	router.POST("/content", handle_content)

	router.GET("/difficulty", handle_difficulty)
	router.POST("/difficulty", handle_difficulty)
	router.GET("/difficulty/:id", handle_difficulty)
	router.POST("/difficulty/:id", handle_difficulty)

	router.GET("/class", handle_class)
	router.POST("/class", handle_class)
	router.GET("/class/:id", handle_class)
	router.POST("/class/:id", handle_class)

	router.GET("/session", handle_session)
	router.POST("/session", handle_session)
	router.GET("/session/:id", handle_session)
	router.POST("/session/:id", handle_session)

	router.GET("/time", handle_time)
	router.POST("/time", handle_time)
	router.GET("/time/:id", handle_time)
	router.POST("/time/:id", handle_time)

	router.GET("/event", handle_event)
	router.POST("/event", handle_event)
	router.GET("/event/:id", handle_event)
	router.POST("/event/:id", handle_event)

	router.Run(":3030")
}

package sm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Dba Dbaccess

func Route_Config(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Config(form)
	}

	c.HTML(http.StatusOK, "config.htm", gin.H{
		"page": "config",
		"form": Dba.Select_Config(),
	})
}

func Route_Content(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Content(form)
	}

	c.HTML(http.StatusOK, "content.htm", gin.H{
		"page": "content",
		"form": Dba.Select_Config(),
	})
}

func Route_Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.htm", gin.H{
		"page":       "index",
		"is_running": false,
	})
}

func Route_Difficulty(c *gin.Context) {
	form := User_Difficulty{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Difficulty(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("difficulty_name") != "" {
			id := Dba.Insert_Difficulty(c.PostForm("difficulty_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/difficulty/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			Dba.Update_Difficulty(form)
		}
	}
	c.HTML(http.StatusOK, "difficulty.htm", gin.H{
		"page": "difficulty",
		"list": Dba.Select_DifficultyList(false),
		"form": form,
	})
}

func Route_Class(c *gin.Context) {
	form := User_Class{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Class_Entries(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("class_name") != "" {
			id := Dba.Insert_Class(c.PostForm("class_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/class/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("entries")), &form.Entries)
			Dba.Update_Class(form)
		}
	}
	c.HTML(http.StatusOK, "class.htm", gin.H{
		"page":     "class",
		"list":     Dba.Select_ClassList(false),
		"car_data": Dba.Select_Cache_Cars(), // TODO: only select needed data
		"form":     form,
	})
}

func Route_Session(c *gin.Context) {
	form := User_Session{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Session(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("session_name") != "" {
			id := Dba.Insert_Session(c.PostForm("session_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/session/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			log.Print("Updating session")
			Dba.Update_Session(form)
		}
	}
	c.HTML(http.StatusOK, "session.htm", gin.H{
		"page": "session",
		"list": Dba.Select_SessionList(false),
		"form": form,
	})
}

func Route_Time(c *gin.Context) {
	form := User_Time{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Time_Weather(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("time_name") != "" {
			id := Dba.Insert_Time(c.PostForm("time_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/time/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("weather")), &form.Weathers)
			Dba.Update_Time(form)
		}
	}

	if len(form.Weathers) == 0 {
		form.Weathers = append(form.Weathers, User_Time_Weather{})
	}

	c.HTML(http.StatusOK, "time.htm", gin.H{
		"page":        "time",
		"list":        Dba.Select_TimeList(false),
		"weatherlist": Dba.Select_Cache_Weathers(),
		"form":        form,
	})
}

func Route_Event(c *gin.Context) {
	form := User_Event{}

	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Event(id)

		if form.Id == nil {
			c.HTML(http.StatusNotFound, "404.htm", gin.H{})
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("id") != "" {
			log.Print("updating event")
			Dba.Update_Event(form)
		} else if c.ShouldBind(&form) == nil {
			log.Print("insertin event")
			id := Dba.Insert_Event(form)
			c.Redirect(http.StatusFound, fmt.Sprint("/event/", id))
			return
		}
	}

	c.HTML(http.StatusOK, "event.htm", gin.H{
		"page":         "event",
		"form":         form,
		"events":       Dba.Select_Events(),
		"difficulties": Dba.Select_DifficultyList(true),
		"sessions":     Dba.Select_SessionList(true),
		"times":        Dba.Select_TimeList(true),
		"classes":      Dba.Select_ClassList(true),
		"max_clients":  2, // TODO
		"track_data":   Dba.Select_Cache_Tracks(),
	})

}

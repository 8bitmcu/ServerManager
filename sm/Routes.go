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
var Stats Server_Stats

func Route_Config(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Config(form)
	}

	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/config.htm", gin.H{
		"page":          "config",
		"form":          Dba.Select_Config(),
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Content(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Content(form)
	}

	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/content.htm", gin.H{
		"page":          "content",
		"form":          Dba.Select_Config(),
		"track_data":    Dba.Select_Cache_Tracks(), // TODO: we only need a list at this point
		"car_data":      Dba.Select_Cache_Cars(),
		"weather_data":  Dba.Select_Cache_Weathers(),
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Index(c *gin.Context) {
	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/index.htm", gin.H{
		"page":          "index",
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Difficulty(c *gin.Context) {
	form := User_Difficulty{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Difficulty(id)

		if form.Id == nil {
			NoRoute(c)
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
	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/difficulty.htm", gin.H{
		"page":          "difficulty",
		"list":          Dba.Select_DifficultyList(false),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Delete_Difficulty(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Dba.Delete_Difficulty(id)
	c.Redirect(http.StatusFound, "/difficulty")
	return
}

func Route_Class(c *gin.Context) {
	form := User_Class{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Class_Entries(id)

		if form.Id == nil {
			NoRoute(c)
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
	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/class.htm", gin.H{
		"page":          "class",
		"list":          Dba.Select_ClassList(false),
		"car_data":      Dba.Select_Cache_Cars(), // TODO: only select needed data
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Delete_Class(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Dba.Delete_Class(id)
	c.Redirect(http.StatusFound, "/class")
	return
}

func Route_Session(c *gin.Context) {
	form := User_Session{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Session(id)

		if form.Id == nil {
			NoRoute(c)
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
	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/session.htm", gin.H{
		"page":          "session",
		"list":          Dba.Select_SessionList(false),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Delete_Session(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Dba.Delete_Session(id)
	c.Redirect(http.StatusFound, "/session")
	return
}

func Route_Time(c *gin.Context) {
	form := User_Time{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Time_Weather(id)

		if form.Id == nil {
			NoRoute(c)
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

	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/time.htm", gin.H{
		"page":          "time",
		"list":          Dba.Select_TimeList(false),
		"weatherlist":   Dba.Select_Cache_Weathers(),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

func Route_Delete_Time(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Dba.Delete_Time(id)
	c.Redirect(http.StatusFound, "/time")
	return
}

func Route_Event(c *gin.Context) {
	form := User_Event{}

	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Event(id)

		if form.Id == nil {
			NoRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("id") != "" {
			Dba.Update_Event(form)
		} else if c.ShouldBind(&form) == nil {
			id := Dba.Insert_Event(form)
			c.Redirect(http.StatusFound, fmt.Sprint("/event/", id))
			return
		}
	}

	Stats.Status = Is_Running()
	c.HTML(http.StatusOK, "/htm/event.htm", gin.H{
		"page":          "event",
		"form":          form,
		"events":        Dba.Select_Events(),
		"difficulties":  Dba.Select_DifficultyList(true),
		"sessions":      Dba.Select_SessionList(true),
		"times":         Dba.Select_TimeList(true),
		"classes":       Dba.Select_ClassList(true),
		"max_clients":   2, // TODO
		"track_data":    Dba.Select_Cache_Tracks(),
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})

}

func Route_Delete_Event(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Dba.Delete_Event(id)
	c.Redirect(http.StatusFound, "/event")
	return
}

func NoRoute(c *gin.Context) {
	Stats.Status = Is_Running()
	c.HTML(http.StatusNotFound, "/htm/404.htm", gin.H{
		"config_filled": Dba.Select_Config_Filled(),
		"stats":         Stats,
	})
}

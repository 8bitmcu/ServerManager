package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Route_Config(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Config(form)
	}

	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/config.htm", gin.H{
		"page":          "config",
		"form":          Dba.Select_Config(),
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Content(c *gin.Context) {
	var form User_Config
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_Content(form)
	}

	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/content.htm", gin.H{
		"page":          "content",
		"form":          Dba.Select_Config(),
		"track_data":    Dba.Select_Cache_Tracks(),
		"car_data":      Dba.Select_Cache_Cars(),
		"weather_data":  Dba.Select_Cache_Weathers(),
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/server")
	return
}

func Route_Server(c *gin.Context) {
	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/server.htm", gin.H{
		"page":          "server",
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Queue(c *gin.Context) {
	if c.Request.Method == "POST" {
		if c.PostForm("event") != "" {
			// insert single event from category
			id, err := strconv.Atoi(c.PostForm("event"))
			if id > 0 && err == nil {
				Dba.Insert_Server_Event(id)
			}
		} else {
			// insert all events from category
			id, err := strconv.Atoi(c.PostForm("category"))
			if id > 0 && err == nil {
				Dba.Insert_Server_Event_Category(id)
			}
		}
	}
	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/queue.htm", gin.H{
		"page":          "queue",
		"config_filled": Dba.Select_Config_Filled(),
		"event_cat":     Dba.Select_Event_CategoryList(false),
		"event_list":    Dba.Select_EventList(),
		"server_events": Dba.Select_Server_Events(false),
		"status":        Status,
	})
}

func Route_Delete_Queue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Server_Event(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/queue")
	}
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
	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/difficulty.htm", gin.H{
		"page":          "difficulty",
		"list":          Dba.Select_DifficultyList(false),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Delete_Difficulty(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Difficulty(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/difficulty")
	}
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
	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/class.htm", gin.H{
		"page":          "class",
		"list":          Dba.Select_ClassList(false),
		"car_data":      Dba.Select_Cache_Cars(),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Delete_Class(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Class(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/class")
	}
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
			Dba.Update_Session(form)
		}
	}
	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/session.htm", gin.H{
		"page":          "session",
		"list":          Dba.Select_SessionList(false),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Delete_Session(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Session(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/session")
	}
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

	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/time.htm", gin.H{
		"page":          "time",
		"list":          Dba.Select_TimeList(false),
		"weatherlist":   Dba.Select_Cache_Weathers(),
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Delete_Time(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Time(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/time")
	}
}

func Route_Event_Category(c *gin.Context) {
	form := User_Event_Category{}
	Status.Refresh()

	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Category_Events(id)

		if form.Id == nil {
			NoRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("category_name") != "" {
			id := Dba.Insert_Event_Category(c.PostForm("category_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/event_cat/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			Dba.Update_Event_Category(form)
		}
	}

	max_clients := Dba.Select_Config().Max_Clients

	c.HTML(http.StatusOK, "/htm/event_cat.htm", gin.H{
		"page":          "event_cat",
		"list":          Dba.Select_Event_CategoryList(false),
		"form":          form,
		"track_data":    Dba.Select_Cache_Tracks(),
		"config_filled": Dba.Select_Config_Filled(),
		"max_clients":   max_clients,
		"status":        Status,
	})
}

func Route_Delete_Event_Category(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.Delete_Event_Category(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/event_cat")
	}
}

func Route_Event(c *gin.Context) {
	form := User_Event{}

	category_id, err := strconv.Atoi(c.Param("category_id"))
	if category_id <= 0 || err != nil {
		NoRoute(c)
		return
	}
	form.Event_Category_Id = &category_id

	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.Select_Event(id)

		if form.Id == nil {
			NoRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		if c.PostForm("id") != "" {
			Dba.Update_Event(form)
			c.Redirect(http.StatusFound, "/event_cat/"+strconv.Itoa(category_id))
			return
		} else {
			Dba.Insert_Event(form)
			c.Redirect(http.StatusFound, "/event_cat/"+strconv.Itoa(category_id))
			return
		}
	}

	Status.Refresh()
	c.HTML(http.StatusOK, "/htm/event.htm", gin.H{
		"page":          "event",
		"form":          form,
		"difficulties":  Dba.Select_DifficultyList(true),
		"sessions":      Dba.Select_SessionList(true),
		"times":         Dba.Select_TimeList(true),
		"classes":       Dba.Select_ClassList(true),
		"max_clients":   Dba.Select_Config().Max_Clients,
		"track_data":    Dba.Select_Cache_Tracks(),
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Delete_Event(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	evt := Dba.Select_Event(id)
	cat_id := *evt.Event_Category_Id
	_, err := Dba.Delete_Event(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.Select_Config_Filled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/event_cat"+strconv.Itoa(cat_id))
	}
}

func Route_Login(c *gin.Context) {
	if c.Request.Method == "POST" {
		usr := c.PostForm("name")
		pwd := c.PostForm("password")

		user := Dba.Select_User(usr)

		if user.Name == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		res := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(pwd))

		if res == nil {
			claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.Name,
				"iss": "servermanager",
				"aud": "admin",
				"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
				"iat": time.Now().Unix(),
			})

			tokenString, err := claims.SignedString(SecretKey)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error creating token")
				return
			}
			c.SetCookie("token", tokenString, 3600*24*30, "/", "", false, true)

			log.Print("Login successful for user: " + *user.Name)

			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	c.HTML(http.StatusOK, "/htm/login.htm", gin.H{})
}

func Route_Logout(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
	return
}

func Route_User(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	form := Dba.Select_User(user.(string))
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.Update_User(form)
	}

	c.HTML(http.StatusOK, "/htm/user.htm", gin.H{
		"page":          "user",
		"form":          form,
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

func Route_Admin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if user != "admin" {
		NoRoute(c)
		return
	}

	c.HTML(http.StatusOK, "/htm/admin.htm", gin.H{
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})

}

func NoRoute(c *gin.Context) {
	Status.Refresh()
	c.HTML(http.StatusNotFound, "/htm/404.htm", gin.H{
		"config_filled": Dba.Select_Config_Filled(),
		"status":        Status,
	})
}

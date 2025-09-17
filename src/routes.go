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

func routeConfig(c *gin.Context) {
	var form UserConfig
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.updateConfig(form)
	}

	Status.refresh()
	c.HTML(http.StatusOK, "/htm/config.htm", gin.H{
		"page":          "config",
		"form":          Dba.selectConfig(),
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeContent(c *gin.Context) {
	var form UserConfig
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.updateContent(form)
	}

	Status.refresh()
	c.HTML(http.StatusOK, "/htm/content.htm", gin.H{
		"page":          "content",
		"form":          Dba.selectConfig(),
		"track_data":    Dba.selectCacheTracks(),
		"car_data":      Dba.selectCacheCars(),
		"weather_data":  Dba.selectCacheWeathers(),
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/server")
}

func routeServer(c *gin.Context) {
	Status.refresh()
	c.HTML(http.StatusOK, "/htm/server.htm", gin.H{
		"page":          "server",
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeQueue(c *gin.Context) {
	if c.Request.Method == "POST" {
		if c.PostForm("event") != "" {
			// insert single event from category
			id, err := strconv.Atoi(c.PostForm("event"))
			if id > 0 && err == nil {
				Dba.insertServerEvent(id)
			}
		} else {
			// insert all events from category
			id, err := strconv.Atoi(c.PostForm("category"))
			if id > 0 && err == nil {
				Dba.insertServerEventCategory(id)
			}
		}
	}
	Status.refresh()
	c.HTML(http.StatusOK, "/htm/queue.htm", gin.H{
		"page":          "queue",
		"config_filled": Dba.selectConfigFilled(),
		"event_cat":     Dba.selectEventCategoryList(false),
		"event_list":    Dba.selectEventList(),
		"server_events": Dba.selectServerEvents(false),
		"status":        Status,
	})
}

func routeDeleteQueue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteServerEvent(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/queue")
	}
}

func routeDifficulty(c *gin.Context) {
	form := UserDifficulty{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.selectDifficulty(id)

		if form.Id == nil {
			noRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("difficulty_name") != "" {
			id := Dba.insertDifficulty(c.PostForm("difficulty_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/difficulty/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			Dba.updateDifficulty(form)
		}
	}
	Status.refresh()
	c.HTML(http.StatusOK, "/htm/difficulty.htm", gin.H{
		"page":          "difficulty",
		"list":          Dba.selectDifficultyList(false),
		"form":          form,
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeDeleteDifficulty(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteDifficulty(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/difficulty")
	}
}

func routeClass(c *gin.Context) {
	form := UserClass{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.selectClassEntries(id)

		if form.Id == nil {
			noRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("class_name") != "" {
			id := Dba.insertClass(c.PostForm("class_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/class/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("entries")), &form.Entries)
			Dba.updateClass(form)
		}
	}
	Status.refresh()
	c.HTML(http.StatusOK, "/htm/class.htm", gin.H{
		"page":          "class",
		"list":          Dba.selectClassList(false),
		"car_data":      Dba.selectCacheCars(),
		"form":          form,
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeDeleteClass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteClass(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/class")
	}
}

func routeSession(c *gin.Context) {
	form := UserSession{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.selectSession(id)

		if form.Id == nil {
			noRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("session_name") != "" {
			id := Dba.insertSession(c.PostForm("session_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/session/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			Dba.updateSession(form)
		}
	}
	Status.refresh()
	c.HTML(http.StatusOK, "/htm/session.htm", gin.H{
		"page":          "session",
		"list":          Dba.selectSessionList(false),
		"form":          form,
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeDeleteSession(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteSession(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/session")
	}
}

func routeTime(c *gin.Context) {
	form := UserTime{}
	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.selectTimeWeather(id)

		if form.Id == nil {
			noRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("time_name") != "" {
			id := Dba.insertTime(c.PostForm("time_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/time/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("weather")), &form.Weathers)
			Dba.updateTime(form)
		}
	}

	if len(form.Weathers) == 0 {
		form.Weathers = append(form.Weathers, UserTimeWeather{})
	}

	Status.refresh()
	c.HTML(http.StatusOK, "/htm/time.htm", gin.H{
		"page":          "time",
		"list":          Dba.selectTimeList(false),
		"weatherlist":   Dba.selectCacheWeathers(),
		"form":          form,
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeDeleteTime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteTime(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/time")
	}
}

func routeEventCategory(c *gin.Context) {
	form := UserEventCategory{}
	Status.refresh()

	id, err := strconv.Atoi(c.Param("id"))
	if id > 0 && err == nil {
		form = Dba.selectCategoryEvents(id)

		if form.Id == nil {
			noRoute(c)
			return
		}
	}
	if c.Request.Method == "POST" {
		if c.PostForm("category_name") != "" {
			id := Dba.insertEventCategory(c.PostForm("category_name"))
			c.Redirect(http.StatusFound, fmt.Sprint("/event/", id))
			return
		} else if c.ShouldBind(&form) == nil {
			json.Unmarshal([]byte(c.PostForm("events")), &form.Events)
			Dba.updateEventCategory(form)
		}
	}

	if len(form.Events) == 0 {
		form.Events = append(form.Events, UserEvent{})
	}

	c.HTML(http.StatusOK, "/htm/event.htm", gin.H{
		"page":          "event",
		"list":          Dba.selectEventCategoryList(false),
		"form":          form,
		"difficulties":  Dba.selectDifficultyList(true),
		"sessions":      Dba.selectSessionList(true),
		"times":         Dba.selectTimeList(true),
		"classes":       Dba.selectClassList(true),
		"max_clients":   Dba.selectConfig().MaxClients,
		"track_data":    Dba.selectCacheTracks(),
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeDeleteEventCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := Dba.deleteEventCategory(id)

	if err != nil {
		c.HTML(http.StatusOK, "/htm/dberror.htm", gin.H{
			"status":        Status,
			"config_filled": Dba.selectConfigFilled(),
			"error":         err.Error(),
		})
	} else {
		c.Redirect(http.StatusFound, "/event")
	}
}

func routeLogin(c *gin.Context) {
	if c.Request.Method == "POST" {
		usr := c.PostForm("name")
		pwd := c.PostForm("password")

		user := Dba.selectUser(usr)

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

func routeLogout(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func routeUser(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	form := Dba.selectUser(user.(string))
	if c.Request.Method == "POST" && c.ShouldBind(&form) == nil {
		Dba.updateUser(form)
	}

	c.HTML(http.StatusOK, "/htm/user.htm", gin.H{
		"page":          "user",
		"form":          form,
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

func routeAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if user != "admin" {
		noRoute(c)
		return
	}

	c.HTML(http.StatusOK, "/htm/admin.htm", gin.H{
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})

}

func noRoute(c *gin.Context) {
	Status.refresh()
	c.HTML(http.StatusNotFound, "/htm/404.htm", gin.H{
		"config_filled": Dba.selectConfigFilled(),
		"status":        Status,
	})
}

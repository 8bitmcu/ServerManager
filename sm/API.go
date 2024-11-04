package sm

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func API_Car_Image(c *gin.Context) {
	car := c.Param("car")
	skin := c.Param("skin")

	file := filepath.Join(Dba.Basepath(), "content", "cars", car, "skins", skin, "preview.jpg")
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.FileAttachment(file, car+"_"+skin+".jpg")
		return
	}
}

func API_Car(c *gin.Context) {
	key := c.Param("key")

	car_data := Dba.Select_Cache_Car(key)

	var power []any
	var torque []any
	var labels []any

	for _, l := range car_data.Power {
		labels = append(labels, l[0])
	}
	for _, p := range car_data.Power {
		power = append(power, p[1])
	}
	for _, t := range car_data.Torque {
		torque = append(torque, t[1])
	}

	c.PureJSON(http.StatusOK, gin.H{
		"key":    key,
		"desc":   car_data.Desc,
		"power":  power,
		"torque": torque,
		"labels": labels,
	})
}

func API_Track_Preview_Image(c *gin.Context) {
	track := c.Param("track")
	config := c.Param("config")

	file := ""
	fileName := ""
	if config != "" {
		file = filepath.Join(Dba.Basepath(), "content", "tracks", track, "ui", config, "preview.png")
		fileName = track + "_" + config
	} else {
		file = filepath.Join(Dba.Basepath(), "content", "tracks", track, "ui", "preview.png")
		fileName = track
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.FileAttachment(file, fileName+".png")
		return
	}
}

func API_Track_Outline_Image(c *gin.Context) {
	track := c.Param("track")
	config := c.Param("config")

	file := ""
	fileName := ""
	if config != "" {
		file = filepath.Join(Dba.Basepath(), "content", "tracks", track, "ui", config, "outline.png")
		fileName = track + "_" + config
	} else {
		file = filepath.Join(Dba.Basepath(), "content", "tracks", track, "ui", "outline.png")
		fileName = track
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.FileAttachment(file, fileName+".png")
		return
	}
}

func API_Difficulty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}

	data := Dba.Select_Difficulty(id)
	if data.Id == nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Session(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}

	data := Dba.Select_Session(id)
	if data.Id == nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Class(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}

	data := Dba.Select_Class_Entries(id)
	if data.Id == nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Time(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}

	data := Dba.Select_Time_Weather(id)
	if data.Id == nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Recache_Cars(c *gin.Context) {
	result := Parse_Cars(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
		"value":  result,
	})
}

func API_Recache_Tracks(c *gin.Context) {
	result := Parse_Tracks(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
		"value":  result,
	})
}

func API_Recache_Weathers(c *gin.Context) {
	result := Parse_Weathers(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
		"value":  result,
	})
}
func API_Recache_Content(c *gin.Context) {
	track := Parse_Tracks(Dba)
	car := Parse_Cars(Dba)
	weather := Parse_Weathers(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result":         "ok",
		"tracks_total":   track,
		"cars_total":     car,
		"weathers_total": weather,
	})
}

func API_Validate_Installpath(c *gin.Context) {
	// TODO eventually check for server binary, based on OS
	path := filepath.Join(c.PostForm("path"), "acs.exe")
	exists := true
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		exists = false
	}

	c.PureJSON(http.StatusOK, gin.H{
		"result": exists,
	})
}

func API_Console_Start(c *gin.Context) {
	Start()
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

func API_Console_Stop(c *gin.Context) {
	Stop()
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

func API_Console_Status(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

// Warning: non determistic as some of the items can be randomized/shuffled.
func API_Entry_List(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	res := Cr.Render_Ini(idInt)
	c.String(http.StatusOK, res.EntryList_Result)
}

func API_Server_Cfg(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	res := Cr.Render_Ini(idInt)
	c.String(http.StatusOK, res.ServerCfg_Result)
}

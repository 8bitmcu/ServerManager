package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func API_Car_Image(c *gin.Context) {
	car := c.Param("car")
	skin := c.Param("skin")

	var zf ZipFile
	zi := zf.FindZipFile("cars/" + car + "/skins/" + skin + "/preview.jpg")
	if zi != nil {
		r, err := zi.Open()
		if err != nil {
			log.Print(err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/jpg", r, nil)
	} else {
		NoRoute(c)
	}
	zf.Close()
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

	filePath := "tracks/" + track + "/preview.png"
	if config != "" {
		filePath = "tracks/" + track + "/" + config + "/preview.png"
	}

	var zf ZipFile
	zi := zf.FindZipFile(filePath)
	if zi != nil {
		r, err := zi.Open()
		if err != nil {
			log.Print(err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/png", r, nil)
	} else {
		NoRoute(c)
	}
	zf.Close()
}

func API_Track_Outline_Image(c *gin.Context) {
	track := c.Param("track")
	config := c.Param("config")

	filePath := "tracks/" + track + "/outline.png"
	if config != "" {
		filePath = "tracks/" + track + "/" + config + "/outline.png"
	}

	var zf ZipFile
	zi := zf.FindZipFile(filePath)
	if zi != nil {
		r, err := zi.Open()
		if err != nil {
			log.Print(err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/png", r, nil)
	} else {
		NoRoute(c)
	}
	zf.Close()
}

func API_Difficulty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NoRoute(c)
		return
	}

	data := Dba.Select_Difficulty(id)
	if data.Id == nil {
		NoRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Session(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NoRoute(c)
		return
	}

	data := Dba.Select_Session(id)
	if data.Id == nil {
		NoRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Class(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NoRoute(c)
		return
	}

	data := Dba.Select_Class_Entries(id)
	if data.Id == nil {
		NoRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Time(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NoRoute(c)
		return
	}

	data := Dba.Select_Time_Weather(id)
	if data.Id == nil {
		NoRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func API_Recache_Content(c *gin.Context) {

	Parse_Content(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result":         "ok",
		"tracks_total":   len(Dba.Select_Cache_Tracks()),
		"cars_total":     len(Dba.Select_Cache_Cars()),
		"weathers_total": len(Dba.Select_Cache_Weathers()),
	})
}

func API_Validate_Installpath(c *gin.Context) {
	binary := "acServer"
	if runtime.GOOS == "windows" {
		binary = "acServer.exe"
	}
	path := filepath.Join(c.PostForm("path"), "server", binary)
	exists := true
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		exists = false
	}

	c.PureJSON(http.StatusOK, gin.H{
		"result": exists,
	})
}

func API_Server_Start(c *gin.Context) {
	if !Is_Running() {
		Status.Server_ApplyTrack()
		Start()
		// Hang the request until the UDP Server becomes online
		for start := time.Now(); time.Since(start) < time.Minute; {
			if Udp.online {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

func API_Server_Stop(c *gin.Context) {
	Stop()
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

func API_Server_Status(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

// Warning: non determistic as some of the items can be randomized/shuffled.
func API_Entry_List(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	Cr.Render_Ini(idInt)
	c.String(http.StatusOK, Cr.EntryList_Result)
}

func API_Server_Cfg(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	Cr.Render_Ini(idInt)
	c.String(http.StatusOK, Cr.ServerCfg_Result)
}

func API_Queue_MoveUp(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	Dba.Update_ServerEvent_MoveUp(idInt)

	c.String(http.StatusOK, "ok")
}

func API_Queue_MoveDown(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	Dba.Update_ServerEvent_MoveDown(idInt)

	c.String(http.StatusOK, "ok")
}


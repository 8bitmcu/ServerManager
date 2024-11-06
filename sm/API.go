package sm

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func API_Car_Image(c *gin.Context) {
	car := c.Param("car")
	skin := c.Param("skin")

	file := filepath.Join(Dba.Basepath(), "content", "cars", car, "skins", skin, "preview.jpg")
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		NoRoute(c)
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
		NoRoute(c)
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
		NoRoute(c)
		return
	} else {
		c.FileAttachment(file, fileName+".png")
		return
	}
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

func API_Recache_Cars(c *gin.Context) {
	Parse_Cars(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

func API_Recache_Tracks(c *gin.Context) {
	Parse_Tracks(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

func API_Recache_Weathers(c *gin.Context) {
	Parse_Weathers(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
func API_Recache_Content(c *gin.Context) {

	cars := 0
	tracks := 0
	weathers := 0
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		tracks = Parse_Tracks(Dba)
	}()
	go func() {
		defer wg.Done()
		cars = Parse_Cars(Dba)
	}()
	go func() {
		defer wg.Done()
		weathers = Parse_Weathers(Dba)
	}()
	wg.Wait()

	c.PureJSON(http.StatusOK, gin.H{
		"result":         "ok",
		"tracks_total":   tracks,
		"cars_total":     cars,
		"weathers_total": weathers,
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
		next_event := Dba.Select_Event_Next()
		cr := Cr.Render_Ini(*next_event.Id)
		cr.Write_Ini()

		tm := time.Now().Unix()
		next_event.ServerCfg = &cr.ServerCfg_Result
		next_event.EntryList = &cr.EntryList_Result


		next_event.Started_At = &tm

		Dba.Update_Event(next_event)

		Start()
	}
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text":       Get_Content(),
	})
}

func API_Server_Stop(c *gin.Context) {
	Stop()
	Dba.Update_Event_SetComplete()
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

	res := Cr.Render_Ini(idInt)
	c.String(http.StatusOK, res.EntryList_Result)
}

func API_Server_Cfg(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	res := Cr.Render_Ini(idInt)
	c.String(http.StatusOK, res.ServerCfg_Result)
}

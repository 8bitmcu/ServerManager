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

func noRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":    "PAGE_NOT_FOUND",
		"message": "Page not found",
	})
}

func apiCarImage(c *gin.Context) {
	car := c.Param("car")
	skin := c.Param("skin")

	var zf ZipFile
	zi := zf.FindZipFile("cars/" + car + "/skins/" + skin + "/preview.jpg")
	if zi != nil {
		r, err := zi.Open()
		if err != nil {
			log.Print("Cannot open car preview skin file in zipfile: ", err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/jpg", r, nil)
	} else {
		noRoute(c)
	}
	zf.Close()
}

func apiCar(c *gin.Context) {
	key := c.Param("key")

	cardata, err := Dba.selectCacheCar(key)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var power []any
	var torque []any
	var labels []any

	for _, l := range cardata.Power {
		labels = append(labels, l[0])
	}
	for _, p := range cardata.Power {
		power = append(power, p[1])
	}
	for _, t := range cardata.Torque {
		torque = append(torque, t[1])
	}

	c.PureJSON(http.StatusOK, gin.H{
		"key":    key,
		"desc":   cardata.Desc,
		"power":  power,
		"torque": torque,
		"labels": labels,
	})
}

func apiTrackPreviewImage(c *gin.Context) {
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
			log.Print("Could not open track preview from zipfile", err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/png", r, nil)
	} else {
		noRoute(c)
	}
	zf.Close()
}

func apiTrackOutlineImage(c *gin.Context) {
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
			log.Print("Could not open track outline from zipfile", err)
		}
		defer r.Close()
		c.DataFromReader(http.StatusOK, int64(zi.UncompressedSize64), "image/png", r, nil)
	} else {
		noRoute(c)
	}
	zf.Close()
}

func apiDifficulty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		noRoute(c)
		return
	}

	data, err := Dba.selectDifficulty(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if data.Id == nil {
		noRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func apiSession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		noRoute(c)
		return
	}

	data, err := Dba.selectSession(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if data.Id == nil {
		noRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func apiClass(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		noRoute(c)
		return
	}

	data, err := Dba.selectClassEntries(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if data.Id == nil {
		noRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func apiTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		noRoute(c)
		return
	}

	data, err := Dba.selectTimeWeather(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if data.Id == nil {
		noRoute(c)
		return
	}
	c.PureJSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func apiRecacheContent(c *gin.Context) {

	tracks, err := Dba.selectCacheTracks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	cars, err := Dba.selectCacheCars()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	weathers, err := Dba.selectCacheWeathers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	parseContent(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result":         "ok",
		"tracks_total":   len(tracks),
		"cars_total":     len(cars),
		"weathers_total": len(weathers),
	})
}

func apiValidateInstallpath(c *gin.Context) {
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

func apiServerStart(c *gin.Context) {
	if !isRunning() {
		Status.serverApplyTrack()
		start()
		// Hang the request until the UDP Server becomes online
		for start := time.Now(); time.Since(start) < time.Minute; {
			if Udp.online {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": isRunning(),
		"text":       getContent(),
	})
}

func apiServerStop(c *gin.Context) {
	stop()
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": isRunning(),
		"text":       getContent(),
	})
}

func apiServerStatus(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": isRunning(),
		"text":       getContent(),
	})
}

func apiServerLogfile(c *gin.Context) {
	logpath := filepath.Join(ConfigFolder, "logfile.log")
	c.FileAttachment(logpath, "logfile.log")
}

func apiServerSmdata(c *gin.Context) {
	smdata := filepath.Join(ConfigFolder, "smdata.db")
	c.FileAttachment(smdata, "smdata.db")
}

func apiServerSmcontent(c *gin.Context) {
	smcontent := filepath.Join(ConfigFolder, "smcontent.zip")
	c.FileAttachment(smcontent, "smcontent.zip")
}

// Warning: non determistic as some of the items can be randomized/shuffled.
func apiEntryList(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	Cr.renderIni(idInt)
	c.String(http.StatusOK, Cr.entryListResult)
}

func apiServerCfg(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)

	Cr.renderIni(idInt)
	c.String(http.StatusOK, Cr.serverCfgResult)
}

func apiQueueMoveUp(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	Dba.updateServerEventMoveUp(idInt)

	c.String(http.StatusOK, "ok")
}

func apiQueueMoveDown(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	Dba.updateServerEventMoveDown(idInt)

	c.String(http.StatusOK, "ok")
}

func apiQueueSkipEvent(c *gin.Context) {
	Status.serverChangeTrack()
	c.String(http.StatusOK, "ok")
}

func apiQueueClearCompleted(c *gin.Context) {
	Dba.deleteServerEventsCompleted()
	c.String(http.StatusOK, "ok")
}

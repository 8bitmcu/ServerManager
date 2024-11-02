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

	file := filepath.Join(basepath(), "content", "cars", car, "skins", skin, "preview.jpg")
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.FileAttachment(file, car+"_"+skin+".jpg")
		return
	}
}

func API_Track_Preview_Image(c *gin.Context) {
	track := c.Param("track")
	config := c.Param("config")

	file := ""
	fileName := ""
	if config != "" {
		file = filepath.Join(basepath(), "content", "tracks", track, "ui", config, "preview.png")
		fileName = track + "_" + config
	} else {
		file = filepath.Join(basepath(), "content", "tracks", track, "ui", "preview.png")
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
		file = filepath.Join(basepath(), "content", "tracks", track, "ui", config, "outline.png")
		fileName = track + "_" + config
	} else {
		file = filepath.Join(basepath(), "content", "tracks", track, "ui", "outline.png")
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
		"value": result,
	})
}

func API_Recache_Tracks(c *gin.Context) {
	result := Parse_Tracks(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
		"value": result,
	})
}

func API_Recache_Weathers(c *gin.Context) {
	result := Parse_Weathers(Dba)
	c.PureJSON(http.StatusOK, gin.H{
		"result": "ok",
		"value": result,
	})
}

func API_Validate_Installpath(c *gin.Context) {
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
		"text": Get_Content(),
	})
}

func API_Console_Stop(c *gin.Context) {
	Stop()
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text": Get_Content(),
	})
}

func API_Console_Status(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"is_running": Is_Running(),
		"text": Get_Content(),
	})
}

func API_Entry_List(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	c.String(http.StatusOK, Render_Entry_List(idInt))
}

func API_Server_Cfg(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	str, _ := Render_Ini(idInt)

	c.String(http.StatusOK, str)
}


/*
@app.route("/api/get_vehicle/<string:key>")
@require_config_set
def get_vehicle(key):
    car_data = dba.get_car(key)

    raw_power = json.loads(car_data['power'])

    torque = []
    for t in json.loads(car_data['torque']):
        torque.append(int(t[1]))

    power = []
    for p in raw_power:
        power.append(int(p[1]))

    labels = []
    for p in raw_power:
        labels.append(int(p[0]))

    json_data = {
        'key': key,
        'desc': car_data['desc'],
        'power': power,
        'torque': torque,
        'labels': labels
    }

    return jsonify(json_data)



*/

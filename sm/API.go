package sm

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// TODO: return 404 when id not found
func API_Difficulty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.PureJSON(200, gin.H{
			"data": Dba.Select_Difficulty(id),
		})
	}
}

func API_Session(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.PureJSON(200, gin.H{
			"data": Dba.Select_Session(id),
		})
	}
}

func API_Class(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.PureJSON(200, gin.H{
			"data": Dba.Select_Class_Entries(id),
		})
	}
}

func API_Time(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.htm", gin.H{})
		return
	} else {
		c.PureJSON(200, gin.H{
			"data": Dba.Select_Time_Weather(id),
		})
	}
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



@app.route("/api/validate_installpath", methods=('POST', ))
def validate_installpath():
    if fsa.validate_installpath(request.form['path']):
        return jsonify({"result": "ok"})
    else:
        return jsonify({"result": "no"})

@app.route("/api/recache_vehicles")
def recache_vehicles():
    result = fsa.parse_cars_folder(dba)
    return jsonify({'result': 'ok', 'value': result})

@app.route("/api/recache_tracks")
def recache_tracks():
    result = fsa.parse_tracks_folder(dba)
    return jsonify({'result': 'ok', 'value': result})

@app.route("/api/recache_weathers")
def recache_weathers():
    result = fsa.parse_weathers_folder(dba)
    return jsonify({'result': 'ok', 'value': result})


if __name__ == "__main__":
    app.run(debug=True)

*/

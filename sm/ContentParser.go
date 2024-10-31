package sm

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/kaptinlin/jsonrepair"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func basepath() string {
	return "/mnt/smb"
}

func Parse_Weathers(dba Dbaccess) int {
	r := regexp.MustCompile("^NAME")
	r2 := regexp.MustCompile("^NAME=(.*)$")
	r3 := regexp.MustCompile("; .*")

	weathers := make([]Cache_Weather, 0)

	weather_path := filepath.Join(basepath(), "content", "weather")
	entries, err := os.ReadDir(weather_path)
	if err != nil {
		log.Print(err)
	}

	for _, element := range entries {

		if !element.IsDir() {
			continue
		}

		ini_path := filepath.Join(weather_path, element.Name(), "weather.ini")

		if _, err := os.Stat(ini_path); errors.Is(err, os.ErrNotExist) {
			continue
		}

		file, err := os.Open(ini_path)
		if err != nil {
			log.Print(err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		name := ""
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Print(err)
				continue
			}
			t := scanner.Text()
			if r.MatchString(t) {
				matches := r2.FindStringSubmatch(t)
				name = r3.ReplaceAllString(matches[1], "")
			}
		}

		if name != "" {
			key := element.Name()
			weather := Cache_Weather{
				Key:  &key,
				Name: &name,
			}

			weathers = append(weathers, weather)
		}
	}

	dba.Update_Cache_Weathers(weathers)

	return len(weathers)
}

func Parse_Tracks(dba Dbaccess) int {
	parse_json := func(json_path string, key string, config string) Cache_Track {
		r := regexp.MustCompile("^[0-9]")
		jsonBytes, err := os.ReadFile(json_path)
		if err != nil {
			log.Print(err)
		}
		jsonStr := string(jsonBytes)

		data, err := jsonrepair.JSONRepair(jsonStr)
		if err != nil {
			log.Print(err)
		}

		var result map[string]string
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			//log.Print(err)
		}
		var tagsJson map[string][]string
		err = json.Unmarshal([]byte(data), &tagsJson)
		if err != nil {
			//log.Print(err)
		}

		tagsRes, err := json.Marshal(tagsJson["tags"])
		if err != nil {
			log.Print(err)
		}

		name := result["name"]
		desc := result["desc"]
		tags := string(tagsRes)
		country := result["country"]
		city := result["city"]
		strlen := r.ReplaceAllString(result["length"], "")
		length, err := strconv.Atoi(strlen)
		width, err := strconv.Atoi(result["width"])
		pitboxes, err := strconv.Atoi(result["pitboxes"])
		run := result["run"]

		track := Cache_Track{
			Key:      &key,
			Config:   &config,
			Name:     &name,
			Desc:     &desc,
			Tags:     &tags,
			Country:  &country,
			City:     &city,
			Length:   &length,
			Width:    &width,
			Pitboxes: &pitboxes,
			Run:      &run,
		}

		return track
	}

	tracks := make([]Cache_Track, 0)
	tracks_path := filepath.Join(basepath(), "content", "tracks")
	entries, err := os.ReadDir(tracks_path)
	if err != nil {
		log.Print(err)
	}

	for _, element := range entries {
		if !element.IsDir() {
			continue
		}

		// if skins is missing, assume it's a missing dlc and avoid listing/saving it
		skins := filepath.Join(tracks_path, element.Name(), "skins")
		if _, err := os.Stat(skins); errors.Is(err, os.ErrNotExist) {
			continue
		}

		json_path := filepath.Join(tracks_path, element.Name(), "ui", "ui_track.json")
		if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
			// track has many configs which we need to parse
			configs_folder := filepath.Join(tracks_path, element.Name(), "ui")
			configs, err := os.ReadDir(configs_folder)
			if err != nil {
				log.Print(err)
			}

			for _, config := range configs {
				if !config.IsDir() {
					continue
				}

				json_path = filepath.Join(tracks_path, element.Name(), "ui", config.Name(), "ui_track.json")

				if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
					// Does not exist, try dlc_ui_track.json
					json_path = filepath.Join(tracks_path, element.Name(), "ui", config.Name(), "dlc_ui_track.json")
					if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
						log.Print(err)
						continue
					}
				}

				track := parse_json(json_path, element.Name(), config.Name())
				tracks = append(tracks, track)
			}
		} else {
			// track has no config / only one selection
			track := parse_json(json_path, element.Name(), "")
			tracks = append(tracks, track)
		}
	}
	dba.Update_Cache_Tracks(tracks)

	return len(tracks)
}

func Parse_Cars(dba Dbaccess) int {
	cars := make([]Cache_Car, 0)
	cars_path := filepath.Join(basepath(), "content", "cars")
	entries, err := os.ReadDir(cars_path)
	if err != nil {
		log.Print(err)
	}

	for _, element := range entries {
		if !element.IsDir() {
			continue
		}

		// if data.acd is missing, assume it's a missing dlc and avoid listing/saving it
		data_acd := filepath.Join(cars_path, element.Name(), "data.acd")
		if _, err := os.Stat(data_acd); errors.Is(err, os.ErrNotExist) {
			continue
		}

		json_path := filepath.Join(cars_path, element.Name(), "ui", "ui_car.json")
		if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
			// file not exist, try dlc json
			json_path := filepath.Join(cars_path, element.Name(), "ui", "dlc_ui_car.json")
			if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
				log.Print(err)
			}
		}

		jsonBytes, err := os.ReadFile(json_path)
		if err != nil {
			log.Print(err)
		}
		jsonStr := string(jsonBytes)

		data, err := jsonrepair.JSONRepair(jsonStr)
		if err != nil {
			log.Print(err)
		}

		type jsonCar struct {
			Name  string `json:"name"`
			Desc  string `json:"desc"`
			Brand string `json:"brand"`
			Class string `json:"class"`

			Tags  []string `json:"tags"`
			Specs struct {
				Bhp          string `json:"bhp"`
				Torque       string `json:"torque"`
				Weight       string `json:"weight"`
				Topspeed     string `json:"topspeed"`
				Acceleration string `json:"acceleration"`
				Pwratio      string `json:"pwratio"`
				Range        int    `json:"range"`
			} `json:"specs"`
			Power  []any `json:"powerCurve"`
			Torque []any `json:"torqueCurve"`
		}
		var result jsonCar
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			log.Print(err)
		}

		specsRes, err := json.Marshal(result.Specs)
		if err != nil {
			log.Print(err)
		}
		specs := string(specsRes)

		tagsRes, err := json.Marshal(result.Tags)
		if err != nil {
			log.Print(err)
		}
		tags := string(tagsRes)

		pwrRes, err := json.Marshal(result.Power)
		if err != nil {
			log.Print(err)
		}
		power := string(pwrRes)

		trqRes, err := json.Marshal(result.Torque)
		if err != nil {
			log.Print(err)
		}
		torque := string(trqRes)

		skins_path := filepath.Join(cars_path, element.Name(), "skins")

		skins, err := os.ReadDir(skins_path)
		if err != nil {
			log.Print(err)
		}

		type jsonCarSkin struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		}

		carSkins := make([]jsonCarSkin, 0)
		for _, skin := range skins {

			skin_json := filepath.Join(skins_path, skin.Name(), "ui_skin.json")
			skin_name := skin.Name()
			if _, err := os.Stat(json_path); errors.Is(err, os.ErrNotExist) {
				//log.Print(err)
			} else {
				jsonBytes, err := os.ReadFile(skin_json)
				if err != nil {
					log.Print(err)
				}
				jsonStr := string(jsonBytes)

				data, err := jsonrepair.JSONRepair(jsonStr)
				if err != nil {
					log.Print(err)
				}

				var result map[string]string
				err = json.Unmarshal([]byte(data), &result)
				if err != nil {
					//log.Print(err)
				}

				skin_name = result["skinname"]
			}

			carSkin := jsonCarSkin{
				Key:  skin.Name(),
				Name: skin_name,
			}

			carSkins = append(carSkins, carSkin)
		}

		skinsRes, err := json.Marshal(carSkins)
		if err != nil {
			log.Print(err)
		}
		skinsStr := string(skinsRes)

		key := element.Name()
		name := result.Name
		desc := result.Desc
		brand := result.Brand
		class := result.Class

		car := Cache_Car{
			Key:    &key,
			Name:   &name,
			Brand:  &brand,
			Desc:   &desc,
			Tags:   &tags,
			Class:  &class,
			Specs:  &specs,
			Torque: &torque,
			Power:  &power,
			Skins:  &skinsStr,
		}
		cars = append(cars, car)
	}

	dba.Update_Cache_Cars(cars)
	return len(cars)
}

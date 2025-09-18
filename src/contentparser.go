package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/kaptinlin/jsonrepair"
)

func parseContent(dba Dbaccess) {
	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	basepath, err := dba.basepath()
	if err != nil {
		log.Print("Database Error: ", err)
		return
	}

	log.Print("Recreating smcontent.zip... please wait....")
	log.Print("Parsing ", basepath)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		tracks := parseTracks(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, tracks)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		cars := parseCars(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, cars)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		weathers := parseWeathers(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, weathers)
		mutex.Unlock()
	}()
	wg.Wait()

	zipfiles[filepath.Join(basepath, "server", "acServer")] = "acServer"
	zipfiles[filepath.Join(basepath, "server", "acServer.exe")] = "acServer.exe"
	zipfiles[filepath.Join(basepath, "system", "data", "surfaces.ini")] = "system/data/surfaces.ini"

	Zf.UpdateZipfile(zipfiles)
}

func parseWeathers(dba Dbaccess) map[string]string {
	zipfiles := map[string]string{}

	r := regexp.MustCompile("^NAME")
	r2 := regexp.MustCompile("^NAME=(.*)$")
	r3 := regexp.MustCompile("; .*")

	parseWeather := func(element fs.DirEntry, file io.Reader) CacheWeather {

		scanner := bufio.NewScanner(file)

		name := ""
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Print("Weather scanner failed to read: ", file, err)
				continue
			}
			t := scanner.Text()
			if r.MatchString(t) {
				matches := r2.FindStringSubmatch(t)
				name = r3.ReplaceAllString(matches[1], "")
				break
			}
		}

		key := element.Name()
		weather := CacheWeather{
			Key:  &key,
			Name: &name,
		}

		return weather
	}

	weathers := make([]CacheWeather, 0)

	basepath, err := dba.basepath()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil
	}
	weatherpath := filepath.Join(basepath, "content", "weather")
	entries, err := os.ReadDir(weatherpath)
	if err != nil {
		log.Print("Could not read directory: ", weatherpath, err)
	}

	var wg sync.WaitGroup
	for _, element := range entries {

		if !element.IsDir() {
			continue
		}

		inipath := filepath.Join(weatherpath, element.Name(), "weather.ini")
		previewpath := filepath.Join(weatherpath, element.Name(), "preview.jpg")

		if _, err := os.Stat(previewpath); errors.Is(err, os.ErrNotExist) {
			//log.Print("(warning) weather preview file missing: ", err)
		} else {
			zipfiles[previewpath] = "weather/" + element.Name() + "/preview.jpg"
		}

		if _, err := os.Stat(inipath); errors.Is(err, os.ErrNotExist) {
			log.Print("weather ini file missing: ", err)
			continue
		}
		file, err := os.Open(inipath)
		if err != nil {
			log.Print("Could not open weather ini file: ", inipath, err)
			continue
		}
		defer file.Close()

		wg.Add(1)
		go func() {
			defer wg.Done()
			weathers = append(weathers, parseWeather(element, file))
		}()
	}

	wg.Wait()

	dba.updateCacheWeathers(weathers)

	return zipfiles
}

// adds all .ini files to map[string]string recursively
func recurseAddIniZip(absPath string, relPath string) map[string]string {
	zipfiles := map[string]string{}
	files, err := os.ReadDir(absPath)
	if err != nil {
		log.Print("Could not read directory content at: ", absPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			zips := recurseAddIniZip(filepath.Join(absPath, file.Name()), relPath+"/"+file.Name())
			maps.Copy(zipfiles, zips)
		} else if strings.HasSuffix(file.Name(), ".ini") {
			if strings.HasPrefix(file.Name(), "models") || file.Name() == "drs_zones.ini" || file.Name() == "surfaces.ini" {
				zipfiles[filepath.Join(absPath, file.Name())] = relPath + "/" + file.Name()
			}
		}
	}
	return zipfiles
}

func parseTracks(dba Dbaccess) map[string]string {
	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	parsejson := func(jsonpath string, key string, config string) CacheTrack {
		r := regexp.MustCompile("[^0-9]")
		jsonBytes, err := os.ReadFile(jsonpath)
		if err != nil {
			log.Print("Could not read track json file: ", jsonpath, err)
		}
		jsonStr := string(jsonBytes)

		data, err := jsonrepair.JSONRepair(jsonStr)
		if err != nil {
			log.Print("Could not repair track json file:", jsonpath, err)
		}

		var result map[string]string
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			//log.Print("Could not unmarshal track json file to map[string]string: ", jsonpath, err)
		}

		var track CacheTrack
		err = json.Unmarshal([]byte(data), &track)
		if err != nil {
			//log.Print("Could not unmarshal track json file to CacheTrack: ", jsonpath, err)
		}
		tracklen := r.ReplaceAllString(result["length"], "")
		parsedLen, err := strconv.Atoi(tracklen)
		if err != nil {
			log.Print("Could not convert track length to integer: ", tracklen, err)
		}
		track.Key = &key
		track.Config = &config
		track.Length = &parsedLen

		return track
	}

	basepath, err := dba.basepath()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil
	}
	tracks := make([]CacheTrack, 0)
	trackspath := filepath.Join(basepath, "content", "tracks")

	parseTrack := func(element fs.DirEntry) {
		if !element.IsDir() {
			return
		}
		zips := recurseAddIniZip(filepath.Join(trackspath, element.Name()), "tracks/"+element.Name())
		mutex.Lock()
		maps.Copy(zipfiles, zips)
		mutex.Unlock()

		jsonpath := filepath.Join(trackspath, element.Name(), "ui", "ui_track.json")
		if _, err := os.Stat(jsonpath); errors.Is(err, os.ErrNotExist) {
			// track has many configs which we need to parse
			configsfolder := filepath.Join(trackspath, element.Name(), "ui")
			configs, err := os.ReadDir(configsfolder)
			if err != nil {
				log.Print("Could not read track config folder: ", configsfolder, err)
			}

			for _, config := range configs {
				if !config.IsDir() {
					continue
				}

				jsonpath = filepath.Join(trackspath, element.Name(), "ui", config.Name(), "ui_track.json")

				if _, err := os.Stat(jsonpath); errors.Is(err, os.ErrNotExist) {
					// Does not exist, try dlc_ui_track.json
					jsonpath = filepath.Join(trackspath, element.Name(), "ui", config.Name(), "dlc_ui_track.json")
					if _, err := os.Stat(jsonpath); errors.Is(err, os.ErrNotExist) {
						log.Print("No valid track config json file found for: ", err)
						continue
					}
				}

				mutex.Lock()
				outline := filepath.Join(trackspath, element.Name(), "ui", config.Name(), "outline.png")
				if _, err := os.Stat(outline); errors.Is(err, os.ErrNotExist) {
					//log.Print("(warning) Track outline file does not exist: ", err)
				} else {
					zipfiles[outline] = "tracks/" + element.Name() + "/" + config.Name() + "/outline.png"
				}
				preview := filepath.Join(trackspath, element.Name(), "ui", config.Name(), "preview.png")
				if _, err := os.Stat(preview); errors.Is(err, os.ErrNotExist) {
					//log.Print("(warning) Track preview file does not exist: ", err)
				} else {
					zipfiles[preview] = "tracks/" + element.Name() + "/" + config.Name() + "/preview.png"
				}
				mutex.Unlock()

				track := parsejson(jsonpath, element.Name(), config.Name())
				tracks = append(tracks, track)
			}
		} else {
			// track has no config / only one selection
			mutex.Lock()
			zipfiles[filepath.Join(trackspath, element.Name(), "ui", "outline.png")] = "tracks/" + element.Name() + "/outline.png"
			zipfiles[filepath.Join(trackspath, element.Name(), "ui", "preview.png")] = "tracks/" + element.Name() + "/preview.png"
			mutex.Unlock()

			track := parsejson(jsonpath, element.Name(), "")
			tracks = append(tracks, track)
		}
	}

	entries, err := os.ReadDir(trackspath)
	if err != nil {
		log.Print("Can not read tracks directory: ", trackspath, err)
	}

	var wg sync.WaitGroup
	for _, element := range entries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parseTrack(element)
		}()
	}
	wg.Wait()

	dba.updateCacheTracks(tracks)

	return zipfiles
}

func parseCars(dba Dbaccess) map[string]string {

	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	type JsonCarSkin struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}

	parseSkin := func(skin fs.DirEntry, skinspath string) JsonCarSkin {
		skinjson := filepath.Join(skinspath, skin.Name(), "ui_skin.json")
		skinname := skin.Name()
		if _, err := os.Stat(skinjson); errors.Is(err, os.ErrNotExist) {
			log.Print("Car ui_skin.json file missing: ", skinjson, skinname, err)
		} else {
			jsonBytes, err := os.ReadFile(skinjson)
			if err != nil {
				log.Print("Cannot read ui_skin.json file: ", skinjson, skinname, err)
			}
			jsonStr := string(jsonBytes)

			data, err := jsonrepair.JSONRepair(jsonStr)
			if err != nil {
				log.Print("Cannot repair ui_skin.json file: ", skinjson, skinname, err)
			}

			var result map[string]string
			err = json.Unmarshal([]byte(data), &result)
			if err != nil {
				//log.Print("Cannot unmarshal ui_skin.json file: ", skinjson, skinname, err)
			}

			if result["skinname"] != "" {
				skinname = result["skinname"]
			}
		}

		carSkin := JsonCarSkin{
			Key:  skin.Name(),
			Name: skinname,
		}

		return carSkin
	}

	parseCar := func(element fs.DirEntry, carspath string) CacheCar {

		jsonpath := filepath.Join(carspath, element.Name(), "ui", "ui_car.json")
		if _, err := os.Stat(jsonpath); errors.Is(err, os.ErrNotExist) {
			// file not exist, try dlc json
			jsonpath := filepath.Join(carspath, element.Name(), "ui", "dlc_ui_car.json")
			if _, err := os.Stat(jsonpath); errors.Is(err, os.ErrNotExist) {
				log.Print("No json file found for car: ", element.Name(), err)
			}
		}

		jsonBytes, err := os.ReadFile(jsonpath)
		if err != nil {
			log.Print("Could not read car json file: ", jsonpath, err)
		}
		jsonStr := string(jsonBytes)

		data, err := jsonrepair.JSONRepair(jsonStr)
		if err != nil {
			log.Print("Could not repair car json file: ", jsonpath, err)
		}

		var result CacheCar
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			//log.Print("Could not unmarshal car json file: ", jsonpath, err)
		}

		skinspath := filepath.Join(carspath, element.Name(), "skins")

		skins, err := os.ReadDir(skinspath)
		if err != nil {
			log.Print("Could not read car skin directory: ", skinspath, err)
		}

		var wg sync.WaitGroup
		for _, skin := range skins {
			if !skin.IsDir() {
				continue
			}

			mutex.Lock()
			zipfiles[filepath.Join(skinspath, skin.Name(), "preview.jpg")] = "cars/" + element.Name() + "/skins/" + skin.Name() + "/preview.jpg"
			mutex.Unlock()

			wg.Add(1)
			go func() {
				defer wg.Done()
				result.Skins = append(result.Skins, parseSkin(skin, skinspath))
			}()
		}

		wg.Wait()

		key := element.Name()
		result.Key = &key

		return result
	}

	basepath, err := dba.basepath()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil
	}
	cars := make([]CacheCar, 0)
	carspath := filepath.Join(basepath, "content", "cars")
	entries, err := os.ReadDir(carspath)
	if err != nil {
		log.Print("Could not read cars directory: ", carspath, err)
	}

	var wg sync.WaitGroup
	for _, element := range entries {
		if !element.IsDir() {
			continue
		}

		// if data.acd or data dir is missing, assume it's a missing dlc and avoid listing/saving it
		dataacd := filepath.Join(carspath, element.Name(), "data.acd")
		if _, err := os.Stat(dataacd); errors.Is(err, os.ErrNotExist) {
			data := filepath.Join(carspath, element.Name(), "data")
			if _, err := os.Stat(data); errors.Is(err, os.ErrNotExist) {
				continue
			}
		} else {
			// append data.acd to our zipfile
			mutex.Lock()
			zipfiles[dataacd] = "cars/" + element.Name() + "/data.acd"
			mutex.Unlock()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			cars = append(cars, parseCar(element, carspath))
		}()
	}

	wg.Wait()

	dba.updateCacheCars(cars)

	return zipfiles
}

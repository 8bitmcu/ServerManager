package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/akutz/sortfold"
	"github.com/kaptinlin/jsonrepair"
	"golang.org/x/image/draw"
)

func updateZipfile(filesToZip map[string]string) {
	zipfilename := "smcontent.zip"
	if _, err := os.Stat(zipfilename); errors.Is(err, os.ErrNotExist) {
		newfile, err := os.Create(zipfilename)
		if err != nil {
			log.Print(err)
		}
		defer newfile.Close()

		w := zip.NewWriter(newfile)

		f, err := w.Create("readme.txt")
		if err != nil {
			log.Print(err)
		}
		_, err = f.Write([]byte("This archive is maintained by servermanager."))
		if err != nil {
			log.Print(err)
		}

		err = w.Close()
		if err != nil {
			log.Print(err)
		}
	}

	zr, err := zip.OpenReader(zipfilename)
	if err != nil {
		log.Print(err)
	}
	defer zr.Close()
	zwf, err := os.Create(zipfilename + "_")
	defer zwf.Close()
	zw := zip.NewWriter(zwf)
	defer zwf.Close()

	log.Print("Compressing files...")

	defer zw.Close()
	var wg sync.WaitGroup

	keys := make([]string, 0, len(filesToZip))
	for k := range filesToZip {
		keys = append(keys, k)
	}

	// sort keys by name, insensitively
	sort.Slice(keys, func(i, j int) bool {
		return sortfold.CompareFold(keys[i], keys[j]) < 0
	})

	newFiles := make([]string, 0)
	for _, filepath := range keys {
		destination := filesToZip[filepath]
		wg.Add(1)
		func() {
			defer wg.Done()

			src, err := os.Open(filepath)
			if err != nil {
				log.Print(err)
			}

			fi, err := src.Stat()
			if err != nil {
				log.Print(err)
			}

			// Skip file if it already exists and has the same timestamp
			for _, zipItem := range zr.File {
				if zipItem.Name == destination && zipItem.Modified.Unix() == fi.ModTime().Unix() {
					return
				}
			}

			fih := &zip.FileHeader{
				Name:     destination,
				Method:   zip.Deflate,
				Modified: fi.ModTime(),
			}
			if err != nil {
				log.Print(err)
			}

			dest, err := zw.CreateHeader(fih)
			if err != nil {
				log.Print(err)
			}

			defer src.Close()
			if strings.HasSuffix(filepath, ".jpg") || strings.HasSuffix(filepath, ".jpeg") {
				img, err := jpeg.Decode(src)

				if err != nil {
					if _, err := io.Copy(dest, src); err != nil {
						log.Print(err)
					}
					return
				}

				width := 640
				height := int(float32(img.Bounds().Max.Y) / float32(img.Bounds().Max.X) * float32(width))
				if img.Bounds().Max.X < width {
					width = img.Bounds().Max.X
					height = img.Bounds().Max.Y
				}

				dst := image.NewRGBA(image.Rect(0, 0, width, height))
				draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

				jpeg.Encode(dest, dst, &jpeg.Options{
					Quality: jpeg.DefaultQuality,
				})
			} else if strings.HasSuffix(filepath, ".png") {
				img, err := png.Decode(src)

				if err != nil {
					if _, err := io.Copy(dest, src); err != nil {
						log.Print(err)
					}
					return
				}

				width := 640
				height := int(float32(img.Bounds().Max.Y) / float32(img.Bounds().Max.X) * float32(width))
				if img.Bounds().Max.X < width {
					width = img.Bounds().Max.X
					height = img.Bounds().Max.Y
				}

				dst := image.NewRGBA(image.Rect(0, 0, width, height))
				draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

				png.Encode(dest, dst)
			} else {
				if _, err := io.Copy(dest, src); err != nil {
					log.Print(err)
				}
			}

			newFiles = append(newFiles, destination)
		}()
	}

	wg.Wait()

	log.Print("Copying zipfile content...")
	inNewFiles := func(value string) bool {
		for _, item := range newFiles {
			if strings.ToLower(item) == strings.ToLower(value) {
				return true
			}
		}
		return false
	}

	for _, zipItem := range zr.File {
		if inNewFiles(zipItem.Name) {
			continue
		}

		zipItemReader, err := zipItem.OpenRaw()
		if err != nil {
			log.Print(err)
		}

		header := zipItem.FileHeader
		targetItem, err := zw.CreateRaw(&header)
		_, err = io.Copy(targetItem, zipItemReader)
		if err != nil {
			log.Print(err)
		}
	}

	os.Remove(zipfilename)
	os.Rename(zipfilename+"_", zipfilename)
}

func Parse_Content(dba Dbaccess) {
	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	log.Print("Parsing ", dba.Basepath())
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		tracks := parse_Tracks(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, tracks)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		cars := parse_Cars(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, cars)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		weathers := parse_Weathers(Dba)
		mutex.Lock()
		maps.Copy(zipfiles, weathers)
		mutex.Unlock()
	}()
	wg.Wait()

	zipfiles[filepath.Join(dba.Basepath(), "server", "acServer")] = "acServer"
	zipfiles[filepath.Join(dba.Basepath(), "server", "acServer.exe")] = "acServer.exe"

	updateZipfile(zipfiles)
	log.Print("Content Updated")
}

func parse_Weathers(dba Dbaccess) map[string]string {
	zipfiles := map[string]string{}

	r := regexp.MustCompile("^NAME")
	r2 := regexp.MustCompile("^NAME=(.*)$")
	r3 := regexp.MustCompile("; .*")

	parseWeather := func(element fs.DirEntry, file io.Reader) Cache_Weather {

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

		key := element.Name()
		weather := Cache_Weather{
			Key:  &key,
			Name: &name,
		}

		return weather
	}

	weathers := make([]Cache_Weather, 0)

	weather_path := filepath.Join(Dba.Basepath(), "content", "weather")
	entries, err := os.ReadDir(weather_path)
	if err != nil {
		log.Print(err)
	}

	var wg sync.WaitGroup
	for _, element := range entries {

		if !element.IsDir() {
			continue
		}

		ini_path := filepath.Join(weather_path, element.Name(), "weather.ini")
		preview_path := filepath.Join(weather_path, element.Name(), "preview.jpg")

		if _, err := os.Stat(preview_path); errors.Is(err, os.ErrNotExist) {
		} else {
			zipfiles[preview_path] = "weather/" + element.Name() + "/preview.jpg"
		}

		if _, err := os.Stat(ini_path); errors.Is(err, os.ErrNotExist) {
			continue
		}
		file, err := os.Open(ini_path)
		if err != nil {
			log.Print(err)
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

	dba.Update_Cache_Weathers(weathers)

	return zipfiles
}

// adds all .ini files to map[string]string recursively
func recurseAddIniZip(absPath string, relPath string) map[string]string {
	zipfiles := map[string]string{}
	files, err := os.ReadDir(absPath)
	if err != nil {
		log.Print(err)
	}

	for _, file := range files {
		if file.IsDir() {
			zips := recurseAddIniZip(filepath.Join(absPath, file.Name()), relPath+"/"+file.Name())
			maps.Copy(zipfiles, zips)
		} else if strings.HasSuffix(file.Name(), ".ini") {
			zipfiles[filepath.Join(absPath, file.Name())] = relPath + "/" + file.Name()
		}
	}
	return zipfiles
}

func parse_Tracks(dba Dbaccess) map[string]string {
	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	parse_json := func(json_path string, key string, config string) Cache_Track {
		r := regexp.MustCompile("[^0-9]")
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

		var track Cache_Track
		err = json.Unmarshal([]byte(data), &track)
		if err != nil {
			//log.Print(err)
		}

		parsedLen, err := strconv.Atoi(r.ReplaceAllString(result["length"], ""))
		track.Key = &key
		track.Config = &config
		track.Length = &parsedLen

		return track
	}

	tracks := make([]Cache_Track, 0)
	tracks_path := filepath.Join(Dba.Basepath(), "content", "tracks")

	parseTrack := func(element fs.DirEntry) {
		if !element.IsDir() {
			return
		}
		zips := recurseAddIniZip(filepath.Join(tracks_path, element.Name()), "tracks/"+element.Name())
		mutex.Lock()
		maps.Copy(zipfiles, zips)
		mutex.Unlock()

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

				mutex.Lock()
				outline := filepath.Join(tracks_path, element.Name(), "ui", config.Name(), "outline.png")
				if _, err := os.Stat(outline); errors.Is(err, os.ErrNotExist) {
				} else {
					zipfiles[outline] = "tracks/" + element.Name() + "/" + config.Name() + "/outline.png"
				}
				preview := filepath.Join(tracks_path, element.Name(), "ui", config.Name(), "preview.png")
				if _, err := os.Stat(preview); errors.Is(err, os.ErrNotExist) {
				} else {
					zipfiles[preview] = "tracks/" + element.Name() + "/" + config.Name() + "/preview.png"
				}
				mutex.Unlock()

				track := parse_json(json_path, element.Name(), config.Name())
				tracks = append(tracks, track)
			}
		} else {
			// track has no config / only one selection
			mutex.Lock()
			zipfiles[filepath.Join(tracks_path, element.Name(), "ui", "outline.png")] = "tracks/" + element.Name() + "/outline.png"
			zipfiles[filepath.Join(tracks_path, element.Name(), "ui", "preview.png")] = "tracks/" + element.Name() + "/preview.png"
			mutex.Unlock()

			track := parse_json(json_path, element.Name(), "")
			tracks = append(tracks, track)
		}
	}

	entries, err := os.ReadDir(tracks_path)
	if err != nil {
		log.Print(err)
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

	dba.Update_Cache_Tracks(tracks)

	return zipfiles
}

func parse_Cars(dba Dbaccess) map[string]string {

	zipfiles := map[string]string{}
	var mutex = &sync.RWMutex{}

	type jsonCarSkin struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}

	parseSkin := func(skin fs.DirEntry, skins_path string) jsonCarSkin {
		skin_json := filepath.Join(skins_path, skin.Name(), "ui_skin.json")
		skin_name := skin.Name()
		if _, err := os.Stat(skin_json); errors.Is(err, os.ErrNotExist) {
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

			if result["skinname"] != "" {
				skin_name = result["skinname"]
			}
		}

		carSkin := jsonCarSkin{
			Key:  skin.Name(),
			Name: skin_name,
		}

		return carSkin
	}

	parseCar := func(element fs.DirEntry, cars_path string) Cache_Car {

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

		var result Cache_Car
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			//log.Print(err)
		}

		skins_path := filepath.Join(cars_path, element.Name(), "skins")

		skins, err := os.ReadDir(skins_path)
		if err != nil {
			log.Print(err)
		}

		var wg sync.WaitGroup
		for _, skin := range skins {
			if !skin.IsDir() {
				continue
			}

			mutex.Lock()
			zipfiles[filepath.Join(skins_path, skin.Name(), "preview.jpg")] = "cars/" + element.Name() + "/skins/" + skin.Name() + "/preview.jpg"
			mutex.Unlock()

			wg.Add(1)
			go func() {
				defer wg.Done()
				result.Skins = append(result.Skins, parseSkin(skin, skins_path))
			}()
		}

		wg.Wait()

		key := element.Name()
		result.Key = &key

		return result
	}

	cars := make([]Cache_Car, 0)
	cars_path := filepath.Join(Dba.Basepath(), "content", "cars")
	entries, err := os.ReadDir(cars_path)
	if err != nil {
		log.Print(err)
	}

	var wg sync.WaitGroup
	for _, element := range entries {
		if !element.IsDir() {
			continue
		}

		// if data.acd or data dir is missing, assume it's a missing dlc and avoid listing/saving it
		data_acd := filepath.Join(cars_path, element.Name(), "data.acd")
		if _, err := os.Stat(data_acd); errors.Is(err, os.ErrNotExist) {
			data := filepath.Join(cars_path, element.Name(), "data")
			if _, err := os.Stat(data); errors.Is(err, os.ErrNotExist) {
				//log.Print("  skipping folder '", element.Name(), "' (most likely a missing DLC)")
				continue
			}

			// append all *.ini files in data/ to our zipfile
			dfiles, err := os.ReadDir(data)
			if err != nil {
				log.Print(err)
			}
			for _, fn := range dfiles {
				if fn.IsDir() {
					continue
				}
				if strings.HasSuffix(fn.Name(), ".ini") {

					mutex.Lock()
					zipfiles[data+"/"+fn.Name()] = "cars/" + element.Name() + "/data/" + fn.Name()
					mutex.Unlock()
				}
			}
		} else {
			// append data.acd to our zipfile
			mutex.Lock()
			zipfiles[data_acd] = "cars/" + element.Name() + "/data.acd"
			mutex.Unlock()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			cars = append(cars, parseCar(element, cars_path))
		}()
	}

	wg.Wait()

	dba.Update_Cache_Cars(cars)

	return zipfiles
}

package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type ConfigRenderer struct {
	serverCfgIni    *template.Template
	entryListIni    *template.Template
	serverCfgResult string
	entryListResult string
	class           UserClass
	track           CacheTrack
	cspRequired     bool
	maxClients      int
	serverEvent     ServerEvent
}

// 8:00 AM = -80
// 18:00 PM = 80
// increment of 8 every 30 minutes
func (cr *ConfigRenderer) timeToSunAngle(timeStr *string) int {
	time, err := time.Parse("15:04", *timeStr)
	if err != nil {
		log.Print("Could not parse time: ", *timeStr, err)
	}

	angle := -80 + (16 * (time.Hour() - 8))
	angle = angle + int(math.Round(float64(time.Minute())/15))*4

	return angle
}

func (cr *ConfigRenderer) writeIni() {
	cfgfolder := filepath.Join(TempFolder, "cfg")
	err := os.MkdirAll(cfgfolder, os.ModePerm)
	if err != nil {
		log.Print("Could not create temp folder: ", cfgfolder, err)
	}

	err = os.WriteFile(filepath.Join(cfgfolder, "server_cfg.ini"), []byte(cr.serverCfgResult), 0644)
	if err != nil {
		log.Print("Could not write server_cfg.ini: ", err)
	}

	err = os.WriteFile(filepath.Join(cfgfolder, "entry_list.ini"), []byte(cr.entryListResult), 0644)
	if err != nil {
		log.Print("Could not write entry_list.ini: ", err)
	}
}

func (cr *ConfigRenderer) renderIni(eventId int) {
	r := regexp.MustCompile(`\d{1,3}`)

	event, err := Dba.selectEvent(eventId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	eventcat, err := Dba.selectEventCategory(*event.EventCategoryId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	tm, err := Dba.selectTimeWeather(*event.TimeId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	diff, err := Dba.selectDifficulty(*event.DifficultyId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	cfg, err := Dba.selectConfig()
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	session, err := Dba.selectSession(*event.SessionId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	class, err := Dba.selectClassEntries(*event.ClassId)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	track, err := Dba.selectCacheTrack(*event.CacheTrackKey, *event.CacheTrackConfig)
	if err != nil {
		log.Print("Database error: ", err)
		return
	}

	if session.QualifyMaxWaitPerc == nil {
		var val = 120
		session.QualifyMaxWaitPerc = &val
	}

	// csp required? build a cspstr to be concat with the track name
	cspstr := ""
	if cfg.CspRequired != nil && *cfg.CspRequired > 0 {
		Cr.cspRequired = true
		cspLetter := ""

		if *cfg.CspPhycars > 0 && *cfg.CspPhytracks > 0 && *cfg.CspHidepit > 0 {
			cspLetter = "/../H"
		} else if *cfg.CspPhycars > 0 && *cfg.CspPhytracks > 0 {
			cspLetter = "/../D"
		} else if *cfg.CspPhycars > 0 && *cfg.CspHidepit > 0 {
			cspLetter = "/../F"
		} else if *cfg.CspPhytracks > 0 && *cfg.CspHidepit > 0 {
			cspLetter = "/../G"
		} else if *cfg.CspPhycars > 0 {
			cspLetter = "/../B"
		} else if *cfg.CspPhytracks > 0 {
			cspLetter = "/../C"
		} else if *cfg.CspHidepit > 0 {
			cspLetter = "/../E"
		}

		cspstr = "csp/" + strconv.Itoa(*cfg.CspVersion) + cspLetter + "/../"
	}

	// Weather CSP? build new graphics string
	if *tm.CspEnabled == 1 {
		t := "13:00" // sets the sun angle to zero; a "nice" default/backup value
		tm.Time = &t
		todm := 1
		tm.TimeOfDayMulti = &todm
		for _, wt := range tm.Weathers {
			cspTime, err := time.Parse("15:04", *wt.CspTime)
			if err != nil {
				log.Print("Could not parse csp time:", *wt.CspTime, err)
			}
			cspTimeInt := (cspTime.Hour() * 3600) + (cspTime.Minute() * 60) + cspTime.Second()
			seconds := strconv.Itoa(cspTimeInt)
			mult := strconv.Itoa(*wt.CspTimeOfDayMulti)

			dateStr := ""
			if wt.CspDate != nil && *wt.CspDate != "" {
				cspDate, err := time.Parse("2006-01-02", *wt.CspDate)
				if err != nil {
					log.Print("Could not parse csp date:", *wt.CspDate, err)
				}

				dateStr = "_start=" + strconv.FormatInt(cspDate.Unix(), 10)
			}

			matches := r.FindStringSubmatch(*wt.Graphics)
			*wt.Graphics = matches[0] + "_time=" + seconds + "_mult=" + mult + dateStr
		}
	}

	// Maximum clients defined as the minimum between maxclients, pitboxes and vehicles in class
	stratneeded := false
	maxclients := *cfg.MaxClients

	if *track.Pitboxes < *cfg.MaxClients {
		stratneeded = true
		maxclients = *track.Pitboxes
	}
	if len(class.Entries) < maxclients {
		stratneeded = true
		maxclients = len(class.Entries)
	}
	if len(class.Entries) > maxclients {
		stratneeded = true
	}

	// Strategy needed? re-order cars in the entry list as per selected strategy
	if stratneeded {
		// Random
		if *event.Strategy == 2 {
			rand.Shuffle(len(class.Entries), func(i, j int) { class.Entries[i], class.Entries[j] = class.Entries[j], class.Entries[i] })
		}
		// Cut the list by the max number of clients
		class.Entries = class.Entries[:maxclients]
	}

	funcMap := template.FuncMap{
		"derefInt": func(i *int) int {
			return *i
		},
	}

	if Cr.serverCfgIni == nil {
		file := FindFile("/ini/server_cfg.ini")
		tmplStr, err := io.ReadAll(file)
		if err != nil {
			log.Print("Could not read template file server_cfg.ini: ", err)
		}
		Cr.serverCfgIni, err = template.New("server_cfg.ini").Funcs(funcMap).Parse(string(tmplStr))
		if err != nil {
			log.Print("Error parsing server_cfg.ini template: ", err)
		}
	}

	data := map[string]any{
		"event":       event,
		"config":      cfg,
		"diff":        diff,
		"session":     session,
		"time":        tm,
		"class":       class,
		"track":       track,
		"max_clients": maxclients,
		"sunangle":    cr.timeToSunAngle(tm.Time),
		"cspstr":      cspstr,
		"name":        eventcat.Name,
	}

	var b bytes.Buffer
	err = Cr.serverCfgIni.Execute(&b, data)
	if err != nil {
		log.Print("Error executing server_cfg.ini template: ", err)
	}

	if Cr.entryListIni == nil {
		file := FindFile("/ini/entry_list.ini")
		tmplStr, err := io.ReadAll(file)
		if err != nil {
			log.Print("Could not read template file entry_list.ini: ", err)
		}
		Cr.entryListIni, err = template.New("entry_list.ini").Funcs(funcMap).Parse(string(tmplStr))
		if err != nil {
			log.Print("Error parsing entry_list.ini template: ", err)
		}
	}

	var b2 bytes.Buffer
	err = Cr.entryListIni.Execute(&b2, class)
	if err != nil {
		log.Print("Error executing entry_list.ini template: ", err)
	}

	cr.serverCfgResult = b.String()
	cr.entryListResult = b2.String()
	cr.class = class
	cr.track = track
	cr.maxClients = maxclients
}

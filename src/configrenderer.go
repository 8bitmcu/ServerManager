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
	server_cfg_ini   *template.Template
	entry_list_ini   *template.Template
	ServerCfg_Result string
	EntryList_Result string
	Class            User_Class
	Track            Cache_Track
	Csp_Required     bool
}

// 8:00 AM = -80
// 18:00 PM = 80
// increment of 8 every 30 minutes
func (cr *ConfigRenderer) Time_To_SunAngle(time_str *string) int {
	time, err := time.Parse("15:04", *time_str)
	if err != nil {
		log.Print(err)
	}

	angle := -80 + (16 * (time.Hour() - 8))
	angle = angle + int(math.Round(float64(time.Minute())/15))*4

	return angle
}

func (cr *ConfigRenderer) Write_Ini() {
	cfgfolder := filepath.Join(TempFolder, "cfg")
	err := os.MkdirAll(cfgfolder, os.ModePerm)
	if err != nil {
		log.Print(err)
	}

	err = os.WriteFile(filepath.Join(cfgfolder, "server_cfg.ini"), []byte(cr.ServerCfg_Result), 0644)
	if err != nil {
		log.Print(err)
	}

	err = os.WriteFile(filepath.Join(cfgfolder, "entry_list.ini"), []byte(cr.EntryList_Result), 0644)
	if err != nil {
		log.Print(err)
	}
}

func (cr *ConfigRenderer) Render_Ini(event_id int) {
	r := regexp.MustCompile("\\d{1,3}")

	event := Dba.Select_Event(event_id)
	event_cat := Dba.Select_Event_Category(*event.Event_Category_Id)
	tm := Dba.Select_Time_Weather(*event.Time_Id)
	diff := Dba.Select_Difficulty(*event.Difficulty_Id)
	cfg := Dba.Select_Config()
	session := Dba.Select_Session(*event.Session_Id)
	class := Dba.Select_Class_Entries(*event.Class_Id)
	track := Dba.Select_Cache_Track(*event.Cache_Track_Key, *event.Cache_Track_Config)

	if session.Qualify_Max_Wait_Perc == nil {
		var val = 120
		session.Qualify_Max_Wait_Perc = &val
	}

	// csp required? build a cspstr to be concat with the track name
	cspstr := ""
	if cfg.Csp_Required != nil && *cfg.Csp_Required > 0 {
		Cr.Csp_Required = true
		csp_letter := ""

		if *cfg.Csp_Phycars > 0 && *cfg.Csp_Phytracks > 0 && *cfg.Csp_Hidepit > 0 {
			csp_letter = "/../H"
		} else if *cfg.Csp_Phycars > 0 && *cfg.Csp_Phytracks > 0 {
			csp_letter = "/../D"
		} else if *cfg.Csp_Phycars > 0 && *cfg.Csp_Hidepit > 0 {
			csp_letter = "/../F"
		} else if *cfg.Csp_Phytracks > 0 && *cfg.Csp_Hidepit > 0 {
			csp_letter = "/../G"
		} else if *cfg.Csp_Phycars > 0 {
			csp_letter = "/../B"
		} else if *cfg.Csp_Phytracks > 0 {
			csp_letter = "/../C"
		} else if *cfg.Csp_Hidepit > 0 {
			csp_letter = "/../E"
		}

		cspstr = "csp/" + strconv.Itoa(*cfg.Csp_Version) + csp_letter + "/../"
	}

	// Weather CSP? build new graphics string
	if *tm.Csp_Enabled == 1 {
		t := "13:00" // sets the sun angle to zero; a "nice" default/backup value
		tm.Time = &t
		todm := 1
		tm.Time_Of_Day_Multi = &todm
		for _, wt := range tm.Weathers {
			csp_time, err := time.Parse("15:04", *wt.Csp_Time)
			if err != nil {
				log.Print(err)
			}
			csp_timeInt := (csp_time.Hour() * 3600) + (csp_time.Minute() * 60) + csp_time.Second()
			seconds := strconv.Itoa(csp_timeInt)
			mult := strconv.Itoa(*wt.Csp_Time_Of_Day_Multi)

			dateStr := ""
			if wt.Csp_Date != nil && *wt.Csp_Date != "" {
				csp_date, err := time.Parse("2006-01-02", *wt.Csp_Date)
				if err != nil {
					log.Print(err)
				}

				dateStr = "_start=" + strconv.FormatInt(csp_date.Unix(), 10)
			}

			matches := r.FindStringSubmatch(*wt.Graphics)
			*wt.Graphics = matches[0] + "_time=" + seconds + "_mult=" + mult + dateStr
		}
	}

	// Maximum clients defined as the minimum between max_clients, pitboxes and vehicles in class
	strat_needed := false
	max_clients := *cfg.Max_Clients

	if *track.Pitboxes < *cfg.Max_Clients {
		strat_needed = true
		max_clients = *track.Pitboxes
	}
	if len(class.Entries) < max_clients {
		strat_needed = true
		max_clients = len(class.Entries)
	}
	if len(class.Entries) > max_clients {
		strat_needed = true
	}

	// Strategy needed? re-order cars in the entry list as per selected strategy
	if strat_needed {
		// Random
		if *event.Strategy == 2 {
			rand.Shuffle(len(class.Entries), func(i, j int) { class.Entries[i], class.Entries[j] = class.Entries[j], class.Entries[i] })
		}
		// Cut the list by the max number of clients
		class.Entries = class.Entries[:max_clients]
	}

	funcMap := template.FuncMap{
		"derefInt": func(i *int) int {
			return *i
		},
	}

	if Cr.server_cfg_ini == nil {
		file := FindFile("/ini/server_cfg.ini")
		tmplStr, err := io.ReadAll(file)
		if err != nil {
			log.Print(err)
		}
		Cr.server_cfg_ini, err = template.New("server_cfg.ini").Funcs(funcMap).Parse(string(tmplStr))
		if err != nil {
			log.Print(err)
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
		"max_clients": max_clients,
		"sunangle":    cr.Time_To_SunAngle(tm.Time),
		"cspstr":      cspstr,
		"name":        event_cat.Name,
	}

	var b bytes.Buffer
	err := Cr.server_cfg_ini.Execute(&b, data)
	if err != nil {
		log.Print(err)
	}

	if Cr.entry_list_ini == nil {
		file := FindFile("/ini/entry_list.ini")
		tmplStr, err := io.ReadAll(file)
		if err != nil {
			log.Print(err)
		}
		Cr.entry_list_ini, err = template.New("entry_list.ini").Funcs(funcMap).Parse(string(tmplStr))
		if err != nil {
			log.Print(err)
		}
	}

	var b2 bytes.Buffer
	err = Cr.entry_list_ini.Execute(&b2, class)
	if err != nil {
		log.Print(err)
	}

	cr.ServerCfg_Result = b.String()
	cr.EntryList_Result = b2.String()
	cr.Class = class
	cr.Track = track
}

package sm

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
)

func Time_To_SunAngle(*string) int {
	return -80
}

func Render_Ini(event_id int) (string, string) {

	event := Dba.Select_Event(event_id)
	time := Dba.Select_Time_Weather(*event.Time_Id)
	diff := Dba.Select_Difficulty(*event.Difficulty_Id)
	cfg := Dba.Select_Config()
	session := Dba.Select_Session(*event.Session_Id)
	class := Dba.Select_Class_Entries(*event.Class_Id)
	track := Dba.Select_Cache_Track(*event.Cache_Track_Key, *event.Cache_Track_Config)

	data := map[string]any{
		"event":    event,
		"config":   cfg,
		"diff":     diff,
		"session":  session,
		"time":     time,
		"class":    class,
		"track":    track,
		"sunangle": Time_To_SunAngle(time.Time),
		"cspstr":   "",
	}

	// csp required? build a cspstr to be concat with the track name
	if cfg.Csp_Required != nil && *cfg.Csp_Required > 0 {
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

		data["cspstr"] = "csp/" + strconv.Itoa(*cfg.Csp_Version) + csp_letter + "/../"
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

	data["max_clients"] = max_clients

	// Strategy needed? re-order cars in the entry list as per selected strategy
	if strat_needed {
		// Random
		if *event.Strategy == 2 {
			rand.Shuffle(len(class.Entries), func(i, j int) { class.Entries[i], class.Entries[j] = class.Entries[j], class.Entries[i] })

			// Cut the list by the max number of clients
			class.Entries = class.Entries[:max_clients]
		}
	}


	funcMap := template.FuncMap{
		"derefInt": func (i *int) int {
			return *i
		},
	}
	var tmplFile = filepath.Join("ini", "server_cfg.ini")
	tmpl, err := template.New("server_cfg.ini").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		log.Print(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		log.Print(err)
	}

	return b.String(), ""
}

func Render_Entry_List(event_id int) string {
	event := Dba.Select_Event(event_id)
	class := Dba.Select_Class_Entries(*event.Class_Id)

	// TODO respect max clients

	var tmplFile = filepath.Join("ini", "entry_list.ini")
	tmpl, err := template.New("entry_list.ini").ParseFiles(tmplFile)
	if err != nil {
		log.Print(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, class)
	if err != nil {
		log.Print(err)
	}

	return b.String()
}

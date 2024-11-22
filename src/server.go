package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

type Server_Status struct {
	Status    bool   `json:"status"`
	Public_Ip string `json:"public_ip"`
	Players   int    `json:"players"`
	Session   SessionInfo
}

var Public_Ip string

func (stats Server_Status) Refresh() {
	Status.Status = Is_Running()
	Status.Public_Ip = Public_Ip

}

func (status Server_Status) Update_Public_Ip() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				res, err := http.Get("https://api.ipify.org")
				if err != nil {
					log.Print(err)
				}
				ip, err := io.ReadAll(res.Body)
				if err != nil {
					log.Print(err)
				}
				Public_Ip = string(ip)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (status Server_Status) Server_ApplyTrack() {
	next_events := Dba.Select_Server_Events(true)

	if len(next_events) < 0 {
		log.Print("No events in queue")
		return
	}
	next_event := next_events[0]

	Cr.Render_Ini(*next_event.User_Event.Id)
	Cr.Write_Ini()

	tm := time.Now().Unix()
	next_event.Started_At = &tm

	next_event.ServerCfg = &Cr.ServerCfg_Result
	next_event.EntryList = &Cr.EntryList_Result

	Dba.Update_Server_Event(next_event)

	// Populate TempFolder
	exec := "acServer"
	if runtime.GOOS == "windows" {
		exec = "acServer.exe"
	}
	Zf.ExtractFile(Zf.FindZipFile(exec), TempFolder)

	// Extract cars
	for _, e := range Cr.Class.Entries {
		Zf.ExtractFiles(Zf.FindZipFiles("cars/"+*e.Cache_Car_Key+"/"), filepath.Join(TempFolder, "content"))
	}

	// Extract track
	if *Cr.Track.Config == "" {
		Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/models.ini"), filepath.Join(TempFolder, "content"))
		Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/data/drs_zones.ini"), filepath.Join(TempFolder, "content"))

		if Cr.Csp_Required {
			Zf.ExtractFileToSubfolder(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/data/surfaces.ini"), filepath.Join(TempFolder, "content", "tracks", "csp", *Cr.Track.Key, *Cr.Track.Config, "data"))
		} else {
			Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/data/surfaces.ini"), filepath.Join(TempFolder, "content"))
		}
	} else {
		Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/models_"+*Cr.Track.Config+".ini"), filepath.Join(TempFolder, "content"))
		Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/"+*Cr.Track.Config+"/data/drs_zones.ini"), filepath.Join(TempFolder, "content"))

		if Cr.Csp_Required {
			Zf.ExtractFileToSubfolder(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/"+*Cr.Track.Config+"/data/surfaces.ini"), filepath.Join(TempFolder, "content", "tracks", "csp", *Cr.Track.Key, *Cr.Track.Config, "data"))
		} else {
			Zf.ExtractFile(Zf.FindZipFile("tracks/"+*Cr.Track.Key+"/"+*Cr.Track.Config+"/data/surfaces.ini"), filepath.Join(TempFolder, "content"))
		}
	}

	// Extract surfaces.ini
	Zf.ExtractFile(Zf.FindZipFile("system/data/surfaces.ini"), filepath.Join(TempFolder))

	Zf.Close()

}

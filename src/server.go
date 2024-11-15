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
	Dba.Update_Event_SetComplete()
	next_event := Dba.Select_Event_Next()
	Cr.Render_Ini(*next_event.Id)
	Cr.Write_Ini()

	tm := time.Now().Unix()
	next_event.ServerCfg = &Cr.ServerCfg_Result
	next_event.EntryList = &Cr.EntryList_Result

	next_event.Started_At = &tm

	//Dba.Update_Event(next_event)

	// Populate TempFolder
	exec := "acServer"
	if runtime.GOOS == "windows" {
		exec = "acServer.exe"
	}
	Zf.ExtractFileToFolder(Zf.FindZipFile(exec), TempFolder)

	// Extract cars
	for _, e := range Cr.Class.Entries {
		Zf.ExtractFilesToFolder(Zf.FindZipFiles("cars/"+*e.Cache_Car_Key+"/"), filepath.Join(TempFolder, "content"))
	}

	// Extract track
	// TODO: CSP compensation
	Zf.ExtractFilesToFolder(Zf.FindZipFiles("tracks/"+*Cr.Track.Key+"/"), filepath.Join(TempFolder, "content"))

	// Extract surfaces.ini
	Zf.ExtractFileToFolder(Zf.FindZipFile("/system/data/surfaces.ini"), filepath.Join(TempFolder, "system", "data", "surfaces.ini"))

	Zf.Close()

}

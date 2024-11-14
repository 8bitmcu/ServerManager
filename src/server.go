package main

import (
	"io"
	"log"
	"net/http"
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
	cr := Cr.Render_Ini(*next_event.Id)
	cr.Write_Ini()

	tm := time.Now().Unix()
	next_event.ServerCfg = &cr.ServerCfg_Result
	next_event.EntryList = &cr.EntryList_Result

	next_event.Started_At = &tm

	Dba.Update_Event(next_event)
}

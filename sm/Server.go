package sm

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Server_Stats struct {
	Status    bool   `json:"status"`
	Public_Ip string `json:"public_ip"`
}

func (stats Server_Stats) Refresh() {
	Stats.Status = Is_Running()
}

func (stats Server_Stats) Update_Public_Ip() {
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
				if len(string(ip)) > 0 {
					stats.Public_Ip = string(ip)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

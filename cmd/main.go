package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"campushub-live-data/internal/stream"
	"campushub-live-data/internal/ws"
)

func main() {
	fmt.Println("CampusHub Live Data Service başlatılıyor...")

	// 1. Sağlık Kontrolü ve WebSocket Uç Noktaları
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Servis ayakta ve veri akisina hazir!"))
	})
	http.HandleFunc("/ws", ws.ServeWS)

	// 2. İrem'in veri akışı için kanalını (Channel) oluştur
	eventChan := make(chan stream.Event)

	// 3. İrem'in Veri Akış Motorunu Arka Planda Başlat (Saniyede 5 veri)
	go stream.StartDataStream("data/events_converted.csv", 5, eventChan)

	// 4. SİSTEMLERİN BİRLEŞTİĞİ NOKTA
	// İrem'den gelen ham veriyi alıp JSON'a çeviriyoruz ve senin WebSocket Havuzuna basıyoruz
	go func() {
		for e := range eventChan {
			// Terminalden akışı takip etmek için log
			fmt.Printf("⚡ CANLI: %s kullanıcısı %s ürününü %s yaptı.\n", e.UserID, e.EventID, e.Action)

			// Go Struct yapısını web tarayıcılarının anlayacağı JSON formatına dönüştür
			veriJSON, err := json.Marshal(e)
			if err != nil {
				fmt.Println("JSON çevirme hatası:", err)
				continue
			}

			// Begüm'ün havuzuna gönder (Bağlı olan herkese yayınla)
			ws.BroadcastMessage(veriJSON)
		}
	}()

	// Sunucuyu ayağa kaldır
	port := ":8080"
	fmt.Printf("Sunucu %s portunda dinleniyor...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}

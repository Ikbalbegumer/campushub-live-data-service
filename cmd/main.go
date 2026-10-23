package main

import (
	"campushub-live-data/internal/stream"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("CampusHub Live Data Service başlatılıyor...")

	// 1. Sağlık Kontrolü Uç Noktası
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Servis ayakta!"))
	})

	// 2. Veri akışı için bir kanal (Channel) oluştur
	eventChan := make(chan stream.Event)

	// 3. İrem'in Veri Akış Motorunu Arka Planda (Goroutine) Başlat
	// Saniyede 5 veri fırlatacak şekilde ayarladık
	go stream.StartDataStream("data/events_converted.csv", 5, eventChan)

	// 4. Akıp gelen veriyi ekranda canlı izle (İleride Begüm bu veriyi WebSocket'e basacak)
	go func() {
		for e := range eventChan {
			fmt.Printf("⚡ CANLI: %s kullanıcısı %s ürününü %s yaptı.\n", e.UserID, e.EventID, e.Action)
		}
	}()

	// Sunucuyu ayağa kaldır
	port := ":8080"
	fmt.Printf("Sunucu %s portunda dinleniyor...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
package main

import (
	"campushub-live-data/internal/ws"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("CampusHub Live Data Service başlatılıyor...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Servis ayakta ve veri akisina hazir!"))
	})

	http.HandleFunc("/ws", ws.ServeWS)

	port := ":8080"
	fmt.Printf("Sunucu %s portunda dinleniyor...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"campushub-live-data/internal/stream"

	"github.com/gorilla/websocket"
)

// Tarayıcı bağlantıları için WebSocket ayarı (CORS engeline takılmamak için Origin kontrolünü kapatıyoruz)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Global olarak bağlı olan istemcileri tutmak için mekanizma
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan stream.Event)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// HTTP isteğini WebSocket bağlantısına yükselt
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade hatası: %v", err)
		return
	}
	defer ws.Close()

	// Yeni istemciyi kaydet
	clients[ws] = true
	log.Printf("Yeni bir tarayıcı istemcisi bağlandı!")

	for {
		// Bağlantının kopmadığını anlamak için okuma döngüsü
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			log.Printf("İstemci bağlantısı kesildi.")
			break
		}
	}
}

// Kanaldan gelen verileri bağlı tüm tarayıcılara WebSocket üzerinden yayınlayan fonksiyon
func handleMessages() {
	for {
		event := <-broadcast
		// Her bir bağlı istemciye JSON formatında veriyi gönder
		for client := range clients {
			err := client.WriteJSON(event)
			if err != nil {
				log.Printf("WebSocket yazma hatası: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	fmt.Println("CampusHub Live Data Service başlatılıyor...")

	// 1. Sağlık Kontrolü Uç Noktası
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Servis ayakta ve veri akisina hazir!"))
	})

	// 2. WebSocket Bağlantı Noktası (Arayüzün bağlandığı adres)
	http.HandleFunc("/ws", handleConnections)

	// 3. Veri akışı için bir kanal (Channel) oluştur
	eventChan := make(chan stream.Event)

	// 4. İrem'in Veri Akış Motorunu Arka Planda Başlat (Saniyede 5 veri)
	go stream.StartDataStream("data/events_converted.csv", 5, eventChan)

	// 5. SİSTEMLERİN BİRLEŞTİĞİ NOKTA
	// Kanaldan gelen verileri yakalayıp hem terminale bas hem de broadcast kanalına aktar
	go func() {
		for e := range eventChan {
			// Terminalden akışı takip etmek için log
			fmt.Printf("⚡ CANLI: %s kullanıcısı %s ürününü %s yaptı.\n", e.UserID, e.EventID, e.Action)

			// WebSocket'e bağlı olan tüm tarayıcılara fırlat
			broadcast <- e
		}
	}()

	// Mesaj yayın mekanizmasını arka planda çalıştır
	go handleMessages()

	// Sunucuyu 8080 portunda ayağa kaldır
	port := ":8080"
	fmt.Printf("Sunucu %s portunda dinleniyor...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}

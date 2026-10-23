package stream

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

// İleride Beyza'nın modellerine geçeceğiz, şimdilik motorun kendi basit yapısı:
type Event struct {
	Timestamp string
	UserID    string
	EventID   string
	Action    string
}

// StartDataStream CSV'yi okur ve kanala (channel) saniyede 'rate' kadar veri pompalar
func StartDataStream(filepath string, ratePerSecond int, outChan chan<- Event) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Hata: Veri dosyası açılamadı: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Başlık satırını atla

	// Saniyede gönderilecek veri hızını (Rate Limit) ayarlayan kronometre
	delay := time.Duration(1000/ratePerSecond) * time.Millisecond
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	fmt.Printf("🚀 Veri akışı saniyede %d kayıt hızıyla başlatıldı...\n", ratePerSecond)

	for range ticker.C {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Veri dosyasının sonuna gelindi. Akış tamamlandı.")
			close(outChan) // Kanalı kapat
			break
		}
		if err != nil {
			continue // Hatalı satırı atla
		}

		// CSV satırını alıp Event objesine çevir
		event := Event{
			Timestamp: record[0],
			UserID:    record[1],
			EventID:   record[2],
			Action:    record[3],
		}

		// Veriyi anlık olarak kanala fırlat
		outChan <- event
	}
}
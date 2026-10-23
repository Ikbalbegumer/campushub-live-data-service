package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	// 5.28 GB'lık devasa kaynak dosyayı bellek dostu şekilde açıyoruz
	inputFile, err := os.Open("data/2019-Oct.csv")
	if err != nil {
		fmt.Println("Hata: Kaynak dosya bulunamadı. Lütfen data/2019-Oct.csv yolunu kontrol et.", err)
		return
	}
	defer inputFile.Close()

	// Oluşturulacak yeni ve temiz dosya
	outputFile, err := os.Create("data/events_converted.csv")
	if err != nil {
		fmt.Println("Hata: Çıktı dosyası oluşturulamadı.", err)
		return
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Projenin istediği zorunlu başlıkları yazıyoruz
	writer.Write([]string{"timestamp", "userId", "eventId", "action"})

	// Orijinal dosyanın ilk satırını (eski başlıkları) atlıyoruz
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Başlık okuma hatası:", err)
		return
	}

	fmt.Println("Dönüştürme işlemi başladı. Sadece 1 Milyon satır alınacak...")

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Bozuk veya hatalı satırları atla
		}

		// 1. Aksiyon (Action) Eşleştirmesi
		action := "VIEW"
		switch record[1] {
		case "view":
			action = "VIEW"
		case "cart":
			action = "CLICK"
		case "purchase":
			action = "REGISTER"
		case "remove_from_cart":
			action = "CANCEL"
		}

		// 2. Zaman (Timestamp) Dönüştürmesi (2019-10-01 00:00:00 UTC -> Unix Timestamp)
		t, err := time.Parse("2006-01-02 15:04:05 MST", record[0])
		timestamp := record[0] // Hata olursa orijinali kalsın
		if err == nil {
			timestamp = fmt.Sprintf("%d", t.Unix())
		}

		// 3. ID Atamaları
		userId := "user_" + record[7]
		eventId := "event_" + record[2]

		// Yeni satırı yazdır
		writer.Write([]string{timestamp, userId, eventId, action})

		count++
		
		// Terminalde süreci takip etmek için her 100 binde bir bilgi ver
		if count%100000 == 0 {
			fmt.Printf("%d satır başarıyla işlendi...\n", count)
		}

		// Proje gereksinimi için 1 Milyon satıra ulaşıldığında döngüyü bitir
		if count >= 1000000 {
			break
		}
	}
	
	fmt.Printf("İşlem tamamlandı! Tertemiz %d satırlık veri data/events_converted.csv olarak kaydedildi.\n", count)
}
# CampusHub Live Data Service (Grup 6)

CampusHub projesi kapsamında gerçek zamanlı kullanıcı hareketlerini (clickstream) simüle eden ve WebSocket üzerinden yayınlayan Live Data Service.

Bu servis, e-ticaret platformlarındaki anlık kullanıcı etkileşimlerini yüksek performanslı bir veri akış motoru ile işler ve istemcilere canlı olarak iletir.

## 🚀 Özellikler

- **Gerçek Zamanlı Veri Akışı:** Arka planda Go (Golang) motoru kullanılarak büyük veri setleri bellek dostu (streaming) bir yaklaşımla işlenir.
- **WebSocket Yayını:** `gorilla/websocket` altyapısı sayesinde istemcilere kesintisiz yayın (broadcast) yapılır.
- **Canlı Test Arayüzü:** HTML/JS ile geliştirilmiş kontrol paneli üzerinden veriler anlık olarak izlenebilir.
- **İstatistik ve Filtreleme:** Arayüz üzerinde anlık olarak "Toplam İşlem" ve "Benzersiz Kullanıcı" sayıları hesaplanır. Ayrıca gelen veri akışı "Sadece VIEW" veya "Sadece PURCHASE" olarak filtrelenebilir.

## 🛠️ Kullanılan Teknolojiler

- **Backend:** Go (Golang)
- **Frontend:** HTML, CSS, Vanilla JavaScript
- **Haberleşme:** WebSocket (`gorilla/websocket`)

## 📂 Veri Seti Hakkında Önemli Not

Projenin temelini oluşturan e-ticaret log veri seti (`2019-Oct.csv`) yaklaşık 1.31 GB boyutundadır. GitHub'ın dosya boyutu limitlerinden dolayı bu dosya Git takibinden çıkartılmış ve depoya eklenmemiştir (`.gitignore`). Projeyi yerelde çalıştırabilmek için ilgili CSV dosyasının `data/` klasörü altına eklenmesi gerekmektedir.

## 💻 Kurulum ve Çalıştırma

Projeyi yerel ortamda çalıştırmak için şu adımları izleyin:

```bash
# Bağımlılıkları indirin
go mod tidy

# Sunucuyu başlatın
go run cmd/main.go
```

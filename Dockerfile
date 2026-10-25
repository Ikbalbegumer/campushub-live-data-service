# Go resmi imajını temel al
FROM golang:1.22-alpine

# Çalışma dizinini ayarla
WORKDIR /app

# Bağımlılık dosyalarını kopyala
COPY go.mod go.sum ./
RUN go mod download

# Tüm proje dosyalarını içeri aktar
COPY . .

# Go uygulamasını derle
RUN go build -o main cmd/main.go

# Sunucunun çalışacağı portu belirt
EXPOSE 8080

# Konteyner ayağa kalkınca çalışacak komut
CMD ["./main"]
# Gunakan image Go berbasis Alpine
FROM golang:alpine

# Set working directory
WORKDIR /app

# Salin file go mod dan sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi
RUN go build -o main

# Setting Permission
RUN chmod +x main

# Jalankan aplikasi
CMD ["./main -c 1"]
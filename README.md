# book-stock-manager

## WORK IN PROGRESS

Backend REST API untuk manajemen stok buku.

Proyek pembuatan **Sistem Managemen Stok dan Penjualan Berbasis QR Code** yang diadakan pada acara *BeSmart* oleh UKM Intermedia Universitas Amikom Purwokerto.

Dokumentasi API dapat diakses [disini](https://crazydw4rf.github.io/book-stock-manager).

## Cara Penggunaan

### Setup Config

1. Copy file config contoh ke file yang akan digunakan
   ```bash
   cp .env.example .env
   ```

   Note: App tetap bisa jalan tanpa file `.env` jika environment variable sudah diset di sistem.

2. Setting nilai di file `.env` sesuai kebutuhan
   ```bash
   # App config
   APP_HOST=127.0.0.1
   APP_PORT=8080

   # JWT config
   JWT_ACCESS_TOKEN_SECRET=secret_key_here
   JWT_REFRESH_TOKEN_SECRET=secret_key_here

   # Database config
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=username
   DB_PASSWORD=password
   DB_NAME=book_stock
   ```

### Database Migration

Jalankan migration untuk membuat tabel yang diperlukan:

```bash
# Di root folder project
go run db/migrate.go db/migrations up
```

### Build App (Linux)

```bash
# Simpan nama modul ke variabel untuk memudahkan
export MODULE="github.com/crazydw4rf/book-stock-manager"

# Build app (debug mode)
go build -o book-stock-manager ./cmd/app/main.go

# Build with custom version
go build -ldflags="-X ${MODULE}/internal/config.APP_VERSION=1.0.0" -o book-stock-manager ./cmd/app/main.go

# Build for production (release mode)
go build -ldflags="-s -w -X ${MODULE}/internal/config.APP_VERSION=1.0.0 -X ${MODULE}/internal/config.APP_ENV=production" -o book-stock-manager ./cmd/app/main.go
```

### Run App (Linux)

```bash
# Run dari hasil build
./book-stock-manager

# Atau langsung run
go run ./cmd/app/main.go
```

Setelah app running, Swagger UI untuk dokumentasi API bisa diakses di: `http://localhost:8080/docs/`

## TODO
- [ ] Menambahkan file aksi CI/CD untuk otomatisasi proses build, test, dan deployment.
- [ ] Menambahkan unit test dan integration test.
- [ ] Mengimplementasikan fitur autentikasi pengguna.
- [ ] Memperbaiki dokumentasi.

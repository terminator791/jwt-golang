# README - Sistem E-Ticketing Transportasi Publik

## 📌 Deskripsi Proyek

Sistem E-Ticketing Transportasi Publik adalah aplikasi backend berbasis Golang yang dirancang untuk mengelola sistem tiket elektronik pada jaringan transportasi publik. Aplikasi ini mendukung operasi 24 jam dengan fitur offline capability, sehingga tetap dapat beroperasi meskipun koneksi internet terputus.

Sistem ini menangani:
- Manajemen kartu prepaid untuk pengguna
- Transaksi check-in dan check-out di terminal
- Perhitungan tarif berdasarkan terminal asal dan tujuan
- Sinkronisasi data antara terminal dan server pusat
- Manajemen pengguna dan autentikasi dengan JWT

## 🏗️ Arsitektur Sistem

Sistem E-Ticketing menggunakan arsitektur hybrid yang menggabungkan:
1. **Edge Computing** - Server lokal di setiap terminal
2. **Centralized Microservices** - Layanan terpusat untuk manajemen data
3. **API Gateway** - Entry point untuk interaksi dengan sistem
4. **Database PostgreSQL** - Penyimpanan data persisten
5. **JWT Authentication** - Sistem otentikasi berbasis token

## 🚀 Fitur Utama

- **Manajemen Pengguna**
  - Registrasi pengguna baru
  - Login dengan JWT authentication
  - Melihat profil pengguna
  - Logout (blacklist token)

- **Manajemen Terminal & Gate**
  - Membuat terminal baru
  - Mengelola gate di terminal

- **Keamanan**
  - Rate limiting untuk mencegah brute force attack
  - Validasi password yang kuat
  - Perlindungan terhadap XSS
  - Token blacklisting untuk logout yang aman

- **Offline Capability**
  - Operasi check-in/check-out tetap berjalan saat offline
  - Sinkronisasi otomatis saat koneksi tersedia kembali

## 🛠️ Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go (Golang) v1.21+
- **Framework Web**: Gin Gonic
- **ORM Database**: GORM
- **Database**: PostgreSQL
- **Autentikasi**: JWT (JSON Web Token)
- **Enkripsi**: bcrypt

## 📋 Persyaratan Sistem

- Go (Golang) versi 1.21 atau lebih baru
- PostgreSQL versi 13 atau lebih baru
- Terminal/Command Line dengan Git

## 📥 Cara Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/IqbalBPH/golang-e-ticketing.git
cd golang-e-ticketing
```

### 2. Instal Dependensi

```bash
go mod download
```

### 3. Siapkan Database

Buat database PostgreSQL baru untuk aplikasi:

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database baru
CREATE DATABASE e_ticketing;

# Keluar dari psql
\q
```

### 4. Konfigurasi Lingkungan

Buat file `.env` di root directory:

```bash
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=e_ticketing
DB_PORT=5432

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET_KEY=rahasia_jwt_yang_sangat_aman_dan_panjang
```

### 5. Jalankan Migrasi Database

```bash
go run main.go
```

Aplikasi akan secara otomatis membuat skema tabel dan user admin default:

```
Email: admin@e-ticketing.com
Password: admin123
```

> ⚠️ **Penting**: Segera ubah password admin default setelah instalasi!

## 🚆 Cara Menjalankan Aplikasi

```bash
go run main.go
```

Atau build dan jalankan executable:

```bash
go build -o e-ticketing
./e-ticketing
```

Server akan berjalan di `http://localhost:8080`

## 🌐 Dokumentasi API Endpoint

### 🔑 Autentikasi

#### 1. Register - Mendaftar Pengguna Baru

**Endpoint**: `POST /api/auth/register`

**Request Body**:
```json
{
  "full_name": "Budi Santoso",
  "email": "budi@example.com",
  "password": "Rahasia123",
  "phone": "081234567890",
  "date_of_birth": "1990-01-15T00:00:00Z"
}
```

**Response (201 Created)**:
```json
{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "full_name": "Budi Santoso",
    "email": "budi@example.com",
    "phone": "081234567890",
    "user_type": "CUSTOMER"
  }
}
```

**Validasi Password**:
- Minimal 6 karakter
- Minimal 1 huruf besar
- Minimal 1 huruf kecil
- Minimal 1 angka

#### 2. Login - Masuk ke Sistem

**Endpoint**: `POST /api/auth/login`

**Request Body**:
```json
{
  "email": "budi@example.com",
  "password": "Rahasia123"
}
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "full_name": "Budi Santoso",
      "email": "budi@example.com",
      "user_type": "CUSTOMER"
    }
  }
}
```

#### 3. Cek Profil - Melihat Informasi Akun

**Endpoint**: `GET /api/user/profile`

**Header**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Profil user berhasil didapatkan",
  "data": {
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "full_name": "Budi Santoso",
    "email": "budi@example.com",
    "phone": "081234567890",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "user_type": "CUSTOMER",
    "created_at": "2023-07-07T15:30:45.123456Z"
  }
}
```

#### 4. Logout - Keluar dari Sistem

**Endpoint**: `POST /api/auth/logout`

**Header**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "Logout berhasil"
}
```

### 🏢 Manajemen Terminal

#### 1. Buat Terminal - Membuat Terminal Baru

**Endpoint**: `POST /api/terminal`

**Header**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Request Body**:
```json
{
  "terminal_name": "Terminal Pusat",
  "terminal_code": "TPT01",
  "location": "Jakarta Pusat",
  "latitude": -6.175110,
  "longitude": 106.865036,
  "is_active": true
}
```

**Response (201 Created)**:
```json
{
  "status": "success",
  "message": "Terminal berhasil dibuat",
  "data": {
    "terminal_id": 1,
    "terminal_name": "Terminal Pusat",
    "terminal_code": "TPT01",
    "location": "Jakarta Pusat",
    "latitude": -6.175110,
    "longitude": 106.865036,
    "is_active": true,
    "created_at": "2025-07-07T14:30:45.123456Z",
    "updated_at": "2025-07-07T14:30:45.123456Z"
  }
}
```

## 🛡️ Fitur Keamanan

### Rate Limiting

API dilindungi dari serangan brute force dengan rate limiting:
- Endpoint autentikasi: 5 request per menit
- Endpoint umum: 30 request per menit

### Validasi Password

Password harus memenuhi kriteria keamanan minimum:
- Panjang minimal 6 karakter
- Mengandung huruf besar dan kecil
- Mengandung angka

### Perlindungan XSS

Semua input disanitasi untuk mencegah serangan XSS (Cross-Site Scripting).

## 📚 Skema Database

Sistem menggunakan skema database yang mencakup 9 tabel utama:

1. **USER** - Informasi pengguna
2. **CARD** - Data kartu prepaid
3. **TERMINAL** - Informasi terminal transportasi
4. **GATE** - Data gate di setiap terminal
5. **TRANSACTION** - Rekaman transaksi perjalanan
6. **FARE_MATRIX** - Matriks tarif antar terminal
7. **CARD_BALANCE_LOG** - Log perubahan saldo kartu
8. **TOP_UP** - Transaksi pengisian saldo
9. **SYNC_LOG** - Log sinkronisasi data

## 🔍 Best Practices & Pertimbangan

### Arsitektur Clean
Proyek menggunakan pendekatan clean architecture dengan pemisahan layer:
- Controllers: Menangani HTTP request/response
- Services: Berisi logika bisnis
- Models: Representasi struktur data
- Middleware: Komponen untuk menangani proses sebelum/sesudah request

### Keamanan
- Token JWT untuk autentikasi
- Password di-hash menggunakan bcrypt
- Rate limiting untuk mencegah brute force
- Sanitasi input untuk mencegah XSS
- Blacklist token untuk logout yang aman

### Database
- UUID untuk primary key tabel utama (mencegah collision di lingkungan terdistribusi)
- Pemilihan tipe data yang tepat (decimal untuk nilai keuangan)
- Penggunaan indeks untuk optimasi query
- Hubungan relasional yang jelas antar entitas

## 🤝 Kontribusi

Silakan berkontribusi pada proyek ini dengan mengikuti langkah-langkah berikut:
1. Fork repository
2. Buat branch fitur baru (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Menambahkan fitur amazing'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📃 Lisensi

Proyek ini dilisensikan di bawah Lisensi MIT - lihat file LICENSE untuk detail.

---

Dibuat dengan ❤️ oleh Tim E-Ticketing | © 2025 | [GitHub](https://github.com/IqbalBPH)
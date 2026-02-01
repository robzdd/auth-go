# Auth System (React + Golang)

Sistem autentikasi fullstack dengan React (Frontend) dan Golang (Backend).

## Prasyarat

- Go 1.20+
- Node.js 18+
- MySQL 8.0+

## Setup & Run

### 1. Database Setup

Buat database di MySQL:

```sql
CREATE DATABASE auth_go;
```

### 2. Backend (Golang)

1.  Masuk ke folder backend:
    ```bash
    cd backend
    ```
2.  Setup Environment Variables:
    - Buka file `.env`
    - Isi `DB_PASSWORD` (password MySQL Anda)
    - Isi Gmail Credentials (`SMTP_EMAIL` dan `SMTP_PASSWORD`)
      - _Note: Gunakan App Password dari Google Account, bukan password login biasa._
3.  Jalankan server:
    ```bash
    go run cmd/api/main.go
    ```
    server akan berjalan di port `8080`.

### 3. Frontend (React)

1.  Masuk ke folder frontend:
    ```bash
    cd frontend
    ```
2.  Install dependencies (jika belum):
    ```bash
    npm install
    ```
3.  Jalankan server development:
    ```bash
    npm run dev
    ```
    Aplikasi bisa diakses di `http://localhost:5173`.

## Struktur Project

- **Backend**: Clean Architecture (`cmd`, `internal`, `pkg`)
- **Frontend**: Feature-Based (`features/auth`, `features/dashboard`)

## Fitur

- Login & Register (JWT)
- Dashboard Terproteksi
- Forgot Password (Email OTP link)
- Reset Password

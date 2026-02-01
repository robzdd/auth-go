# Dokumentasi Arsitektur Project (Auth-Go)

Dokumen ini menjelaskan struktur folder, alur data, dan koneksi antar komponen dalam sistem Fullstack ini.

## 1. Struktur Folder

### A. Backend (`/backend`)

Menggunakan prinsip **Clean Architecture** (atau Standard Go Project Layout). Tujuannya adalah memisahkan "Logic Bisnis" dari "Framework/Database" agar kode mudah dites dan diganti-ganti.

- **`cmd/api/main.go`**: Titik awal (Entry Point). Di sini kita load config, connect database, dan menyambungkan semua komponen (Repo -> Service -> Handler).
- **`internal/`**: Kode inti aplikasi yang tidak boleh di-import oleh project lain.
  - **`config/`**: Membaca file `.env` (Database credentials, JWT secret).
  - **`domain/`**: "Jantung" aplikasi. Berisi **Struct** (Model User) dan **Interface** (Kontrak fungsi). Ini tidak tergantung pada library apapun.
  - **`repository/`**: Layer akses data (Database). Hanya di folder ini kita menyentuh SQL/GORM.
  - **`service/`**: Layer logika bisnis. Contoh: Hashing password sebelum simpan, validasi input, kirim email. Service tidak tahu soal HTTP atau SQL, dia cuma tahu logic.
  - **`handler/`**: Layer transportasi HTTP (menggunakan Gin). Tugasnya baca Request Body (JSON), panggil Service, dan balikin Response JSON.
  - **`middleware/`**: Pengecekan di tengah jalan (contoh: Cek token JWT sebelum masuk handler).
- **`pkg/`**: Library bantuan (Helper) yang bisa dipakai ulang (contoh: fungsi JWT, Hashing Password).

### B. Frontend (`/frontend`)

Menggunakan **Feature-Based Architecture**. Kode dikelompokkan berdasarkan "Fitur", bukan berdasarkan jenis file.

- **`src/features/`**: Folder utama fitur.
  - **`auth/`**: Semua tentang login/register (Component Form, API hook, Types).
  - **`dashboard/`**: Semua tentang dashboard (Tabel user, API hook user listing).
- **`src/components/ui/`**: Komponen visual reusable (Button, Input, Table) dari Shadcn UI.
- **`src/lib/`**: Konfigurasi global (Axios client, Utility class).
- **`src/App.tsx`**: Mengatur Routing (Halaman mana loading komponen apa).

---

## 2. Alur Data (Flow)

Bagaimana data mengalir dari klik tombol di Frontend sampai tersimpan di Database?

**Contoh Kasus: User Login**

1.  **Frontend (UI)**:
    - User isi form di `LoginForm.tsx`.
    - Tombol diklik -> memanggil `useLogin` (React Query Hook).
    - `axios` mengirim request HTTP POST ke `http://localhost:8080/api/auth/login`.

2.  **Backend (Router & Handler)**:
    - `main.go` menerima request dan mengarahkannya ke `AuthHandler.Login` (di `internal/handler`).
    - `AuthHandler` mem-parsing JSON body (email, password) dari request.

3.  **Backend (Service)**:
    - Handler memanggil `AuthService.Login` (di `internal/service`).
    - Service memvalidasi input.
    - Service memanggil `UserRepository.FindByEmail` untuk cari user.

4.  **Backend (Repository & DB)**:
    - `UserRepository` (di `internal/repository`) menyusun query SQL via GORM: `SELECT * FROM users WHERE email = ?`.
    - MySQL mengembalikan data User.

5.  **Backend (Balik ke Service)**:
    - Service menerima data user.
    - Service membandingkan password input vs password hash di DB (pakai `bcrypt`).
    - Jika cocok, Service membuat **JWT Token**.

6.  **Response**:
    - Handler mengirim JSON `{ "token": "ey...", "user": { ... } }` ke Frontend.
    - Frontend menyimpan token dan user redirect ke Dashboard.

---

## 3. Koneksi Database

Koneksi tidak terjadi secara "sihir", tapi melalui konfigurasi eksplisit:

1.  File **`.env`** di backend menyimpan kredensial:
    ```
    DB_USER=root
    DB_PASSWORD=...
    DB_NAME=auth_go
    ```
2.  Saat aplikasi start (`cmd/api/main.go`), fungsi `config.LoadConfig()` membaca file ini.
3.  Fungsi `database.ConnectDB(cfg)` (di `internal/database`) menggunakan driver `gorm.io/driver/mysql` untuk membuka koneksi TCP ke MySQL.
4.  Variabel `db` ini kemudian "disuntikkan" (Dependency Injection) ke dalam Repository.

Jadi `Repo` punya `DB`, `Service` punya `Repo`, dan `Handler` punya `Service`. Rantai ketergantungan ini dibangun di `main.go`.

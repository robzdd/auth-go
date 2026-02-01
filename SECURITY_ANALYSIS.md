# Analisis Security Project (Go + React vs Laravel)

Project ini dibangun menggunakan **Go (Backend)** dan **React (Frontend)**. Tidak seperti Laravel yang "baterai-included" (semua fitur security instan aktif), di Stack ini kita harus merakit pertahanan sendiri.

Berikut adalah analisis komprehensif kerentanan dan cara kita mengatasinya:

## 1. SQL Injection (Database Attack)

- **Laravel:** Dilindungi oleh Eloquent ORM.
- **Project Ini:** Menggunakan **GORM**.
  - ✅ **Aman:** Kami menggunakan Parameter Binding (`Where("name LIKE ?", input)`).
  - ❌ **Risiko:** Jika Anda menulis raw query seperti `db.Exec("SELECT * FROM users WHERE name = " + input)`, itu rentan. Di project ini, semua query sudah menggunakan cara aman GORM.

## 2. Cross-Site Scripting (XSS)

- **Laravel:** Blade engine otomatis melakukan escape string `{{ }}`.
- **Project Ini:** React.
  - ✅ **Aman:** React by default melakukan escaping pada data `{var}`.
  - ⚠️ **Risiko:** Kami menyimpan **JWT Token di LocalStorage**. Jika ada satu saja celah XSS (misal dari library pihak ketiga yang buruk), hacker bisa mencuri token ini.
  - **Mitigasi:** Kami menambahkan Header Keamanan (CSP) di backend untuk membatasi script apa yang boleh jalan.

## 3. Cross-Site Request Forgery (CSRF)

- **Laravel:** Punya `@csrf` token ditiap form.
- **Project Ini:** Stateless JWT.
  - ✅ **Analisis:** Karena kita menyimpan token di LocalStorage (bukan Cookie), aplikasi ini **KEBAL** terhadap CSRF standard. Browser tidak akan mengirim token secara otomatis ke domain lain.
  - Implikasinya: Kita "menukar" risiko CSRF dengan XSS. Dalam arsitektur modern SPA (Single Page App), ini adalah trade-off yang umum diterima.

## 4. Broken Object Level Authorization (IDOR)

- **Laravel:** Menggunakan Gates & Policies (`can('view', $post)`).
- **Project Ini:** Manual Check.
  - ⚠️ **Temuan:** Di endpoint `GetAllUsers`, saat ini **SEMUA user yang login** bisa melihat data semua user lain. Tidak ada pengecekan "Apakah saya Admin?".
  - **Rekomendasi:** Perlu menambahkan kolom `role` di tabel user dan middleware cek role.

## 5. Rate Limiting (Brute Force)

- **Laravel:** Middleware `throttle:api` (default 60 hit/menit).
- **Project Ini:** Belum ada.
  - ⚠️ **Risiko:** Hacker bisa mencoba login 1000x per detik sampai password ketemu.
  - **Rekomendasi:** Pasang library `gin-limiter`.

## 6. Information Disclosure (Penting!)

- **Masalah:** Sebelumnya, jika database error, backend mengirim pesan error mentah ke frontend (misal: "Table 'users' doesn't exist"). Ini memberi petunjuk ke hacker.
- **Perbaikan (Baru Saja Dilakukan):**
  - Endpoint `Forgot Password`: Sekarang selalu membalas "Link dikirim" walaupun email tidak ditemukan (Mencegah Email Enumeration).
  - Endpoint lainnya: Error server sekarang dibalas generic "Internal Server Error", detail aslinya hanya di-log di server.

## 7. Sensitive Data Exposure

- **Temuan:** File `.env` (Password DB) dan `api.exe` sempat ter-commit.
- **Perbaikan:** Sudah dibuatkan `.gitignore` dan file dihapus dari index Git.

---

**Kesimpulan:**
Secara fundamental arsitektur ini solid, tetapi karena Go menuntut kita menulis "plumbing" sendiri, risiko _human error_ (lupa validasi role, lupa sanitasi error) lebih besar daripada framework magic seperti Laravel.

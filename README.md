# EDSA - Go Fiber RESTful API

Proyek ini adalah REST API untuk melakukan tracking progres belajar siswa. Dibangun dengan Go menggunakan framework **Fiber** dan **GORM** untuk interaksi dengan database **MariaDB**, aplikasi ini dirancang dengan struktur yang mudah dipelihara dan siap untuk skala produksi, serta dilengkapi sistem autentikasi modern.

### Fitur Utama

* **Go Fiber**: Framework web yang cepat dan minimalis.
* **Struktur Terlapis**: Memisahkan kode ke dalam lapisan `handler`, `service`, dan `repository`.
* **Registrasi Asinkronus**: Menggunakan worker pool untuk memproses registrasi di latar belakang, membuat API lebih responsif.
* **Penanganan Kegagalan**: Memberikan notifikasi yang jelas jika proses registrasi di latar belakang gagal.
* **Autentikasi Aman dengan Paseto**: Menggunakan token Paseto sebagai alternatif yang lebih aman dari JWT.
* **Password Hashing dengan Argon2**: Mengamankan password pengguna dengan algoritma *hashing* modern.
* **GORM**: ORM yang tangguh untuk operasi database.
* **Migrasi & Seeder**: Menggunakan `golang-migrate` untuk manajemen skema dan *seeder* untuk data awal.
* **Validasi Terpusat**: Logika validasi *request* yang dapat digunakan kembali.
* **DTO (Data Transfer Object)**: Memisahkan model database dari respons API untuk format yang bersih.
* **Konfigurasi Berbasis `.env`**: Menggunakan variabel lingkungan untuk fleksibilitas.
* **Makefile**: Mengotomatiskan tugas-tugas pengembangan umum.

---

### Persyaratan

* [Go](https://go.dev/dl/) (versi 1.18+)
* [MariaDB](https://mariadb.org/download/)
* `migrate` CLI (diinstal dengan `go install`)

---

### Instalasi

1.  **Clone repositori:**
    ```bash
    git clone https://github.com/ipincamp/go-edsa-api.git
    cd go-edsa-api
    ```

2.  **Unduh dependensi Go:**
    ```bash
    go mod tidy
    ```

---

### Konfigurasi

Buat file `.env` di root proyek dan isi dengan konfigurasi berikut. **Pastikan untuk mengganti nilai-nilai placeholder.**

```dotenv
# Application
TZ=Asia/Jakarta
APP_PORT=8080
APP_ENV=development

# Database Connection
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=edsa_db

# Paseto Token Configuration
# HARUS TEPAT 32 KARAKTER
PASETO_SYMMETRIC_KEY="your-super-secret-key-must-be-32-bytes"
PASETO_EXPIRE_IN_HOURS=8

# Admin Seeder Credentials
AdminName=Admin
AdminEmail=admin@edsa.app
AdminPassword=secret
```

---

### Panduan Penggunaan

Gunakan `Makefile` untuk menjalankan perintah-perintah berikut dari terminal Anda:

| Perintah                               | Deskripsi                                                                        |
| -------------------------------------- | -------------------------------------------------------------------------------- |
| `make run`                             | Menjalankan aplikasi dalam mode pengembangan.                                    |
| `make build`                           | Membangun executable aplikasi.                                                   |
| `make run-prod`                        | Menjalankan aplikasi dalam mode produksi.                                        |
| `make new-migration NAME=nama_migrasi` | Membuat file migrasi baru. Contoh: `make new-migration NAME=create_users_table`. |
| `make migrate-up`                      | Menjalankan semua migrasi yang belum dieksekusi.                                 |
| `make migrate-down`                    | Mengembalikan (rollback) migrasi satu langkah.                                   |
| `make seed`                            | Menjalankan seeder database untuk menyuntikkan data awal (termasuk admin).       |
| `make clean`                           | Membersihkan file cache dan file sementara.                                      |

---

### Endpoint API

#### Autentikasi

* **`POST /api/auth/register`**: Mendaftarkan pengguna baru. Proses ini berjalan secara asinkronus. Anda akan menerima respons 202 Accepted yang menandakan permintaan Anda sedang diproses.
* **`POST /api/auth/login`**: Login untuk mendapatkan token Paseto.
* **`POST /api/auth/logout`**: Logout pengguna (memerlukan token).

#### Pengguna (Memerlukan Autentikasi)

* **`GET /api/profile`**: Mendapatkan detail profil pengguna yang sedang login.

---

### Struktur Proyek

```text
.
├── cmd/
│   ├── app/
│   │   └── main.go       # Titik masuk aplikasi utama
│   └── seeder/
│       └── main.go       # Titik masuk untuk seeder database
├── internal/
│   ├── api/
│   │   ├── dto/          # Data Transfer Objects (Request & Response)
│   │   ├── handlers/     # Berisi fungsi-fungsi handler API
│   │   ├── middleware/   # Middleware (e.g., autentikasi)
│   │   ├── response/     # Helper untuk respons JSON standar
│   │   ├── routes/       # Mendefinisikan endpoint API
│   │   └── validator/    # Logika validasi
│   ├── config/           # Mengatur konfigurasi aplikasi (.env)
│   ├── database/         # Koneksi database dan seeder
│   ├── models/           # Definisi model data
│   ├── repositories/     # Logika interaksi dengan database
│   └── services/         # Logika bisnis utama
│   └── utils/            # Utilitas (hashing, token)
├── migrations/           # File-file SQL untuk migrasi
├── scripts/              # Skrip pembantu (misal: untuk migrasi)
├── .env                  # Variabel lingkungan
├── .env.example          # Contoh variabel lingkungan
├── .gitignore
├── go.mod
├── go.sum
└── Makefile              # Skrip untuk otomatisasi

```

## LICENSE

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

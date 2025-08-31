# EDSA - Go Fiber RESTful API

Proyek ini adalah REST API untuk melakukan tracking progres belajar siswa. Dibangun dengan Go menggunakan framework **Fiber** dan **GORM** untuk interaksi dengan database **MariaDB**, aplikasi ini dirancang dengan struktur yang mudah dipelihara dan siap untuk skala produksi, serta dilengkapi sistem autentikasi modern.

### Fitur Utama

* **Go Fiber**: Framework web yang cepat dan minimalis.
* **Struktur Terlapis**: Memisahkan kode ke dalam lapisan `handler`, `service`, dan `repository`.
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
ENV=development
SERVER_HOST=localhost
SERVER_PORT=8000

# Database Connection (PostgreSQL)
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=root
DB_PASS=password
DB_NAME=edsa_db
DB_TZ=Asia/Jakarta

# Paseto Token Configuration
# HARUS TEPAT 32 KARAKTER
PASETO_SECRET_KEY="your-super-secret-key-must-be-32-bytes"
PASETO_TOKEN_TTL_MIN=480

# Admin Seeder Credentials
ADMIN_NAME=Admin
ADMIN_EMAIL=admin@edsa.app
ADMIN_PASSWORD=secret
```

---

### Panduan Penggunaan

Gunakan `Makefile` untuk menjalankan perintah-perintah berikut dari terminal Anda:

| Perintah               | Deskripsi                                                                         |
| ---------------------- | --------------------------------------------------------------------------------- |
| `make help`            | Menampilkan semua perintah yang tersedia.                                         |
| `make clean`           | Membersihkan file cache dan file sementara.                                       |
| `make build`           | Membangun executable aplikasi.                                                    |
| `make run`             | Menjalankan aplikasi dalam mode produksi.                                         |
| `make run-dev`         | Menjalankan aplikasi dalam mode pengembangan.                                     |
| `make debug`           | Menjalankan aplikasi dengan debugger Delve.                                       |
| `make migrate-create`  | Membuat file migrasi baru. Contoh: `make migrate-create name=create_users_table`. |
| `make migrate-up`      | Menjalankan semua migrasi yang belum dieksekusi.                                  |
| `make migrate-down`    | Mengembalikan (rollback) migrasi satu langkah.                                    |
| `make seed-create`     | Membuat file seeder baru. Contoh: `make seed-create name=product`.                |
| `make db-seed`         | Menjalankan seeder database untuk menyuntikkan data awal (termasuk admin).        |

---

### Endpoint API

#### Autentikasi

* **`POST /api/auth/register`**: Mendaftarkan pengguna baru secara langsung dan mengembalikan data pengguna dengan status 201 Created.
* **`POST /api/auth/login`**: Login untuk mendapatkan token Paseto.
* **`POST /api/auth/logout`**: Logout pengguna (memerlukan token).

#### Pengguna (Memerlukan Autentikasi)

* **`GET /api/profile`**: Mendapatkan detail profil pengguna yang sedang login.

---

### Struktur Proyek

```text
.
├── domain/
│   ├── dto/                  # Data Transfer Objects (Request & Response)
│   └── user.domain.go        # Definisi model data utama (entitas)
├── internal/
│   ├── api/
│   │   ├── handler/          # Pengelola permintaan HTTP
│   │   ├── middleware/       # Middleware untuk penanganan permintaan
│   │   └── router/           # Pengaturan rute API
│   ├── config/               # Pengelola konfigurasi dari .env
│   ├── constant/             # Menyimpan nilai konstan aplikasi
│   ├── database/
│   │   ├── connection/       # Koneksi ke database
│   │   ├── migration/
│   │   │   ├── migrations/   # Berkas-berkas migrasi
│   │   │   └── migration.go  # Logika migrasi
│   │   └── seeder/           # Penyuntikan data awal
│   │   │   ├── seeders/      # Berkas-berkas seeder
│   │   │   └── seeder.go     # Logika seeder
│   ├── repository/           # Logika interaksi langsung dengan database
│   ├── service/              # Logika bisnis utama aplikasi
│   └── util/                 # Utilitas
├── .env.example              # Contoh variabel lingkungan
├── .gitignore
├── go.mod
├── go.sum
├── main.go                   # Titik masuk utama aplikasi
├── Makefile                  # Skrip untuk otomatisasi tugas
└── README.md

```

## LICENSE

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

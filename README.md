# EDSA - Go Fiber RESTful API

Proyek ini adalah contoh aplikasi API RESTful yang dibangun dengan Go, menggunakan framework **Fiber** dan **GORM** untuk interaksi dengan database **MariaDB**. Proyek ini mengikuti struktur yang mudah di-maintain dan siap untuk skala produksi.

### Fitur Utama

* **Go Fiber**: Framework web yang cepat dan minimalis.
* **Struktur Terlapis**: Memisahkan kode ke dalam lapisan `handler`, `service`, dan `repository` untuk menjaga pemisahan tanggung jawab.
* **GORM**: ORM yang tangguh untuk operasi database.
* **Migrasi Database**: Menggunakan `golang-migrate` untuk manajemen skema database yang aman dan terstruktur.
* **Seeder Database**: Menyediakan data awal untuk kebutuhan pengembangan dan pengujian.
* **Konfigurasi Berbasis `.env`**: Menggunakan variabel lingkungan untuk konfigurasi yang fleksibel.
* **Makefile**: Mengotomatiskan tugas-tugas pengembangan umum seperti menjalankan aplikasi, migrasi, dan seeder.

---

### Persyaratan

* [Go](https://go.dev/dl/) (versi 1.18+)
* [MariaDB](https://mariadb.org/download/)
* `migrate` CLI (diinstal dengan `go install`)

---

### Instalasi

1.  **Clone repositori:**
    ```bash
    git clone [https://github.com/ipincamp/edsa.git](https://github.com/ipincamp/edsa.git)
    cd edsa
    ```

2.  **Unduh dependensi Go:**
    ```bash
    go mod tidy
    ```

---

### Konfigurasi

Buat file `.env` di root proyek dan tambahkan konfigurasi berikut. Pastikan untuk mengganti nilai `DB_USER` dan `DB_PASS` dengan kredensial database MariaDB Anda.

```dotenv
# Application Port
PORT=8080

# Application Environment
APP_ENV=development

# Database Connection
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=edsa_db
DB_LOCATION=Asia/Jakarta
```

---

### Panduan Penggunaan

Gunakan `Makefile` untuk menjalankan perintah-perintah berikut dari terminal Anda:

| Perintah                                      | Deskripsi                                                                        |
| --------------------------------------------- | -------------------------------------------------------------------------------- |
| `make run`                                    | Menjalankan aplikasi dalam mode pengembangan.                                    |
| `make build`                                  | Membangun executable aplikasi.                                                   |
| `make run-prod`                               | Menjalankan aplikasi dalam mode produksi.                                        |
| `make new-migration NAME=nama_migrasi_anda`   | Membuat file migrasi baru. Contoh: `make new-migration NAME=create_users_table`. |
| `make migrate-up`                             | Menjalankan semua migrasi yang belum dieksekusi.                                 |
| `make migrate-down`                           | Mengembalikan (rollback) migrasi satu langkah.                                   |
| `make seed`                                   | Menjalankan seeder database untuk menyuntikkan data awal.                        |
| `make clean`                                  | Membersihkan file cache dan file sementara.                                      |

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
│   │   ├── handlers/     # Berisi fungsi-fungsi handler API
│   │   └── routes/       # Mendefinisikan endpoint API
│   ├── config/           # Mengatur konfigurasi aplikasi (.env)
│   ├── database/         # Koneksi database dan seeder
│   ├── models/           # Definisi model data
│   ├── repositories/     # Logika interaksi dengan database
│   └── services/         # Logika bisnis utama
├── migrations/           # File-file SQL untuk migrasi
├── scripts/              # Skrip pembantu (misal: untuk migrasi)
├── .env                  # Variabel lingkungan
├── .gitignore
├── go.mod
├── go.sum
└── Makefile              # Skrip untuk otomatisasi

```

## LICENSE

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

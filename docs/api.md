# Dokumentasi API Go EDSA

Dokumentasi ini menjelaskan cara menggunakan dan menguji endpoint API yang tersedia di proyek Go EDSA.

**URL Dasar:** `http://localhost:9000`

## Otentikasi

Beberapa endpoint memerlukan otentikasi menggunakan token Paseto. Untuk mengakses endpoint ini, Anda harus menyertakan header `Authorization` dengan nilai `Bearer <token>`.

## Response Format

Semua response menggunakan format JSON yang konsisten:

### Response Sukses
```json
{
  "status": true,
  "message": "Success message",
  "data": {
    // data response
  }
}
```

### Response Error
```json
{
  "status": false,
  "message": "Error message",
  "error": {
    // error details (opsional)
  }
}
```

## Endpoint

## 🔐 Autentikasi

### 1. Registrasi Pengguna Baru

Endpoint ini digunakan untuk membuat pengguna baru.

- **URL:** `/api/v1/auth/register`
- **Metode:** `POST`
- **Header:**
  - `Content-Type`: `application/json`
- **Body (raw JSON):**

```json
{
  "name": "Nama Lengkap",
  "email": "email@example.com",
  "password": "passwordminimal8karakter"
}
```

- **Validasi:**
  - `name`: Required, minimal 3 karakter
  - `email`: Required, format email valid
  - `password`: Required, minimal 8 karakter

- **Contoh menggunakan cURL:**

```bash
curl -X POST http://localhost:9000/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}'
```

- **Respon Sukses (201 Created):**

```json
{
  "status": true,
  "message": "user registered successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "guest",
    "joined_at": "2025-09-07 10:30:00",
    "updated_at": "2025-09-07 10:30:00"
  }
}
```

- **Respon Error (409 Conflict - Email sudah ada):**

```json
{
  "status": false,
  "message": "email already exists",
  "error": null
}
```

- **Respon Error (400 Bad Request - Validasi gagal):**

```json
{
  "status": false,
  "message": "validation failed",
  "error": [
    {
      "field": "email",
      "message": "Field 'email' must be a valid email address"
    },
    {
      "field": "password",
      "message": "Field 'password' must be at least 8 characters long"
    }
  ]
}
```

### 2. Login Pengguna

Endpoint ini digunakan untuk login dan mendapatkan access token & refresh token.

- **URL:** `/api/v1/auth/login`
- **Metode:** `POST`
- **Header:**
  - `Content-Type`: `application/json`
- **Body (raw JSON):**

```json
{
  "email": "email@example.com",
  "password": "passwordanda"
}
```

- **Validasi:**
  - `email`: Required, format email valid
  - `password`: Required

- **Contoh menggunakan cURL:**

```bash
curl -X POST http://localhost:9000/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john.doe@example.com",
  "password": "password123"
}'
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "login successful",
  "data": {
    "access_token": "v2.local.xxxxxxxxxxxx",
    "refresh_token": "v2.local.yyyyyyyyyyyy"
  }
}
```

- **Respon Error (401 Unauthorized):**

```json
{
  "status": false,
  "message": "invalid credentials",
  "error": null
}
```

### 3. Refresh Token

Endpoint ini digunakan untuk mendapatkan access token baru menggunakan refresh token.

- **URL:** `/api/v1/auth/refresh`
- **Metode:** `POST`
- **Header:**
  - `Content-Type`: `application/json`
- **Body (raw JSON):**

```json
{
  "refresh_token": "v2.local.yyyyyyyyyyyy"
}
```

- **Contoh menggunakan cURL:**

```bash
curl -X POST http://localhost:9000/api/v1/auth/refresh \
-H "Content-Type: application/json" \
-d '{
  "refresh_token": "v2.local.yyyyyyyyyyyy"
}'
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "token refreshed successfully",
  "data": {
    "access_token": "v2.local.newxxxxxxxxx"
  }
}
```

### 4. Logout

Endpoint ini digunakan untuk logout pengguna.

- **URL:** `/api/v1/auth/logout`
- **Metode:** `POST`
- **Header:**
  - `Authorization`: `Bearer <access_token>`

- **Contoh menggunakan cURL:**

```bash
curl -X POST http://localhost:9000/api/v1/auth/logout \
-H "Authorization: Bearer <access_token>"
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "logged out successfully",
  "data": null
}
```

## 👤 User Management

### 5. Get Profile (Self)

Endpoint ini digunakan untuk mendapatkan profil pengguna yang sedang login.

- **URL:** `/api/v1/users/profile`
- **Metode:** `GET`
- **Header:**
  - `Authorization`: `Bearer <access_token>`
- **Access Level:** Authenticated users

- **Contoh menggunakan cURL:**

```bash
curl -X GET http://localhost:9000/api/v1/users/profile \
-H "Authorization: Bearer <access_token>"
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "profile retrieved successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "guest",
    "joined_at": "2025-09-07 10:30:00",
    "updated_at": "2025-09-07 10:30:00"
  }
}
```

### 6. Update Profile (Self)

Endpoint ini digunakan untuk mengupdate profil pengguna yang sedang login.

- **URL:** `/api/v1/users/profile`
- **Metode:** `PUT`
- **Header:**
  - `Authorization`: `Bearer <access_token>`
  - `Content-Type`: `application/json`
- **Access Level:** Authenticated users

- **Body (raw JSON):**

```json
{
  "name": "New Name",
  "old_password": "current_password",
  "new_password": "new_password123",
  "new_password_confirmation": "new_password123"
}
```

- **Validasi:**
  - `name`: Opsional, minimal 3 karakter jika disertakan
  - `old_password`: Required jika ingin ganti password
  - `new_password`: Opsional, minimal 8 karakter
  - `new_password_confirmation`: Required jika new_password disertakan, harus sama dengan new_password

- **Contoh menggunakan cURL:**

```bash
curl -X PUT http://localhost:9000/api/v1/users/profile \
-H "Authorization: Bearer <access_token>" \
-H "Content-Type: application/json" \
-d '{
  "name": "John Smith",
  "old_password": "password123",
  "new_password": "newpassword456",
  "new_password_confirmation": "newpassword456"
}'
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "profile updated successfully",
  "data": null
}
```

- **Respon Error (400 Bad Request - Password salah):**

```json
{
  "status": false,
  "message": "old password is incorrect",
  "error": null
}
```

### 7. Get All Users (Admin Only)

Endpoint ini digunakan untuk mendapatkan daftar semua pengguna dengan pagination.

- **URL:** `/api/v1/users`
- **Metode:** `GET`
- **Header:**
  - `Authorization`: `Bearer <access_token>`
- **Access Level:** Admin only
- **Query Parameters:**
  - `page` (opsional): Halaman (default: 1)
  - `limit` (opsional): Jumlah data per halaman (default: 10, max: 100)
  - `role` (opsional): Filter berdasarkan role (`admin`, `teacher`, `student`, `guest`)

- **Contoh menggunakan cURL:**

```bash
curl -X GET "http://localhost:9000/api/v1/users?page=1&limit=10&role=student" \
-H "Authorization: Bearer <access_token>"
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "users retrieved successfully",
  "data": {
    "data": [
      {
        "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "role": "student",
        "joined_at": "2025-09-07 10:30:00",
        "updated_at": "2025-09-07 10:30:00"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

### 8. Get User by ID

Endpoint ini digunakan untuk mendapatkan data pengguna berdasarkan ID.

- **URL:** `/api/v1/users/{userId}`
- **Metode:** `GET`
- **Header:**
  - `Authorization`: `Bearer <access_token>`
- **Access Level:** Self atau Admin
- **Parameters:**
  - `userId`: UUID pengguna

- **Contoh menggunakan cURL:**

```bash
curl -X GET http://localhost:9000/api/v1/users/a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6 \
-H "Authorization: Bearer <access_token>"
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "user retrieved successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "student",
    "joined_at": "2025-09-07 10:30:00",
    "updated_at": "2025-09-07 10:30:00"
  }
}
```

### 9. Update User by ID

Endpoint ini digunakan untuk mengupdate data pengguna berdasarkan ID.

- **URL:** `/api/v1/users/{userId}`
- **Metode:** `PUT`
- **Header:**
  - `Authorization`: `Bearer <access_token>`
  - `Content-Type`: `application/json`
- **Access Level:** Self atau Admin
- **Parameters:**
  - `userId`: UUID pengguna

- **Body (raw JSON):**

```json
{
  "name": "Updated Name",
  "old_password": "current_password",
  "new_password": "new_password123",
  "new_password_confirmation": "new_password123"
}
```

- **Contoh menggunakan cURL:**

```bash
curl -X PUT http://localhost:9000/api/v1/users/a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6 \
-H "Authorization: Bearer <access_token>" \
-H "Content-Type: application/json" \
-d '{
  "name": "Updated John Doe"
}'
```

- **Respon Sukses (200 OK):**

```json
{
  "status": true,
  "message": "user updated successfully",
  "data": null
}
```

## 🔒 Access Control & Permissions

API menggunakan sistem role-based dan permission-based access control:

### Roles
- **admin**: Akses penuh ke semua endpoint
- **teacher**: Akses terbatas untuk fitur pengajar
- **student**: Akses terbatas untuk fitur siswa
- **guest**: Akses minimal (default untuk user baru)

### Access Levels
- **Public**: Dapat diakses tanpa autentikasi
- **Authenticated**: Memerlukan access token yang valid
- **Self or Admin**: Hanya bisa diakses oleh pemilik resource atau admin
- **Admin Only**: Hanya bisa diakses oleh admin

## ❌ Error Responses

### Common Error Codes

#### 400 Bad Request
```json
{
  "status": false,
  "message": "validation failed",
  "error": [
    {
      "field": "email",
      "message": "Field 'email' must be a valid email address"
    }
  ]
}
```

#### 401 Unauthorized
```json
{
  "status": false,
  "message": "authorization header is not provided",
  "error": null
}
```

#### 403 Forbidden
```json
{
  "status": false,
  "message": "insufficient role permissions",
  "error": null
}
```

#### 404 Not Found
```json
{
  "status": false,
  "message": "user not found",
  "error": null
}
```

#### 409 Conflict
```json
{
  "status": false,
  "message": "email already exists",
  "error": null
}
```

#### 500 Internal Server Error
```json
{
  "status": false,
  "message": "an unexpected error occurred",
  "error": null
}
```

## 🚀 Testing

### Menggunakan cURL
Semua contoh cURL di atas dapat digunakan langsung untuk testing.

### Menggunakan Postman/Insomnia
1. Import collection atau buat request manual
2. Set base URL: `http://localhost:9000`
3. Untuk endpoint yang memerlukan auth, tambahkan header:
   - Key: `Authorization`
   - Value: `Bearer <access_token>`

### Flow Testing yang Disarankan
1. **Register** user baru
2. **Login** dengan user tersebut
3. **Get Profile** untuk memastikan token bekerja
4. **Update Profile** untuk testing update
5. Testing endpoint admin (jika memiliki user admin)

---

**Catatan:** Ganti `<access_token>` dengan token yang Anda dapatkan dari endpoint login. Token memiliki masa berlaku terbatas sesuai konfigurasi di `.env`.

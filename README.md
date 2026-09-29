# Campus Lost & Found API

Backend-only REST API untuk pengelolaan barang hilang dan barang temuan di lingkungan kampus.
Project ini dibuat sebagai project integrasi materi Praktikum Pemrograman Backend Lanjut Modul 1–7 menggunakan Go + Fiber + PostgreSQL.

## Materi yang diintegrasikan
- Modul 1: sintaks Go, struct, method, pointer, slice/map, goroutine/channel dasar.
- Modul 2: REST API, HTTP method/status, PUT vs PATCH, query string, request/response header.
- Modul 3: PostgreSQL, migration, pgx connection pool, repository pattern, parameterized query.
- Modul 4: Clean Architecture dan pemisahan layer.
- Modul 5: bcrypt, JWT access token, refresh token, logout, authentication middleware.
- Modul 6: RBAC, permissions, fail closed, ownership, 401/403.
- Modul 7: validator deklaratif, PATCH pointer + omitnil, cursor pagination, CSV content negotiation, centralized ErrorHandler, error code stabil, request_id.

## Menjalankan
1. Buat database PostgreSQL: `campus_lost_found`.
2. Salin `.env.example` menjadi `.env`, lalu sesuaikan konfigurasi.
3. Jalankan migration SQL secara berurutan.
4. Jalankan `go mod tidy`.
5. Jalankan `go run .`.

## Endpoint utama
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`
- `GET /api/v1/items`
- `GET /api/v1/items/:id`
- `POST /api/v1/items`
- `PUT /api/v1/items/:id`
- `PATCH /api/v1/items/:id`
- `DELETE /api/v1/items/:id`
- `POST /api/v1/items/:id/claims`
- `GET /api/v1/items/:id/claims`
- `PATCH /api/v1/claims/:id/status`

## Content negotiation
Default response adalah JSON. Untuk daftar item, gunakan `Accept: text/csv` untuk mendapatkan CSV.

## Cursor pagination
Contoh: `GET /api/v1/items?limit=5`. Jika masih ada data berikutnya, response memiliki `meta.next_cursor` dan `meta.has_more=true`.

## Testing
`go test ./...`

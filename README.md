# Go Shorty – High Performance URL Shortener Service

![version](https://img.shields.io/badge/version-1.0-blue)
![status](https://img.shields.io/badge/status-stable-brightgreen)
![go](https://img.shields.io/badge/Go-1.24-%2300ADD8?logo=go)
![postgres](https://img.shields.io/badge/PostgreSQL-16-%234169E1?logo=postgresql)
![architecture](https://img.shields.io/badge/Architecture-Clean%20Layered-orange)
![license](https://img.shields.io/badge/License-MIT-green)

---

Introduction – GoShorty v1.0

Live Demo: https://go-shorty-production.up.railway.app/

---

## Table of Contents
1. [Problem Description](#1-mo-ta-bai-toan-problem-description)
2. [Architecture & Core System Design](#2-architecture--core-system-design-goshorty-v10)
3. [Challenges & Solutions](#3-thach-thuc--giai-phap-challenges--solutions)
4. [API Flow](#4-api-flow)
5. [Database Schema](#5-database-schema)
6. [Implementation Overview](#6-implementation-overview)
7. [Limitations & Future Improvements](#7-limitations--future-improvements)
8. [Folder Structure](#8-folder-structure)
9. [Deployment Guide / How to Run](#9-deployment-guide--how-to-run)
10. [Project Status](#10-project-status)

---

# 1. Mô tả bài toán (Problem Description)

Các ứng dụng hiện đại cần chuyển đổi URL dài, phức tạp thành định danh ngắn gọn. Các thách thức kỹ thuật bao gồm:

- **Tính duy nhất của short code:**  
  Phải đảm bảo không trùng lặp, tránh check-and-retry tốn tài nguyên.

- **Click counter an toàn:**  
  Phải đếm click chính xác khi nhiều người dùng truy cập cùng lúc.

- **Tối ưu hiệu năng đọc:**  
  Luồng redirect cần nhanh, tránh bottleneck DB.

> ⚠ GoShorty v1.0: single instance, atomic updates row-level, traffic nhỏ/medium, chưa High Availability.

---


# 2. Architecture & Core System Design (GoShorty v1.0)

GoShorty v1.0 được thiết kế theo **Clean Layered Architecture**, tập trung vào tính rõ ràng, dễ bảo trì, và mở rộng sau này.

<img width="1021" height="286" alt="Screenshot 2025-12-07 at 21 38 47" src="https://github.com/user-attachments/assets/bf21a78f-056b-425d-a43e-7a8a10146ef0" />

## Nguyên tắc thiết kế

* **Đơn giản (Simplicity):** Mô hình hiện tại là **API Server → PostgreSQL**, không dùng microservices, giúp triển khai và vận hành dễ dàng.
* **Hiệu năng hợp lý (Performance-oriented):** Sử dụng **REST API** để tối ưu cache và giảm overhead cho luồng redirect.
* **Mở rộng trong tương lai (Extensible):** API server **stateless**, thuận lợi cho việc scale ngang sau này nếu thêm Docker/K8s/Redis.

> ⚠ Lưu ý: V1.0 là **single instance**, chưa thực sự High Availability hay horizontal scaling.

## Tư duy phân tách luồng xử lý (Read/Write Path – Conceptual)

Hệ thống **chưa tách thành 2 service riêng biệt**, nhưng code được viết theo tư duy **Write Path / Read Path**, giúp dễ mở rộng sau này:

* **Write Path (Tạo link):**

  * Ưu tiên **tính đúng đắn**: validate URL, tránh lưu trữ rác.
  * Base62 encoding từ BIGSERIAL ID → short code không trùng.
  * Lưu vào DB và trả metadata.

* **Read Path (Redirect):**

  * Ưu tiên **tốc độ**: Nhận request → Query DB → HTTP 302 redirect.
  * Ghi nhận lượt click bằng **một phép cộng trong SQL**.

> Đây là **tư duy thiết kế**, chưa tách service thực sự.

## Tính năng cốt lõi (Core Capabilities)

* **Short URL Creation:** Base62 từ BIGSERIAL ID → mã ngắn, dễ đọc, không trùng.
* **Redirect Engine:** HTTP 302 + Location Header → chuyển hướng chuẩn, hỗ trợ SEO và analytics.
* **Click Counter / Analytics:**

  * Code hiện tại dùng `UPDATE links SET clicks_count = clicks_count + 1 WHERE short_code = ?`.
  * Trên **single PostgreSQL instance**, phép cộng này là **atomic trên từng row**, đảm bảo **không mất click khi nhiều request đồng thời**.
  * ⚠ Tuy nhiên, đây **không phải giải pháp phân tán**, và chưa tối ưu cho hàng triệu request/giây hay multi-instance.
* **Data Integrity:** PostgreSQL constraints (Unique Index) → ngăn dữ liệu rác, đảm bảo tính nhất quán.

---


# 3. Thách thức & Giải pháp (Challenges & Solutions)

## A. Vấn đề "Duplicate Key" & Quy trình Lưu trữ

* **Vấn đề:**
  Trong luồng "2-Step Save" (Insert URL → Lấy ID → Update mã ngắn), nếu dùng hàm `Save()` chung chung của GORM, có thể gặp lỗi **Duplicate Key**, do GORM cố gắng insert lại bản ghi đã có ID thay vì chỉ update trường cần thiết.
* **Giải pháp:**
  Tách rõ ràng tầng Repository, sử dụng thao tác **explicit operations**:

  * `Create()` → insert bản ghi mới, lấy ID tự tăng.
  * `UpdateShortCode()` → chỉ update trường `short_code` bằng `UpdateColumn`.
* **Kết quả:**
  Quy trình lưu trữ trở nên **an toàn, nguyên tử**, tránh xung đột duplicate key.

---

## B. Race Condition trong Thống kê lượt click (Analytics)

* **Vấn đề:**
  Logic "đọc → cộng → lưu" truyền thống sẽ gây sai lệch số liệu click nếu nhiều người dùng truy cập cùng lúc.
* **Giải pháp:**
  Sử dụng **atomic update trên từng row** của PostgreSQL:

  ```sql
  UPDATE links SET clicks_count = clicks_count + 1 WHERE short_code = ?
  ```

  PostgreSQL đảm bảo **row-level atomicity**, tránh mất dữ liệu khi nhiều request đồng thời.
* **Kết quả:**

  * Click counter **chính xác trên single instance**.
  * ⚠ Chưa phải giải pháp phân tán, và chưa tối ưu cho traffic cực lớn hoặc multi-instance.

---


# 4. API Flow

GoShorty v1.0 có hai luồng chính: **Create Short URL** và **Redirect**. Cả hai được thiết kế đơn giản, hiệu quả, phù hợp với traffic nhỏ/medium trên single instance PostgreSQL.

---

### A. Create Short URL Flow
<img width="1163" height="387" alt="Screenshot 2025-12-07 at 21 42 10" src="https://github.com/user-attachments/assets/a9d003c3-9a8b-4b79-a16c-8dc30fb257fc" />

**Mục tiêu:** Tạo một short code duy nhất cho URL gốc, lưu metadata và trả về thông tin cho client.

**Flow chi tiết:**

1. **Client gửi POST request** với JSON body:

   ```json
   { "original_url": "https://example.com/long-url" }
   ```

2. **Handler Layer (Gin)** nhận request, validate dữ liệu cơ bản: URL hợp lệ, không rỗng.

3. **Service Layer**:

   * Kiểm tra URL đã tồn tại chưa (optional, nếu muốn tránh duplicate URLs).
   * Tạo bản ghi mới trong DB bằng `Repository.Create()`, lấy `id` tự tăng.
   * Chuyển `id` sang short code bằng **Base62 encoding**.
   * Cập nhật short code trong cùng bản ghi (`Repository.UpdateShortCode()`).

4. **Repository Layer (GORM)** thực hiện thao tác Insert và Update:

   * `Create()` → insert bản ghi với `original_url` và `id`.
   * `UpdateShortCode()` → update `short_code` mà **không ghi đè toàn bộ row**, tránh duplicate key và race condition.

5. **Handler trả về JSON**:

   ```json
   { "short_code": "aB3dE1", "original_url": "...", "clicks_count": 0 }
   ```

**Lý do thiết kế này:**

* Tách rõ ràng **Insert vs Update**, tránh lỗi Duplicate Key khi dùng `Save()` của GORM.
* Base62 từ ID **đảm bảo collision-free** mà không cần check-and-retry.
* Thiết kế này đơn giản, dễ mở rộng, vẫn giữ **atomicity** trên single instance.

---

### B. Redirect Flow

<img width="865" height="484" alt="Screenshot 2025-12-07 at 21 40 23" src="https://github.com/user-attachments/assets/976068f9-15ff-451d-b4bb-b987e7745a1b" />

**Mục tiêu:** Chuyển hướng người dùng từ short URL → URL gốc và ghi nhận click.

**Flow chi tiết:**

1. **Client truy cập short URL**, ví dụ:

   ```
   GET /aB3dE1
   ```

2. **Handler Layer (Gin)** nhận request, extract `short_code` từ path.

3. **Service Layer**:

   * Gọi Repository tìm bản ghi theo `short_code`.
   * Nếu tồn tại, trả về `original_url`.
   * Nếu không, trả lỗi 404.

4. **Repository Layer (GORM / PostgreSQL)**:

   * `SELECT * FROM links WHERE short_code = ?`
   * `UPDATE links SET clicks_count = clicks_count + 1 WHERE short_code = ?`
     → phép cộng **atomic trên row-level**, PostgreSQL đảm bảo không mất click dù nhiều request đồng thời.

5. **Handler Layer** gửi **HTTP 302 redirect**:

   ```http
   HTTP/1.1 302 Found
   Location: https://example.com/long-url
   ```

**Lý do thiết kế này:**

* Tách luồng **Read Path / Write Path**: Redirect ưu tiên tốc độ, ghi click là một phép cộng đơn giản trong DB.
* Sử dụng **atomic update row-level** trên PostgreSQL để đảm bảo số liệu thống kê chính xác.
* Giải pháp này **nhanh, đơn giản, không cần locking** ở tầng ứng dụng.
* Nhược điểm: chưa hỗ trợ multi-instance hoặc traffic cực lớn → cần caching/Redis cho hot path trong tương lai.


### Tổng quan

* **Create URL** → ưu tiên **tính đúng đắn, tránh duplicate, atomic insert/update**.
* **Redirect** → ưu tiên **tốc độ, atomic click count**, luồng đơn giản, O(1) lookup trên Unique Index.
* Kiến trúc này phù hợp với **v1.0: single instance, traffic nhỏ/medium**, dễ nâng cấp sau này với caching và horizontal scaling.

---

## 5. Database Schema

Thiết kế schema tập trung vào hiệu năng đọc (redirect) và tính toàn vẹn dữ liệu ở mức single instance.
V1.0 dùng bảng đơn, tránh JOIN để tối ưu cho luồng redirect.

<img width="300" alt="Schema Diagram" src="https://github.com/user-attachments/assets/948a8d70-394e-4877-9997-d6a9c3b27bd8" />

###  Cấu trúc bảng `links`

| Column         | Type                     | Mô tả                                              |
| -------------- | ------------------------ | -------------------------------------------------- |
| `id`           | BIGSERIAL (PK)           | ID tự tăng, ánh xạ từ `uint primaryKey` trong GORM |
| `short_code`   | VARCHAR(10) UNIQUE       | Mã rút gọn, tra cứu nhanh nhờ Unique Index         |
| `original_url` | TEXT NOT NULL            | URL gốc, không cho phép null                       |
| `clicks_count` | BIGINT                   | Tổng lượt click, phép cộng atomic trên single row  |
| `created_at`   | TIMESTAMP WITH TIME ZONE | Thời điểm tạo bản ghi, tự động sinh bởi GORM       |


---

### Chi tiết & Lý do kỹ thuật

1. **ID (BIGSERIAL)**

   * Dùng làm cơ sở cho Base62 encoding → short code.
   * 64-bit đủ lớn, hỗ trợ hàng tỷ bản ghi, loại bỏ rủi ro overflow.

2. **short_code (VARCHAR(10), UNIQUE Index)**

   * Đảm bảo **collision-free** trên single instance.
   * Tra cứu gần như O(1), tối ưu cho luồng redirect.

3. **original_url (TEXT, NOT NULL)**

   * Lưu URL gốc dài và phức tạp.
   * **Không enforce unique**, có thể lưu trùng URL nhưng short code khác nhau.

4. **clicks_count (BIGINT)**

   * Đếm lượt truy cập.
   * **Atomic row-level update** (`UPDATE ... SET clicks_count = clicks_count + 1`) đảm bảo số liệu không mất khi nhiều request đồng thời **trên single DB instance**.

5. **created_at (TIMESTAMP)**

   * Lưu thời điểm tạo, phục vụ audit hoặc phân tích.

---

### Data Integrity / ACID

* **Atomicity:** chỉ đảm bảo cho mỗi statement riêng lẻ. 2-step create URL chưa wrap transaction → có thể không atomic tuyệt đối.
* **Consistency:** enforced nhờ **Unique Index** trên `short_code` và **NOT NULL** trên `original_url`.
* **Isolation:** row-level atomic update đảm bảo không mất click khi nhiều request đồng thời trên cùng row.
* **Durability:** PostgreSQL đảm bảo dữ liệu đã commit được ghi vào disk, không mất khi crash.

> ⚠ Lưu ý: V1.0 là **single instance**, chưa có multi-node replication hay caching layer, nên ACID chỉ đảm bảo trong phạm vi instance đơn lẻ.

---

## 6. Implementation Overview

GoShorty v1.0 tuân thủ **Clean Layered Architecture**, tách biệt rõ ràng trách nhiệm giữa các layer:

### **1. Handler Layer (Gin)**

* Nhận request từ client (JSON binding, query param).
* Validate dữ liệu cơ bản (URL hợp lệ, không rỗng).
* Gọi Service layer để thực hiện business logic.
* Trả response chuẩn (JSON + HTTP status).
* Không chứa business logic, giữ **stateless**.

### **2. Service Layer**

* Chứa **tất cả business logic**:

  * Kiểm tra URL tồn tại (optional).
  * Base62 encoding từ ID → short code.
  * Orchestrate luồng Write Path / Read Path.
  * Error handling: mapping từ Repository → Service → HTTP status.
* Thực hiện **transaction nhỏ / atomic** khi cần.
* Tư duy **Write Path / Read Path**:

  * Write Path → tạo URL, insert + update short code.
  * Read Path → redirect, atomic update click count.

### **3. Repository Layer (GORM / PostgreSQL)**

* Truy xuất DB: `Create()`, `UpdateShortCode()`, `GetByShortCode()`.
* Abstraction để Service layer **không phụ thuộc vào GORM**.
* Thực hiện **row-level atomic operations**: đảm bảo click counter chính xác.
* Sử dụng Unique Index, constraints để đảm bảo **data integrity**.

### **4. Error Handling & Atomicity**

* Sử dụng Go `errors.Is()` để so sánh lỗi across layers.
* Row-level operations trên PostgreSQL đảm bảo **atomic update clicks_count**.
* Lưu ý: 2-step create URL chưa wrap transaction → có thể không hoàn toàn atomic, cần cải thiện ở v2.0.

### **5. Tóm tắt luồng code**

| Layer      | Chức năng chính                                                                |
| ---------- | ------------------------------------------------------------------------------ |
| Handler    | HTTP I/O, JSON binding, response, validate                                     |
| Service    | Business logic, Base62 encoding, orchestrate Read/Write, error mapping         |
| Repository | DB access, atomic row-level update, enforce constraints, abstract GORM details |

> **Lưu ý:** Implementation này tối ưu cho **v1.0: single instance, traffic nhỏ/medium**. Multi-instance, caching layer và Redis sẽ được thêm vào các phiên bản sau.

---

## 7. Limitations & Future Improvements

### 7.1 Limitations (Nhược điểm v1.0)

* **Single instance:**
  ACID chỉ đảm bảo trong phạm vi một PostgreSQL instance duy nhất; chưa hỗ trợ multi-node replication → chưa High Availability.

* **Traffic nhỏ/medium:**
  Atomic row-level update cho click counter hoạt động tốt, nhưng chưa tối ưu cho hàng triệu request/giây hoặc multi-instance.

* **No caching layer:**
  Luồng redirect vẫn query trực tiếp DB → có thể trở thành bottleneck khi link hot.

* **2-step Create URL:**
  Quy trình Insert + Update short code chưa wrap transaction → có thể không hoàn toàn atomic nếu gặp lỗi giữa các bước.

* **No comprehensive tests:**
  Chưa có đầy đủ unit & integration tests; error handling chỉ cover single instance, chưa đảm bảo cho edge cases.

* **Duplicate URLs:**
  Hiện tại, system cho phép cùng một original URL tạo nhiều short code → có thể dẫn đến duplicate entry nếu muốn tránh trùng.

---

### 7.2 Future Improvements – v1.1

Các cải tiến này tập trung vào việc **khắc phục những hạn chế còn tồn tại của v1.0**, có thể triển khai trong thời gian ngắn (ví dụ 2–3 ngày):

* **Atomic Create URL:**
  Wrap quy trình Insert + Update short code trong transaction để đảm bảo **tính toàn vẹn** tuyệt đối.

* **Duplicate URL Handling:**
  Kiểm tra duplicate original URL trước khi tạo short code để tránh việc cùng một URL có nhiều short code không cần thiết.

* **Basic Caching (Optional):**
  Thêm cache tạm thời cho redirect hot path để giảm query DB, tăng tốc độ cho các link được truy cập nhiều.

* **Improved Error Handling & Logging:**
  Bổ sung logging chi tiết hơn, handle edge cases, tránh silent failure trong các API.

* **Unit & Integration Testing (Minimal):**
  Viết thêm một số test đơn giản cho service & repository để đảm bảo các chức năng cơ bản hoạt động ổn định.

> ⚠ Mục tiêu v1.1: sửa những thiếu sót nhỏ, nâng tính ổn định và trải nghiệm người dùng, **không thêm tính năng mới lớn hay thay đổi kiến trúc**.

---

## 8. Folder Structure

Cấu trúc thư mục GoShorty v1.0 được tổ chức theo **Clean Layered Architecture**, dễ hiểu và dễ bảo trì:

```bash
go-shorty/
├── cmd/                  # Entry point: main.go, start server
├── config/               # Cấu hình ứng dụng, env loader
├── internal/             # Core code theo Clean Architecture
│   ├── handler/          # HTTP Handlers (Gin): nhận request, trả response
│   ├── service/          # Business Logic: validation, encoding, orchestration
│   ├── repository/       # DB access (GORM): truy vấn, insert, update
│   ├── model/            # Entity / struct / DB mapping
│   └── util/             # Helper functions, utils
├── deployments/
│   └── docker/
│       └── init.sql      # Seed / DB initialization, phục vụ Docker/Production
├── .env                  # Environment variables
├── go.mod                # Go module definition
├── go.sum                # Module checksum
└── README.md             # Project documentation
```

**Giải thích nhanh:**

* `cmd/` → nơi bắt đầu server, giữ `main.go`.
* `config/` → load cấu hình từ `.env` hoặc default values.
* `internal/` → code chính, tách theo layers: handler → service → repository → model → util.
* `deployments/docker/` → các file phục vụ deploy, seed DB (`init.sql`) cho Docker container.
* `.env` → cấu hình môi trường (DB URL, PORT…).
* `go.mod` & `go.sum` → quản lý dependency.
* `README.md` → tài liệu tổng quan, hướng dẫn, schema, flow, folder structure…

> ⚠ V1.0 hiện tại sử dụng **Railway deployment**, chưa dùng Docker trực tiếp.
---

## 9. Deployment Guide / How to Run

### 1. Clone repository

```bash
git clone https://github.com/quandang08/go-shorty.git
cd go-shorty
```

### 2. Setup Environment

Tạo file `.env` dựa trên cấu hình của bạn:

```env
DB_HOST=localhost
DB_USER=your_user
DB_PASSWORD=your_pass
DB_NAME=shorty_db
PORT=8080
SHORT_DOMAIN=http://localhost:8080/
```

> ⚠ Lưu ý: Điều chỉnh thông tin DB và port theo môi trường của bạn.

### 3. Run the application

```bash
go run ./cmd/server/main.go
```

Sau đó truy cập:

```
http://localhost:8080
```

---

### 4. Example API

#### A. Create Short URL

```bash
curl -X POST http://localhost:8080/api/v1/links \
     -H 'Content-Type: application/json' \
     -d '{"original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
```

**Response:**

```json
201 Created
{
    "short_code": "1",
    "short_url": "http://localhost:8080/1",
    "original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "clicks_count": 0
}
```

#### B. Redirect

Truy cập trình duyệt:

```
http://localhost:8080/1
```

→ Chuyển hướng đến URL gốc (HTTP 302 Found).

#### C. Get Link Analytics

```bash
curl http://localhost:8080/api/v1/links/1
```

**Response:**

```json
200 OK
{
    "short_code": "1",
    "clicks_count": 5,
    "original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}
```

#### D. Get All Links

```bash
curl http://localhost:8080/api/v1/links
```

**Response:**

```json
[
    {
        "short_code": "1",
        "original_url": "https://devops-is-fun.com/railway-demo-1",
        "clicks_count": 2,
        "created_at": "2025-12-07T02:34:34.872392Z",
        "short_url": "http://localhost:8080/1"
    },
    {
        "short_code": "2",
        "original_url": "https://devops-is-fun.com/railway-demo-1",
        "clicks_count": 10,
        "created_at": "2025-12-07T02:40:23.254892Z",
        "short_url": "http://localhost:8080/2"
    },
    {
        "short_code": "3",
        "original_url": "https://devops-is-fun.com/railway-demo-1",
        "clicks_count": 0,
        "created_at": "2025-12-07T03:03:18.010742Z",
        "short_url": "http://localhost:8080/3"
    }
]
```

> ⚠ Lưu ý: Hai API `Get Link Analytics` và `Get All Links` chỉ phục vụ mục đích xem thông tin, không ảnh hưởng đến luồng chính (Create / Redirect).

---

## 10. Project Status

* **Version:** 1.0
* **Status:** Stable
* **Traffic Capacity:** Single instance PostgreSQL, phù hợp cho traffic nhỏ/medium.
* **Roadmap / Future Features:**

  * Redis caching cho hot path (GET /:code)
  * Docker / Kubernetes deployment, horizontal scaling
  * Observability: Prometheus + Grafana
  * Unit & Integration tests
  * Multi-instance / distributed setup

> ⚠ Lưu ý: v1.0 là phiên bản đầu, tập trung vào **tính đúng đắn, atomic updates, O(1) redirect**, chưa tối ưu cho traffic cực lớn hay High Availability.


---

<br><br><br>

## Ideas / Coming Soon – GoShorty 2.0

> “We’re not just shortening links. We’re redefining how people interact with URLs.” – GoShorty Vision

### Smart Short Code Suggestions – Trải nghiệm mượt mà

Trong thực tế, khi người dùng muốn tự nhập alias cho link của họ, thường gặp một vấn đề: nếu cơ sở dữ liệu có hàng triệu short code, việc nhập thủ công dễ dẫn tới **trùng lặp liên tục**. Hệ thống sẽ buộc người dùng thử đi thử lại, gây **mất thời gian, khó chịu và trải nghiệm tệ**.

GoShorty 2.0 giải quyết vấn đề này một cách thông minh. Hệ thống sẽ **hiểu ý định của người dùng**, gợi ý các lựa chọn khả thi, nhanh chóng và an toàn. Người dùng vẫn chủ động nhập alias, nhưng không còn phải thử đi thử lại nhiều lần nữa.

> Chi tiết thuật toán? Chỉ GoShorty biết.
> Trải nghiệm sẽ nói lên tất cả.

---

### Structured Short Codes – Bảo mật tinh tế

Cho phép người dùng nhập short_code thủ công mở ra rủi ro: SQL Injection, XSS, brute-force… GoShorty 2.0 nói: **“Bạn đặt tên, chúng tôi bảo vệ.”**

Lấy cảm hứng từ JWT (header–payload–signature), short_code được **thiết kế có cấu trúc riêng**, vừa bảo vệ hệ thống vừa vẫn **ngắn gọn và dễ nhớ**.

Cách cấu trúc và thuật toán tối giản như thế nào… đó là một bí mật. Người dùng sẽ cảm nhận hiệu quả bảo vệ ngay từ trải nghiệm, nhưng **chi tiết chỉ GoShorty mới biết**.

---

<br>

Coming Soon – GoShorty 2.0

<br>

**Author:** CodeWithAmu  
**GitHub:** [https://github.com/CodeWithAmu](https://github.com/CodeWithAmu)  
**Made with ❤️ using Golang.**


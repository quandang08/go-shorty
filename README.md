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

<img width="1022" height="408" alt="Screenshot 2025-12-04 at 10 00 55" src="https://github.com/user-attachments/assets/33ebcab9-ddc2-4efc-8bfc-eaa12b5e52d4" />

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

<img width="767" height="282" alt="Screenshot 2025-12-04 at 18 35 23" src="https://github.com/user-attachments/assets/e7c7d60e-3f54-47d9-a363-f501ef6415a2" />

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

<img width="784" height="430" alt="Screenshot 2025-12-04 at 18 36 53" src="https://github.com/user-attachments/assets/c7905d66-9529-4722-9064-7b89d67bee77" />

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


---

## 7. Limitations & Future Improvements



---
Mục bổ sung : 8. Folder Structure, 9. Deployment Guide, 10. Example API, 11. Project Status

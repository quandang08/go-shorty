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

# 1. Vấn đề Cốt lõi (The Core Problem)

Các ứng dụng hiện đại cần một giải pháp tin cậy để chuyển đổi các URL dài và phức tạp thành các định danh ngắn gọn. Mặc dù các dịch vụ bên thứ ba (như Bit.ly) rất phổ biến, chúng đi kèm với rủi ro về phụ thuộc hạ tầng, độ trễ mạng (latency) và thiếu quyền kiểm soát dữ liệu.

GoShorty được xây dựng như một Backend Service tự vận hành (Self-hosted) để giải quyết bài toán rút gọn link với trọng tâm là Độ sẵn sàng cao (High Availability) và Tính toàn vẹn dữ liệu nghiêm ngặt (Strict Data Integrity).

Các Quyết định Kỹ thuật Chính (Key Engineering Decisions)
Thay vì chỉ xây dựng một ứng dụng CRUD đơn giản, chúng tôi tập trung giải quyết các ràng buộc kỹ thuật cụ thể của một hệ thống chịu tải cao:

Đảm bảo Tính Duy nhất (Collision-Free):

Thách thức: Việc sinh chuỗi ngẫu nhiên (như MD5/UUID) đòi hỏi phải truy vấn ngược vào DB để kiểm tra trùng lặp ("check-and-retry"), gây tốn kém tài nguyên và làm chậm hệ thống.

Giải pháp: Sử dụng thuật toán Base62 Encoding dựa trên ID tự tăng (BIGSERIAL) của Database. Giải pháp toán học này đảm bảo tính duy nhất tuyệt đối theo thiết kế, loại bỏ hoàn toàn việc kiểm tra trùng lặp và tối đa hóa hiệu suất ghi.

An toàn Đồng thời trong Thống kê (Concurrency Safety):

Thách thức: Logic đếm click ngây thơ ("đọc-cộng-lưu") sẽ gây ra lỗi "mất dữ liệu" (lost updates) khi có hàng nghìn người dùng click cùng lúc (Race Condition).

Giải pháp: Sử dụng cơ chế Atomic Updates ở cấp độ Database (UPDATE ... SET clicks = clicks + 1). Điều này đảm bảo độ chính xác 100% cho dữ liệu thống kê mà không cần dùng khóa (locking) phức tạp ở tầng ứng dụng Golang.

Đánh đổi: Hiệu năng vs. Độ tin cậy (Performance vs. Reliability):

Quyết định: Chọn PostgreSQL thay vì NoSQL (MongoDB/Redis).

Lý do: Mặc dù NoSQL mở rộng (scale) dễ hơn, nhưng yêu cầu về tuân thủ ACID (để đảm bảo tính nhất quán của dữ liệu quan trọng) được đặt lên hàng đầu. Hệ thống được tối ưu hóa cho tốc độ đọc O(1) trên luồng chuyển hướng (redirect path) bằng cách sử dụng Indexing hiệu quả.

---

# 2. Architecture & Core System Design

Hệ thống tuân thủ Kiến trúc Phân tầng (Clean Layered Architecture) để đảm bảo sự tách biệt trách nhiệm, dễ dàng bảo trì và mở rộng.

<img width="1022" height="408" alt="Screenshot 2025-12-04 at 10 00 55" src="https://github.com/user-attachments/assets/33ebcab9-ddc2-4efc-8bfc-eaa12b5e52d4" />

GoShorty được thiết kế tối giản nhưng đủ mạnh để vận hành dài hạn. Kiến trúc xoay quanh ba ưu tiên quan trọng:

Tốc độ xử lý cực nhanh

Tính ổn định và độ tin cậy cao

Khả năng mở rộng theo chiều ngang mà không phải tái thiết kế

Thay vì phân tán thành nhiều service, GoShorty giữ mô hình API Server → PostgreSQL. Cách tiếp cận này mang lại:

Hệ thống dễ triển khai, dễ vận hành, ít điểm lỗi

Uptime cao vì không phụ thuộc service bên ngoài

Dễ scale bằng cách nhân bản nhiều instance API server

Chi phí thấp nhưng vẫn ổn định trong thời gian dài

GoShorty chọn REST API thay vì serverless hoặc GraphQL vì:

Độ trễ thấp, predictable — rất quan trọng cho redirect

Dễ cache, dễ mở rộng

Logic đơn giản không cần cấu trúc query phức tạp

Tổ chức hai luồng quan trọng nhất: Create & Redirect

GoShorty xoay quanh hai core flows:

Create Short URL – ghi dữ liệu, generate short code

Redirect – đọc dữ liệu, điều hướng nhanh nhất có thể

Tách hai workload này giúp:

Redirect đạt tốc độ tối đa, không bị ảnh hưởng bởi quá trình tạo link

Scale riêng redirect nhiều hơn khi traffic tăng

Tránh xung đột giữa tác vụ đọc và ghi

Giữ codebase rõ ràng: mỗi luồng tối ưu cho đúng nhiệm vụ của nó

Cả hai flow đều được mô tả bằng Sequence Diagram để developer dễ hiểu hệ thống:

Request đi qua Handler → Service → Repository

Làm gì tại mỗi bước

Cách validate và xử lý lỗi

Cách trả response

Điều này đảm bảo ai đọc README cũng đủ hiểu cách hệ thống vận hành.

🚀 Core Features (v1.0)
Short URL Creation

Nhận original_url

Validate input

Lưu vào DB

Encode ID → Base62

Unique index đảm bảo không trùng

Trả về short URL hoàn chỉnh

Redirect Handler

Nhận short_code

Tra cứu trong PostgreSQL

Tăng clicks_count

Trả về HTTP 302 redirect

Luồng cực nhanh, tối thiểu logic

Base62 Encoding

Dựa trên auto-increment ID

Không collision theo thiết kế

Chuỗi ngắn, dễ nhớ

Không cần hash/random phức tạp

Database (PostgreSQL)

1 bảng duy nhất: short_urls

Index:

unique(short_code)

unique(original_url)

Migration tự động bằng GORM

---

# 3. Thách thức & Giải pháp (Challenges & Solutions)

## A. Vấn đề "Duplicate Key" & Quy trình Lưu trữ
  - Vấn đề: Khi triển khai luồng "2-Step Save" (Insert URL $\to$ Lấy ID $\to$ Update Mã ngắn), GORM gặp lỗi Duplicate Key do cố gắng Insert lại bản ghi đã có ID thay vì Update.
  - Giải pháp: Tái cấu trúc tầng Repository để sử dụng các thao tác Tường minh (Explicit Operations). Thay thế hàm Save() chung chung bằng Create() (cho Insert) và UpdateShortCode() (sử dụng UpdateColumn chỉ cập nhật trường cần thiết).
  - Kết quả: Đảm bảo quy trình lưu trữ an toàn, nguyên tử và không có lỗi xung đột.

## B. Race Condition trong Thống kê (Analytics)
  - Vấn đề: Logic "đọc-cộng-lưu" thông thường sẽ làm sai lệch số liệu click khi có nhiều người dùng truy cập cùng lúc.
  - Giải pháp: Sử dụng Atomic Updates tại cấp độ Database (clicks = clicks + 1).
  - Kết quả: Chuyển giao việc khóa (locking) cho PostgreSQL xử lý, đảm bảo chính xác 100% dữ liệu thống kê.
---

# 4. API Flow

### Create Short URL Flow
<img width="767" height="282" alt="Screenshot 2025-12-04 at 18 35 23" src="https://github.com/user-attachments/assets/e7c7d60e-3f54-47d9-a363-f501ef6415a2" />


### Redirect Flow
<img width="784" height="430" alt="Screenshot 2025-12-04 at 18 36 53" src="https://github.com/user-attachments/assets/c7905d66-9529-4722-9064-7b89d67bee77" />


---

## 5. Database Schema
<img width="242" height="308" alt="Screenshot 2025-12-04 at 18 30 30" src="https://github.com/user-attachments/assets/948a8d70-394e-4877-9997-d6a9c3b27bd8" />

Mục đích của thiết kế Schema này không chỉ là lưu trữ dữ liệu, mà là tối ưu hóa cho hai luồng quan trọng nhất: Redirect (đọc) và Create (ghi), đồng thời đảm bảo tính toàn vẹn dữ liệu (Data Integrity) dưới tải cao.

Tận dụng PostgreSQL ACID: Việc chọn PostgreSQL và thiết lập các Unique Index đảm bảo cơ chế ACID (Atomicity, Consistency, Isolation, Durability). Điều này đặc biệt quan trọng để bảo vệ dữ liệu clicks_count khỏi bị sai lệch (Lost Update) dưới tải cao.

Tối ưu hóa Hot Path (Redirect): Thiết kế này sử dụng Single Table để loại bỏ nhu cầu Join bảng, giúp tối giản hóa logic và đạt được tốc độ truy vấn tối đa.

ID: BIGSERIAL (Primary Key): Là cơ sở cho thuật toán Base62. BIGSERIAL đảm bảo có thể lưu trữ hơn 9 triệu triệu link, đủ cho mọi nhu cầu.

short_code: VARCHAR(10) & UNIQUE Index: Giới hạn độ dài tối đa và đặt Unique Index để đảm bảo không bao giờ có hai mã ngắn giống nhau, ngăn chặn xung đột ở tầng DB.

original_url: TEXT & UNIQUE Index: Index này quan trọng để kiểm tra nhanh chóng xem link gốc đã được rút gọn trước đó hay chưa (Duplicate URL Check), tránh lãng phí.

clicks_count: INT: Sử dụng kiểu INT và được bảo vệ bởi Atomic Update trong Tầng Repository.

A. ID (BIGSERIAL - Primary Key)
ID không chỉ là khóa chính mà còn là cơ sở toán học cho thuật toán Base62 Encoding. Việc chọn kiểu BIGSERIAL thay vì SERIAL thông thường đảm bảo hệ thống có khả năng lưu trữ hơn 9 triệu triệu link, loại bỏ hoàn toàn rủi ro tràn số (overflow) trong dài hạn.

B. short_code (VARCHAR(10) - UNIQUE Index)
Đây là cột quan trọng nhất trong luồng Redirect.

Đặt UNIQUE Index trên cột này là bắt buộc để đảm bảo không bao giờ có hai mã ngắn trùng nhau khi lookup, ngăn chặn xung đột ở tầng DB.

Index này cho phép PostgreSQL tìm kiếm và trả về original_url với độ phức tạp O(1) (truy vấn cực nhanh) cho luồng Redirect Hot Path của hệ thống.

C. original_url (TEXT - UNIQUE Index)
UNIQUE Index trên cột này rất quan trọng để thực hiện Duplicate URL Check nhanh chóng ở tầng Service. Mục đích là để kiểm tra xem một link gốc đã được rút gọn trước đó hay chưa, ngăn chặn việc tạo ra các bản ghi trùng lặp và tiết kiệm tài nguyên DB.

Sử dụng kiểu TEXT để chấp nhận độ dài URL linh hoạt và lớn.

D. clicks_count (INT)
Cột này lưu trữ số lần click và là nơi dễ bị lỗi nhất trong tình huống tải cao. Nó được bảo vệ bởi cơ chế Atomic Update ở tầng Repository, đảm bảo tính toàn vẹn và chính xác của dữ liệu dưới mọi điều kiện tải.

---

## 6. Implementation Overview
The implementation strictly follows the Clean Layered Architecture principles:

Handler Layer (Gin): Responsible solely for HTTP I/O (JSON binding, response status, error mapping).

Service Layer: Contains all Business Logic (Validation, Base62 encoding, Existence Check) and orchestrates the transaction flows.

Repository Layer (GORM): Handles database access and abstracts DB operations, ensuring the Service Layer does not depend on GORM specifics.

Error Handling: Utilizes Go's built-in errors.Is() for safe error comparison across layers (e.g., mapping a DB error to a custom Service Error, and finally to an appropriate HTTP status).

---

## 7. Limitations & Future Improvements

For long-term production readiness, the following features are planned for future versions:

Redis Caching Layer: Implement Redis for caching the redirect hot path (GET /:code). This will reduce database latency to near zero for highly-trafficked links, maximizing redirect speed and significantly offloading PostgreSQL.

Containerization (Docker/Kubernetes): Fully implement and test the provided Dockerfile and docker-compose.yml for simplified local environment setup and cloud orchestration (Kubernetes/ECS), ensuring easy horizontal scalability.

Metrics & Observability: Integrate Prometheus and Grafana for monitoring key metrics (Redirect latency, DB queries, Click Volume) to proactively detect failures and capacity issues.

Unit & Integration Testing: Implement comprehensive test suites for the Service and Repository layers to ensure code reliability and prevent regressions during feature expansion.


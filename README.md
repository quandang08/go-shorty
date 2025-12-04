# Go Shorty – High Performance URL Shortener Service

![version](https://img.shields.io/badge/version-1.0-blue)
![status](https://img.shields.io/badge/status-stable-brightgreen)
![go](https://img.shields.io/badge/Go-1.22-%2300ADD8?logo=go)
![postgres](https://img.shields.io/badge/PostgreSQL-16-%234169E1?logo=postgresql)
![architecture](https://img.shields.io/badge/Architecture-Clean%20Layered-orange)
![license](https://img.shields.io/badge/License-MIT-green)

---

## 1. Introduction

GoShorty là hệ thống rút gọn URL được thiết kế theo tiêu chí **đơn giản, rõ ràng và đúng chuẩn backend**.  
Dự án tập trung giải quyết hai vấn đề cốt lõi của bất kỳ URL Shortener hiện đại nào:

- **Tạo short code ngắn, không trùng và hiệu suất cao**
- **Redirect nhanh và chính xác**, đồng thời cập nhật số lượt click an toàn

GoShorty sử dụng kiến trúc và công nghệ tối ưu:

- **Golang** — xử lý concurrency rất tốt (goroutine)
- **PostgreSQL** — đảm bảo tính nhất quán dữ liệu ACID
- **Base62 từ ID tự tăng** — bảo đảm short code duy nhất và truy vấn nhanh

---

## 2. Core Features (Version 1.0)

Phiên bản **v1.0** tập trung xây nền tảng backend ổn định và đúng chuẩn.

### Short URL Creation
- Nhận `original_url`
- Lưu vào DB
- Generate short code bằng Base62 (ID → mã rút gọn)
- Đảm bảo không trùng bằng unique index
- Trả về short URL hoàn chỉnh

### Redirect Handler
- Tra cứu `short_code` trong DB
- Tăng `clicks_count`
- Trả về HTTP 302 redirect tới original URL

### Base62 Encoding
- Không dùng random → tránh collision
- Hiệu suất cao, chuỗi ngắn, đẹp
- Lookup nhanh do dựa trên ID

### Database (PostgreSQL)
- 1 bảng: `short_urls`
- Index: `unique(short_code)`, `unique(original_url)`
- Migration tự động bằng GORM

### Kiến trúc Clean Layered
- Handler (Gin)
- Service layer
- Repository
- Model (GORM)

---

## 3. System Architecture
<img width="1022" height="408" alt="Screenshot 2025-12-04 at 10 00 55" src="https://github.com/user-attachments/assets/33ebcab9-ddc2-4efc-8bfc-eaa12b5e52d4" />


---

## 4. API Flow

### Create Short URL Flow
<img width="767" height="282" alt="Screenshot 2025-12-04 at 18 35 23" src="https://github.com/user-attachments/assets/e7c7d60e-3f54-47d9-a363-f501ef6415a2" />


### Redirect Flow
<img width="784" height="430" alt="Screenshot 2025-12-04 at 18 36 53" src="https://github.com/user-attachments/assets/c7905d66-9529-4722-9064-7b89d67bee77" />


---

## 5. Database Schema
<img width="242" height="308" alt="Screenshot 2025-12-04 at 18 30 30" src="https://github.com/user-attachments/assets/948a8d70-394e-4877-9997-d6a9c3b27bd8" />

---

## 6. Implementation Overview
*(Sẽ bổ sung ở phiên bản 1.1)*

---

## 7. Limitations & Future Improvements
*(Sẽ hoàn thiện trong phiên bản 1.2)*


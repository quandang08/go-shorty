# Go Shorty – High Performance URL Shortener Service

## 1. Introduction

GoShorty là một hệ thống rút gọn URL được xây dựng với tiêu chí **đơn giản, rõ ràng và đúng chuẩn backend**.  
Dự án tập trung giải quyết hai vấn đề cốt lõi của bất kỳ URL Shortener chuyên nghiệp nào:

- **Tạo short code ngắn, không trùng và hiệu suất cao.**
- **Chuyển hướng (redirect) nhanh và chính xác**, đồng thời ghi nhận lượt click một cách an toàn và nhất quán.

Để đạt được điều đó, GoShorty lựa chọn kiến trúc tối ưu và dễ mở rộng:

- **Golang** – phù hợp cho high-performance API nhờ khả năng xử lý concurrency tự nhiên (goroutine).
- **PostgreSQL** – đảm bảo tính toàn vẹn dữ liệu (ACID), hỗ trợ unique index để chống trùng lặp và tránh race-condition.
- **Base62 Encoding trên ID tự tăng** – đảm bảo short code duy nhất, không cần random và truy vấn rất nhanh.

GoShorty được xây dựng như một **minimal but correct backend system**:  
gọn nhẹ, ổn định, dễ bảo trì, và hoạt động tốt ngay cả khi load tăng cao.

## 2. System Architecture
<img width="1022" height="408" alt="Screenshot 2025-12-04 at 10 00 55" src="https://github.com/user-attachments/assets/d06e30db-fe4b-4977-8535-b6be87c0ba9c" />

## 3. API Flow

#Create Short URL Flow
<img width="767" height="282" alt="Screenshot 2025-12-04 at 18 35 23" src="https://github.com/user-attachments/assets/a84ed7cf-3071-42b3-8a12-774a651c6978" />

#Redirect Flow
<img width="784" height="430" alt="Screenshot 2025-12-04 at 18 36 53" src="https://github.com/user-attachments/assets/e394c17d-48af-4dcf-a1b5-6eacac9835d8" />


## 4. Database Schema
<img width="242" height="308" alt="Screenshot 2025-12-04 at 18 30 30" src="https://github.com/user-attachments/assets/f7248540-ef65-420b-8cf9-720de2a55ea6" />

## 5. Implementation Overview
## 6. Limitations & Future Improvements

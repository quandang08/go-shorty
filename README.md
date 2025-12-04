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
## 3. API Flow
## 4. Database Schema
## 5. Implementation Overview
## 6. Limitations & Future Improvements

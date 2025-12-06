# Go Shorty – High Performance URL Shortener Service

![version](https://img.shields.io/badge/version-1.0-blue)
![status](https://img.shields.io/badge/status-stable-brightgreen)
![go](https://img.shields.io/badge/Go-1.22-%2300ADD8?logo=go)
![postgres](https://img.shields.io/badge/PostgreSQL-16-%234169E1?logo=postgresql)
![architecture](https://img.shields.io/badge/Architecture-Clean%20Layered-orange)
![license](https://img.shields.io/badge/License-MIT-green)

---

Introduction – GoShorty v1.0

GoShorty là một hệ thống rút gọn URL được xây dựng với trọng tâm hướng đến sự chủ động, độ ổn định và khả năng vận hành dài hạn. Dù được dùng cho dự án cá nhân, nội bộ doanh nghiệp hay tích hợp như một module trong kiến trúc lớn hơn, hệ thống vẫn hoạt động độc lập mà không phụ thuộc bất kỳ dịch vụ bên ngoài nào.

Ở phiên bản đầu tiên, GoShorty tập trung vào hai khả năng nền tảng của mọi URL shortener:
tạo short code duy nhất với hiệu năng cao và redirect chính xác tuyệt đối.
Không tính năng thừa. Không mơ hồ. Chỉ những gì cần để đảm bảo sự bền vững.

Bên cạnh tính kỹ thuật, rút gọn URL còn mang lại giá trị thực tế: dễ chia sẻ hơn, dễ quản lý hơn, giúp chuẩn hóa đường dẫn trong hệ thống nội bộ, và tạo ra lớp “trừu tượng” để không phơi bày các endpoint dài hoặc nhạy cảm. GoShorty giữ vai trò này bằng cách cung cấp short link gọn, rõ ràng và ổn định cho mọi nhu cầu từ vận hành đến truyền thông.

Việc tự xây thay vì dùng Bit.ly hay Rebrandly xuất phát từ nhu cầu kiểm soát hoàn toàn dữ liệu, tuân thủ chính sách nội bộ và tránh rủi ro phụ thuộc hạ tầng bên thứ ba. Mọi cơ chế sinh mã, lưu trữ và xử lý đều nằm trong tay người vận hành.

GoShorty sử dụng Base62 kết hợp ID tự tăng để tạo short code:

không thể trùng theo thiết kế

tốc độ truy vấn ổn định

dễ mở rộng theo chiều ngang

không tiêu tốn chi phí hashing hoặc random

dữ liệu bền vững trong nhiều năm

Cơ chế tạo mã đơn giản nhưng chắc chắn, và khi kết hợp với PostgreSQL, toàn bộ short link được lưu trữ một cách bền bỉ, nhất quán và không bị sai lệch, kể cả khi hệ thống chịu tải cao trong thời gian dài.

Dưới những tình huống xấu nhất như mất điện hoặc server sập đột ngột, database vẫn đảm bảo tính toàn vẹn. Người dùng chỉ gặp gián đoạn truy cập tạm thời — còn dữ liệu short link vẫn an toàn tuyệt đối và được khôi phục nguyên vẹn khi dịch vụ hoạt động trở lại.

Chỉ với một server nhỏ chạy Go + Postgres, hệ thống có thể xử lý hàng trăm nghìn request mỗi ngày với chi phí cực thấp, phù hợp để vận hành lâu dài mà không phải mở rộng hạ tầng quá mức.

GoShorty v1.0 được đặt trên nền móng tối giản nhưng nghiêm túc, và hướng đến một chặng đường dài hơn: mở rộng theo nhu cầu thực tế mà không phá vỡ kiến trúc hiện tại, bổ sung tính năng một cách có chủ đích, và phát triển thành một nền tảng rút gọn URL nhỏ gọn nhưng bền vững trong nhiều năm.

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


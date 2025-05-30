# Changelog

Tất cả những thay đổi đáng chú ý của MongoDB Provider sẽ được ghi lại trong file này.

Định dạng dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân theo [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-05-30

### Changed
- **BREAKING CHANGE**: Migrated module path from `github.com/go-fork/providers/mongodb` to `go.fork.vn/providers/mongodb`
- Updated all import statements to use the new module path
- Migrated to official go.fork.vn domain

### Updated
- Updated dependency `go.fork.vn/providers/config` to v0.1.0
- Updated dependency `go.fork.vn/di` to v0.1.0
- Updated documentation with new module path
- Updated README.md with new installation instructions

## [v0.0.2-next] - Unreleased


## [v0.0.1] 2025-05-27

### Added
- Cập nhật tài liệu chi tiết với file doc.go
- Thêm các ví dụ sử dụng đầy đủ trong README.md
- Cải thiện test coverage từ 80.7% lên 90.4%
- Tối ưu hóa quy trình kết nối và quản lý session
- Phiên bản ổn định đầu tiên
- Tích hợp với framework dependency injection go-fork
- Quản lý kết nối MongoDB và connection pooling
- Hỗ trợ xác thực và SSL/TLS
- Interface cho các thao tác MongoDB phổ biến
- Hỗ trợ transaction và change streams
- Tiện ích kiểm tra sức khỏe và thống kê
- Mock cho kiểm thử
- Cập nhật phụ thuộc MongoDB driver lên phiên bản mới nhất
- Sửa vấn đề rò rỉ bộ nhớ trong quản lý session
- Xử lý tốt hơn các lỗi không đồng bộ trong quản lý connection pool
- Thêm phương thức QueryContext để trả về kết quả dạng bson.M
- Thêm phương thức HealthCheck để kiểm tra trạng thái kết nối
- Hỗ trợ cấu hình TLS/SSL
- Cải thiện quản lý lỗi khi không thể kết nối tới MongoDB
- Tối ưu hóa quản lý connection pool
- Đóng kết nối đúng cách khi context bị hủy
- Phiên bản thử nghiệm đầu tiên
- Các chức năng cơ bản cho việc kết nối và thao tác với MongoDB
- Hỗ trợ cấu hình qua YAML
- Tích hợp với framework DI của go-fork

# Changelog

## [Unreleased]

## [v0.0.5] - 2025-05-29

### Changed
- Nâng cấp github.com/go-fork/di từ v0.0.4 lên v0.0.5
- Đổi tên method Provides() thành Providers() để phù hợp với interface di.ServiceProvider v0.0.5
- Triển khai method Requires() mới cho di.ServiceProvider v0.0.5
- Cập nhật tài liệu và ghi chú tương thích trong doc.go

### Added
- Bổ sung các bài test cho methods Providers() và Requires()
- Thêm thông tin tương thích di v0.0.5 vào README.md

### Fixed
- Cải thiện cách xử lý lỗi khi loading config

## [v0.0.4] - 2025-05-27

### Added
- Hỗ trợ proxy SMTP cho các môi trường doanh nghiệp
- Cải thiện hiệu năng khi gửi email hàng loạt
- Thêm tùy chọn retry cho các email bị lỗi

### Fixed
- Sửa lỗi khi gửi attachment có dung lượng lớn
- Xử lý đúng các kí tự đặc biệt trong tên email

## [v0.0.3] - 2025-05-25

### Added

- Tích hợp đầy đủ với DI container của ứng dụng
- API fluent và dễ sử dụng
- Hỗ trợ gửi cả email văn bản thuần túy và HTML
- Hỗ trợ render template từ `text/template` và `html/template`
- Hỗ trợ file đính kèm và nhúng hình ảnh
- Hỗ trợ xử lý hàng đợi (queue) cho email
- Dễ dàng test với MockMailer

## v0.0.3 - 2025-05-25

### Added

- Tích hợp đầy đủ với DI container của ứng dụng
- API fluent và dễ sử dụng
- Hỗ trợ gửi cả email văn bản thuần túy và HTML
- Hỗ trợ render template từ `text/template` và `html/template`
- Hỗ trợ file đính kèm và nhúng hình ảnh
- Hỗ trợ xử lý hàng đợi (queue) cho email
- Dễ dàng test với MockMailer

### Changed
- **BREAKING**: Upgraded from `github.com/go-gomail/gomail` to `gopkg.in/gomail.v2`
  - Updated all imports to use the new `gopkg.in/gomail.v2` package
  - The API remains the same, but the underlying dependency has been upgraded
  - This change provides better stability and official gopkg.in versioning

### Updated
- Documentation updated to reference `gopkg.in/gomail.v2` instead of `github.com/go-gomail/gomail`
- README.md updated with new dependency information
- All code comments and documentation reflect the new dependency

### Technical Details
- All `github.com/go-gomail/gomail` imports replaced with `gopkg.in/gomail.v2`
- No API changes required - the upgrade is backward compatible
- All existing tests continue to pass
- Core functionality remains unchanged

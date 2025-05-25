# Changelog

## [Unreleased]

## v0.0.3 - 2025-05-25

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

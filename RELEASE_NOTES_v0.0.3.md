# Release Notes v0.0.3

## Tính năng mới

- **Scheduler Provider**: Thêm mới Provider Scheduler để quản lý các tác vụ định kỳ
- **Queue Integration**: Hoàn thiện chức năng worker với tích hợp scheduler
- **Testing Utilities**: Thêm MockManager cho unit testing cấu hình và nâng cao queue tests

## Cải tiến

- **Mailer Package**: Cập nhật từ `github.com/go-gomail/gomail` sang `gopkg.in/gomail.v2` với API tương thích
- **Scheduler Enhancements**: Nâng cao chức năng scheduler với validation và các phương thức fluent interface

## Nâng cấp dependencies

- Di package: v0.0.2 → v0.0.3
- fsnotify: v1.8.0 → v1.9.0
- Các thư viện phụ thuộc khác đã được cập nhật lên phiên bản mới nhất

## Chuẩn bị tương lai

- Xóa module go.mod cho SMS, chuẩn bị cho việc cải tiến module này

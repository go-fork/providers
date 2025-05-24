# Release Notes v0.0.3

## Tính năng mới

- **Scheduler Provider**: Thêm mới Provider Scheduler để quản lý các tác vụ định kỳ
- **Queue Integration**: Hoàn thiện chức năng worker với tích hợp scheduler
- **Testing Utilities**: Thêm MockManager cho unit testing cấu hình và nâng cao queue tests

## Cải tiến

- **Mailer Package**: Cập nhật từ `github.com/go-gomail/gomail` sang `gopkg.in/gomail.v2` với API tương thích
- **Scheduler Enhancements**: Nâng cao chức năng scheduler với validation và các phương thức fluent interface
- **Config Package**: Bổ sung các API tiện ích mới cho việc truy xuất cấu hình

## Nâng cấp dependencies

- Di package: v0.0.2 → v0.0.3
- fsnotify: v1.8.0 → v1.9.0
- Các thư viện phụ thuộc khác đã được cập nhật lên phiên bản mới nhất

## Thay đổi cấu trúc

- **HTTP và Middleware**: Đã được chuyển sang repository riêng
- **SMS Provider**: Xóa module go.mod cho SMS, chuẩn bị cho việc cải tiến module này trong phiên bản tiếp theo

## Cài đặt

```bash
# Cài đặt từng provider riêng lẻ
go get github.com/go-fork/providers/config@v0.0.3
go get github.com/go-fork/providers/mailer@v0.0.3
go get github.com/go-fork/providers/queue@v0.0.3
go get github.com/go-fork/providers/scheduler@v0.0.3
go get github.com/go-fork/providers/log@v0.0.3

# Cài đặt tất cả providers
go get github.com/go-fork/providers/...@latest
```

## Tương thích

| Provider | v0.0.2 | v0.0.3 | di v0.0.2 | di v0.0.3 |
|----------|--------|--------|-----------|-----------|
| cache    | ✅     | -      | ✅        | ✅        |
| config   | ✅     | ✅     | ✅        | ✅        |
| log      | ✅     | ✅     | ✅        | ✅        |
| mailer   | ✅     | ✅     | ✅        | ✅        |
| queue    | ✅     | ✅     | ✅        | ✅        |
| scheduler| -      | ✅     | ❌        | ✅        |
| sms      | ✅     | -      | ✅        | ❌        |

✅ = Tương thích đầy đủ
❌ = Không tương thích
\- = Không có phiên bản

## Liên kết

- [CHANGELOG](CHANGELOG.md) - Danh sách đầy đủ các thay đổi
- [GitHub Repository](https://github.com/go-fork/providers) - Mã nguồn và issues

---
*Ngày phát hành: 25/05/2025*

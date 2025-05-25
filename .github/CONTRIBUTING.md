# Hướng dẫn đóng góp cho Fork Providers

Cảm ơn bạn đã quan tâm đến việc đóng góp cho Fork Providers! Dưới đây là các bước và quy tắc để đóng góp vào dự án.

## Quy trình đóng góp

1. Fork repository
2. Tạo nhánh tính năng (`git checkout -b feature/amazing-feature`)
3. Commit các thay đổi của bạn (`git commit -m 'feat: Add some amazing feature'`)
4. Push lên nhánh (`git push origin feature/amazing-feature`)
5. Mở Pull Request

## Quy tắc Commit

Chúng tôi sử dụng [Conventional Commits](https://www.conventionalcommits.org/) cho các commit message:

- `feat`: Thêm tính năng mới
- `fix`: Sửa lỗi
- `docs`: Thay đổi tài liệu
- `style`: Thay đổi không ảnh hưởng đến ý nghĩa code (khoảng trắng, định dạng, etc)
- `refactor`: Thay đổi code không sửa lỗi hoặc thêm tính năng
- `perf`: Thay đổi cải thiện hiệu suất
- `test`: Thêm hoặc sửa test
- `chore`: Thay đổi quá trình build hoặc công cụ auxiliary

## Quy tắc code

- Tuân thủ các quy ước đặt tên của Go
- Đảm bảo code của bạn vượt qua linter (`golangci-lint run`)
- Viết tests cho tất cả các chức năng mới
- Đảm bảo coverage ít nhất 70%
- Thêm comments theo chuẩn godoc

## Tạo Issue

- Sử dụng template có sẵn
- Cung cấp thông tin chi tiết nhất có thể
- Nếu có thể, hãy bao gồm các bước để tái hiện vấn đề

## Tiêu chuẩn Pull Request

- PR nên nhỏ và tập trung vào một thay đổi cụ thể
- Mô tả rõ ràng những gì PR thay đổi
- Tham chiếu đến issue tương ứng (nếu có)
- Đảm bảo tất cả các bài kiểm tra tự động vượt qua

## Quy trình review

- Tối thiểu một maintainer sẽ phải review và chấp thuận PR
- Các comment phải được giải quyết trước khi merge
- Chúng tôi có thể yêu cầu thay đổi trước khi chấp nhận PR

## Cấu trúc dự án

Dưới đây là cấu trúc thư mục của dự án:

```
providers/
├── adaptor/      # Các adaptor cho các thành phần khác nhau
├── cache/        # Quản lý cache
├── config/       # Quản lý cấu hình
├── http/         # Xử lý HTTP
├── log/          # Quản lý log
├── mailer/       # Dịch vụ email
├── middleware/   # Middleware HTTP
├── mongodb/      # Kết nối MongoDB
├── queue/        # Quản lý hàng đợi
├── scheduler/    # Lập lịch tác vụ
└── sms/          # Dịch vụ SMS
```

## Phát triển mới

Khi thêm một provider hoặc driver mới:

1. Đảm bảo tuân thủ các interface hiện có
2. Luôn viết tests đầy đủ (unit test, mock nếu cần)
3. Thêm tài liệu cho thành phần mới (doc.go, README, ví dụ nếu có)
4. Đảm bảo code không có lỗi linter và test pass 100%
5. Nếu thay đổi ảnh hưởng hệ thống, cập nhật CHANGELOG.md theo chuẩn (Added/Changed/Fixed/Removed)

## Liên hệ

Nếu bạn có bất kỳ câu hỏi nào, vui lòng tạo một issue hoặc liên hệ với maintainers.

Cảm ơn sự đóng góp của bạn!

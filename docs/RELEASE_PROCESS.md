# Quy trình Phát hành (Release) cho Go-Fork Providers

Tài liệu này mô tả chi tiết quy trình phát hành các module trong repository `github.com/go-fork/providers`.

## Cấu trúc Repository

Repository này tuân theo mô hình "multi-module repository", trong đó có nhiều Go modules độc lập:

- `github.com/go-fork/providers/cache`
- `github.com/go-fork/providers/config`
- `github.com/go-fork/providers/log`
- `github.com/go-fork/providers/mailer`
- `github.com/go-fork/providers/queue`
- `github.com/go-fork/providers/scheduler`
- `github.com/go-fork/providers/sms`
- Và các module khác...

Mỗi module có file `go.mod` riêng và có thể được đánh phiên bản và phát hành độc lập.

## Hệ thống Phiên bản

Chúng ta sử dụng [Semantic Versioning (SemVer)](https://semver.org/) cho tất cả các modules:

- **Major version (X)**: Khi có các thay đổi không tương thích ngược (breaking changes)
- **Minor version (Y)**: Khi thêm tính năng mới nhưng vẫn tương thích ngược
- **Patch version (Z)**: Khi sửa lỗi, tối ưu mà không thêm tính năng mới

### Quy ước đặt tên tag

- **Phát hành toàn bộ repository**: `vX.Y.Z` (ví dụ: `v0.1.0`)
- **Phát hành module riêng lẻ**: `module/vX.Y.Z` (ví dụ: `cache/v0.1.0`)

## Quản lý CHANGELOG

Mỗi module duy trì file `CHANGELOG.md` riêng để theo dõi các thay đổi cụ thể cho module đó. Repository chính cũng có một file `CHANGELOG.md` ở thư mục gốc tóm tắt các thay đổi trên tất cả các modules.

- CHANGELOG của module chứa thông tin chi tiết về các thay đổi cho module cụ thể đó
- CHANGELOG gốc cung cấp liên kết đến tất cả CHANGELOG của module và tóm tắt tổng quan

## Chuẩn bị Phát hành

Trước khi phát hành, hãy đảm bảo:

- Tất cả các thay đổi đã được commit
- CI pipeline đã chạy thành công
- Các bài kiểm tra (tests) đã chạy thành công
- CHANGELOG.md của mỗi module đã được cập nhật
- Kiểm tra các vấn đề tương thích giữa các module

## Quy trình Phát hành

### 1. Kiểm tra Tương thích

Trước khi phát hành nhiều module, nên kiểm tra các vấn đề tương thích:

```bash
# Kiểm tra tương thích giữa các module
./scripts/check_compatibility.sh
```

Script này sẽ:
- Phân tích các phụ thuộc giữa các module
- Kiểm tra sự không khớp về phiên bản
- Phát hiện các phụ thuộc vòng tròn
- Báo cáo các vấn đề được tìm thấy

### 2. Phát hành Module Riêng lẻ

Sử dụng script release:

```bash
# Phát hành module cache phiên bản v0.1.0
./scripts/release.sh --module cache --version v0.1.0
```

Script này sẽ:
- Kiểm tra tính hợp lệ của module
- Tự động tạo hoặc cập nhật file CHANGELOG.md cho module với:
  - Thông tin tính năng trích xuất từ README.md
  - Các thay đổi từ commit messages
- Tạo Git tag theo định dạng `module/version`: `cache/v0.1.0`
- Tùy chọn tạo GitHub Release nếu sử dụng flag `--create-release`

### 3. Phát hành Nhiều Module Cùng lúc

```bash
# Phát hành tất cả các modules với cùng phiên bản v0.1.0
./scripts/release.sh --all v0.1.0
```

Khi sử dụng tùy chọn này, script sẽ:
- Kiểm tra tương thích giữa các module trước khi phát hành
- Cập nhật CHANGELOG.md chính với thông tin tổng hợp từ tất cả module
- Phát hành từng module riêng biệt với cùng phiên bản
- Tạo các tag cho từng module theo định dạng `module/version`

Sau khi script hoàn thành, bạn cần đẩy tất cả các tag lên remote:
```bash
git push origin --tags
```

### 4. Phát hành Toàn bộ Repository

```bash
# Phát hành repository với phiên bản v0.1.0
./scripts/release.sh --repo v0.1.0
```

Script này sẽ:
- Cập nhật CHANGELOG.md ở thư mục gốc
- Tạo một tag chung cho toàn bộ repository
- Tùy chọn tạo GitHub Release cho toàn bộ repository

Sau đó, đẩy tag lên remote:
```bash
git push origin v0.1.0
```

## Các Tùy chọn Script Release

Script `release.sh` có nhiều tùy chọn hữu ích:

```bash
./scripts/release.sh --help
```

Các tùy chọn chính:

- `-m, --module MODULE`: Chỉ định module cần phát hành (ví dụ: cache, log)
- `-v, --version VERSION`: Chỉ định phiên bản phát hành (ví dụ: v0.1.0)
- `-a, --all VERSION`: Phát hành tất cả các module với cùng phiên bản
- `-r, --repo VERSION`: Phát hành toàn bộ repository với phiên bản đã cho
- `-c, --create-release`: Tạo GitHub Release ngoài việc tạo tags
- `-f, --force`: Bỏ qua kiểm tra các thay đổi chưa được commit
- `-o, --overwrite`: Ghi đè tag hiện có nếu đã tồn tại
- `-p, --push-only`: Chỉ đẩy tag hiện có, không tạo tag mới
- `-g, --generate-changelog`: Tạo changelog từ các thông điệp commit

## Ví dụ Sử dụng

```bash
# Phát hành module cache phiên bản v0.1.0 và tạo GitHub Release
./scripts/release.sh --module cache --version v0.1.0 --create-release

# Phát hành tất cả các module với v0.1.0, bỏ qua kiểm tra thay đổi chưa commit
./scripts/release.sh --all v0.1.0 --force

# Phát hành module config v0.1.0, ghi đè tag nếu đã tồn tại
./scripts/release.sh --module config --version v0.1.0 --overwrite

# Chỉ đẩy tag hiện có lên remote mà không tạo tag mới
./scripts/release.sh --module config --version v0.1.0 --push-only

# Tự động tạo changelog từ thông điệp commit
./scripts/release.sh --module config --version v0.1.0 --generate-changelog
```

## Quản lý Changelog Tự động

Hệ thống release hiện hỗ trợ các tính năng nâng cao để quản lý CHANGELOG:

1. **Tạo CHANGELOG tự động cho module mới**: Nếu một module chưa có file CHANGELOG.md, script sẽ tự động tạo một file mới với cấu trúc phù hợp.

2. **Trích xuất tính năng từ README**: Script sẽ trích xuất thông tin tính năng từ phần "Tính năng" hoặc "Features" trong file README.md của module.

3. **Tổng hợp CHANGELOG của repository**: Khi phát hành nhiều module, script `update_main_changelog.sh` sẽ tự động tổng hợp thông tin từ tất cả CHANGELOG của module vào CHANGELOG chính.

4. **Phát hiện các thay đổi giữa các phiên bản**: Script sẽ phân tích các commit giữa phiên bản hiện tại và phiên bản trước đó để tạo changelog chi tiết.

## Quy trình CI/CD cho Phát hành

Khi một tag được đẩy lên GitHub, các workflow sau sẽ được kích hoạt:

1. **Module Release Workflow** (cho tags `*/v*`):
   - Kích hoạt khi tag có định dạng `module/vX.Y.Z`
   - Chạy tests cho module cụ thể
   - Tạo GitHub Release cho module đó

2. **Repository Release Workflow** (cho tags `v*`):
   - Kích hoạt khi tag có định dạng `vX.Y.Z`
   - Tạo GitHub Release cho toàn bộ repository

## Quản lý Phiên bản Module

### Module với Major Version > 1

Đối với modules có major version > 1, tuân theo quy ước Go Modules:

```
github.com/go-fork/providers/module/v2
github.com/go-fork/providers/module/v3
...
```

Mỗi major version sẽ có thư mục riêng với hậu tố `/vX`.

## Kiểm tra Tương thích

Trước khi phát hành, kiểm tra tương thích ngược:

```bash
# Sử dụng apidiff để kiểm tra thay đổi API
go install golang.org/x/exp/cmd/apidiff@latest
cd module
apidiff -incompatible ./... origin/main
```

## Hướng dẫn chọn Phiên bản

- **X.0.0** (Major): Khi có thay đổi API không tương thích ngược (breaking changes)
- **0.Y.0** (Minor): Khi thêm tính năng mới tương thích ngược
- **0.0.Z** (Patch): Khi sửa lỗi, cải thiện hiệu suất mà không thêm tính năng mới

## Sử dụng Releases

### Import module

```go
import "github.com/go-fork/providers/cache" // Phiên bản mới nhất
import "github.com/go-fork/providers/cache/v2" // Major version 2
```

### Go get cho module riêng lẻ

```bash
# Sử dụng định dạng tag mới module/vX.Y.Z
go get github.com/go-fork/providers/cache@cache/v0.1.0
go get github.com/go-fork/providers/config@config/v0.2.0
```

### Go get toàn bộ repo (tất cả module)

```bash
# Lấy tất cả các module với cùng phiên bản
go get github.com/go-fork/providers/...@v0.1.0
```

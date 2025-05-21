# Quy trình Release cho Go-Fork Providers

Tài liệu này mô tả chi tiết quy trình release các module trong repository `github.com/go-fork/providers`.

## Cấu trúc Repository

Repository này tuân theo mô hình "multi-module repository", trong đó có nhiều Go modules độc lập:

- `github.com/go-fork/providers/cache`
- `github.com/go-fork/providers/config`
- `github.com/go-fork/providers/http`
- `github.com/go-fork/providers/log`
- `github.com/go-fork/providers/middleware`
- Và các module khác...

Mỗi module có file `go.mod` riêng và có thể được version, release độc lập.

## Hệ thống đánh phiên bản

Chúng ta sử dụng [Semantic Versioning (SemVer)](https://semver.org/) cho tất cả các modules:

- **Major version (X)**: Khi có các thay đổi không tương thích ngược (breaking changes)
- **Minor version (Y)**: Khi thêm tính năng mới nhưng vẫn tương thích ngược
- **Patch version (Z)**: Khi sửa lỗi, tối ưu mà không thêm tính năng mới

### Convention đánh tag

- **Release toàn bộ repository**: `vX.Y.Z` (ví dụ: `v0.1.0`)
- **Release module riêng lẻ**: `module/vX.Y.Z` (ví dụ: `cache/v0.1.0`)

## Các bước release

### 1. Chuẩn bị release

Trước khi release, đảm bảo:

- Tất cả các thay đổi đã được commit
- CI pipeline đã pass
- Tests đã pass
- CHANGELOG.md đã được cập nhật

### 2. Release module riêng lẻ

Sử dụng script release:

```bash
# Release module cache phiên bản v0.1.0
./scripts/release.sh --module cache --version v0.1.0
```

Script sẽ:
- Kiểm tra tính hợp lệ của module
- Cập nhật CHANGELOG.md (nếu có)
- Tạo Git tag: `cache/v0.1.0`

Sau đó, push tag để trigger GitHub Actions:

```bash
git push origin cache/v0.1.0
```

### 3. Release nhiều module cùng lúc

```bash
# Release tất cả các modules với cùng phiên bản v0.1.0
./scripts/release.sh --all v0.1.0

# Push tất cả các tags
git push origin --tags
```

### 4. Release toàn bộ repository

```bash
# Release repository với phiên bản v0.1.0
./scripts/release.sh --repo v0.1.0

# Push tag
git push origin v0.1.0
```

## CI/CD Pipeline cho releases

Khi một tag được push lên GitHub, các workflow sau sẽ được trigger:

1. **Module Release Workflow** (cho tags `*/v*`):
   - Chạy khi tag có format `module/vX.Y.Z`
   - Chạy tests cho module cụ thể
   - Tạo GitHub Release cho module đó

2. **Repository Release Workflow** (cho tags `v*`):
   - Chạy khi tag có format `vX.Y.Z`
   - Sử dụng GoReleaser để tạo GitHub Release cho toàn bộ repository

## Quản lý Changelog

- Mỗi module *nên* có file `CHANGELOG.md` riêng
- Repository cũng *nên* có một `CHANGELOG.md` ở root
- Script release sẽ tự động cập nhật CHANGELOG với phiên bản mới

## Quản lý phiên bản module

### Module với Major Version > 1

Đối với modules có major version > 1, tuân theo Go Modules conventions:

```
github.com/go-fork/providers/module/v2
github.com/go-fork/providers/module/v3
...
```

Mỗi major version sẽ có directory riêng với suffix `/vX`.

## Kiểm tra tương thích

Trước khi release, kiểm tra tương thích ngược:

```bash
# Sử dụng apidiff để kiểm tra thay đổi API
go install golang.org/x/exp/cmd/apidiff@latest
cd module
apidiff -incompatible ./... origin/main
```

## Guideline chọn phiên bản

- **Patch (vX.Y.Z → vX.Y.Z+1)**: Bug fixes không thay đổi API
- **Minor (vX.Y.Z → vX.Y+1.0)**: Thêm tính năng, không có breaking changes
- **Major (vX.Y.Z → vX+1.0.0)**: Có breaking changes

## Sử dụng Releases

### Import module

```go
import "github.com/go-fork/providers/cache" // Phiên bản mới nhất
import "github.com/go-fork/providers/cache/v2" // Major version 2
```

### Go get

```bash
go get github.com/go-fork/providers/cache@v0.1.0
go get github.com/go-fork/providers/config@v0.2.0
```

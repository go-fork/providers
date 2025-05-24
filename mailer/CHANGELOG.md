# Changelog

## [Unreleased]

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

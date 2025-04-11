# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - Unreleased

### Added

- Snapshot testing system with content-type aware formatters
- Support for JSON, XML, HTML, plain text, and binary response comparison
- Snapshot management commands (test, update, list, cleanup)
- Flexible snapshot update modes (none, all, failed, missing)
- Comprehensive diffing between responses with detailed output
- Configurable header ignoring for reliable snapshot comparisons

## [0.1.0] - 2025-04-11

### Added

- Core domain models for Swagger/OpenAPI documents
- Core domain models for HTTP request files
- Swagger parser with support for JSON and YAML formats
- HTTP file generator with proper formatting
- File writer for generating organized files and directories
- Configuration system with environment variable support
- Command-line interface with generate command
- Clean architecture implementation throughout

### Changed

- Initial release

### Removed

- Nothing removed

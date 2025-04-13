# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - Unreleased

### Added

- Schema validation against OpenAPI specifications
- Sequential tests with dependencies between steps
- Variable extraction from response body, headers, and status
- Test assertions with multiple validation types (equals, contains, matches, etc.)
- Support for watching file changes and running tests continuously
- Advanced test runner with combined features
- Test sequence file format for workflow testing
- Expanded CLI with advanced testing commands
- Comprehensive documentation for advanced testing features
- Integration examples for CI/CD pipelines

### Changed

- Enhanced test runner to support schema validation
- Improved test reporting with more detailed results
- Updated CLI with more flexible options
- Expanded README with advanced testing feature descriptions

## [0.2.0] - Unreleased

### Added

- HTTP executor for executing requests from .http files
- Support for variable substitution in URLs, headers, and body content
- Environment variable loading for configurable requests
- HTTP file parser with robust format support
- HTTP timeout configuration for reliable testing
- Comprehensive error handling for network issues
- Snapshot testing system with content-type aware formatters
- Support for JSON, XML, HTML, plain text, and binary response comparison
- Snapshot management commands (test, update, list, cleanup)
- Flexible snapshot update modes (none, all, failed, missing)
- Comprehensive diffing between responses with detailed output
- Configurable header ignoring for reliable snapshot comparisons
- Complete HTTP request lifecycle management
- Comprehensive documentation and examples
- API reference for developers
- Detailed guides for installation, usage, and configuration
- Complete examples for different usage scenarios
- FAQ section for common questions

### Changed

- Improved README with detailed information and examples
- Enhanced command-line help text
- Added proper code documentation
- Updated snapshot command to use real HTTP executor

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

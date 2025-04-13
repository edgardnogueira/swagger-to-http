# Models Package

This package contains domain models for the swagger-to-http tool. The models are organized into logical groupings to improve maintainability and avoid duplications.

## Model Organization

- **http_models.go**: Contains all HTTP-related models (requests, responses, headers, etc.)
- **snapshot_models.go**: Contains all snapshot-related models (diffs, comparison results, etc.)
- **validation_models.go**: Contains all schema validation models
- **test_models.go**: Contains all test-related models (sequences, assertions, reports, etc.)
- **converters.go**: Contains conversion utilities between different model formats
- **utils.go**: Contains general string and data conversion utilities

## Model Compatibility

To maintain backward compatibility while consolidating models, we've implemented several strategies:

1. **Compatibility Fields**: Core models like `HTTPRequest` and `SnapshotResult` include fields needed for both new and legacy code.
2. **Conversion Methods**: Methods like `ToHTTPFileRequest()` and `ToHTTPRequest()` allow seamless conversion between different model formats.
3. **Helper Functions**: Utility functions in `utils.go` help with common type conversion issues (string/byte conversions, etc.)
4. **Field Synchronization**: The `SyncCompatibilityFields()` method keeps duplicate fields synchronized.

## Deprecated Files

The following files are kept for backward compatibility but are deprecated and should not be used for new development:

- **http.go**: Use http_models.go instead
- **snapshot.go**: Use snapshot_models.go instead
- **schema_validation.go**: Use validation_models.go instead
- **test_report.go**: Use test_models.go instead
- **test_sequence.go**: Use test_models.go instead

## Model Ownership

To avoid duplicate declarations, each model should be defined in exactly one file. When adding new models or making changes, make sure to follow this organization to maintain code clarity.

## Best Practices

1. When adding a new model, carefully consider which file it should go in
2. When extending existing models, update the original definition rather than creating a duplicate
3. Use meaningful comments to explain model fields and their purpose
4. Maintain consistent naming and formatting conventions
5. Use converters and utilities to handle incompatible types instead of duplicating code
6. When a model has compatibility fields, ensure they are kept in sync using the appropriate method

## Field Compatibility Guidelines

When working with models that have compatibility issues:

1. Use the `ToHTTPRequest()` and `ToHTTPFileRequest()` methods to convert between HTTP model formats
2. Always call `SyncCompatibilityFields()` after modifying `SnapshotResult` fields
3. Use `StringToBytes()` and `BytesToString()` for string/byte conversions
4. When working with headers, use appropriate conversion functions based on the expected format
5. Prefer the newer consolidated model structures for new development

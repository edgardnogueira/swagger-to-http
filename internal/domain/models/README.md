# Models Package

This package contains domain models for the swagger-to-http tool. The models are organized into logical groupings to improve maintainability and avoid duplications.

## Model Organization

- **http_models.go**: Contains all HTTP-related models (requests, responses, headers, etc.)
- **snapshot_models.go**: Contains all snapshot-related models (diffs, comparison results, etc.)
- **validation_models.go**: Contains all schema validation models
- **test_models.go**: Contains all test-related models (sequences, assertions, reports, etc.)

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

# Configuration Guide

This guide explains how to configure `swagger-to-http` using configuration files and environment variables.

## Configuration Sources

`swagger-to-http` uses the following configuration lookup paths (in order of precedence):

1. Command-line flags (highest precedence)
2. Environment variables with the `STH_` prefix
3. Configuration files:
   - `./swagger-to-http.yaml` (current directory)
   - `$HOME/.swagger-to-http/swagger-to-http.yaml` (home directory)
   - `/etc/swagger-to-http/swagger-to-http.yaml` (system-wide)

## Configuration File Format

The configuration file should be in YAML format. Here's a complete example:

```yaml
output:
  directory: http-requests
  
generator:
  indent_json: true
  include_auth: false
  auth_header: Authorization
  auth_token: ""
  default_tag: default
  base_url: ""
  
snapshots:
  directory: .snapshots
  update_mode: none  # none, all, failed, missing
  ignore_headers:
    - Date
    - Set-Cookie
  fail_on_missing: false
  cleanup_after_run: false
```

## Configuration Options

### Output Options

| File Key | Env Variable | CLI Flag | Description | Default |
|----------|--------------|----------|-------------|---------|
| `output.directory` | `STH_OUTPUT_DIRECTORY` | `-o, --output` | Directory for HTTP file output | `http-requests` |

### Generator Options

| File Key | Env Variable | CLI Flag | Description | Default |
|----------|--------------|----------|-------------|---------|
| `generator.indent_json` | `STH_INDENT_JSON` | `-i, --indent-json` | Indent JSON in request bodies | `true` |
| `generator.include_auth` | `STH_INCLUDE_AUTH` | `--auth` | Include authentication header in requests | `false` |
| `generator.auth_header` | `STH_AUTH_HEADER` | `--auth-header` | Authentication header name | `Authorization` |
| `generator.auth_token` | `STH_AUTH_TOKEN` | `--auth-token` | Authentication token value | `""` |
| `generator.default_tag` | `STH_DEFAULT_TAG` | `-t, --default-tag` | Default tag for operations without tags | `default` |
| `generator.base_url` | `STH_BASE_URL` | `-b, --base-url` | Base URL for requests | `""` |

### Snapshot Options

| File Key | Env Variable | CLI Flag | Description | Default |
|----------|--------------|----------|-------------|---------|
| `snapshots.directory` | `STH_SNAPSHOT_DIRECTORY` | `--snapshot-dir` | Directory for snapshot storage | `.snapshots` |
| `snapshots.update_mode` | `STH_UPDATE_MODE` | `--update` | Update mode for snapshots | `none` |
| `snapshots.ignore_headers` | `STH_IGNORE_HEADERS` | `--ignore-headers` | Headers to ignore in comparison | `["Date", "Set-Cookie"]` |
| `snapshots.fail_on_missing` | `STH_FAIL_ON_MISSING` | `--fail-on-missing` | Fail when snapshot is missing | `false` |
| `snapshots.cleanup_after_run` | `STH_CLEANUP_AFTER_RUN` | `--cleanup` | Remove unused snapshots after testing | `false` |

## Environment Variables

All configuration options can be set using environment variables with the `STH_` prefix. For nested options in the YAML file, use underscores.

Examples:

```bash
# Set output directory
export STH_OUTPUT_DIRECTORY=api-requests

# Set authentication token
export STH_AUTH_TOKEN="Bearer mytoken123"

# Set snapshot update mode
export STH_UPDATE_MODE=failed

# Set multiple headers to ignore (as a comma-separated list)
export STH_IGNORE_HEADERS="Date,Set-Cookie,X-Request-ID"
```

## Configuration Precedence

When multiple configuration sources define the same option, the following order of precedence applies:

1. Command-line flags have the highest precedence
2. Environment variables override configuration file settings
3. Configuration file settings override defaults

## Example Configurations

### Basic Configuration

```yaml
output:
  directory: http-requests
generator:
  indent_json: true
  default_tag: api
```

### Authentication Configuration

```yaml
generator:
  include_auth: true
  auth_header: Authorization
  auth_token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Snapshot Testing Configuration

```yaml
snapshots:
  directory: .snapshots
  update_mode: failed
  ignore_headers:
    - Date
    - Set-Cookie
    - X-Request-ID
    - X-Rate-Limit-Remaining
```

### CI/CD Configuration

```yaml
output:
  directory: http-requests
snapshots:
  directory: .snapshots
  update_mode: none
  fail_on_missing: true
  cleanup_after_run: true
```

## Next Steps

- Read the [Usage Guide](usage.md) to see how to use these configuration options
- Learn about the [HTTP File Format](http-file-format.md)

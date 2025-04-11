# Installation Guide

This guide covers the various ways to install the `swagger-to-http` tool.

## System Requirements

- Go 1.21 or higher (for source installation)
- Any operating system: Linux, macOS, Windows

## Installation Methods

### Using Go

The recommended way to install `swagger-to-http` is using Go:

```bash
go install github.com/edgardnogueira/swagger-to-http@latest
```

This will download, compile, and install the latest version of the tool. The binary will be placed in your `$GOPATH/bin` directory, which should be in your system's PATH.

### Binary Releases

Pre-built binaries are available for various platforms on the [Releases page](https://github.com/edgardnogueira/swagger-to-http/releases).

1. Download the appropriate binary for your platform:
   - Linux: `swagger-to-http_linux_amd64.tar.gz`
   - macOS: `swagger-to-http_darwin_amd64.tar.gz` or `swagger-to-http_darwin_arm64.tar.gz` (for Apple Silicon)
   - Windows: `swagger-to-http_windows_amd64.zip`

2. Extract the archive:
   ```bash
   # Linux/macOS
   tar -xzf swagger-to-http_*.tar.gz
   
   # Windows (using PowerShell)
   Expand-Archive swagger-to-http_*.zip -DestinationPath .
   ```

3. Move the binary to a location in your PATH:
   ```bash
   # Linux/macOS
   sudo mv swagger-to-http /usr/local/bin/
   
   # Windows
   # Move to a directory in your PATH, such as C:\Windows\System32\
   ```

### Using Docker

For containerized environments, you can use the Docker image:

```bash
# Pull the image
docker pull edgardnogueira/swagger-to-http

# Run the container
docker run -v $(pwd):/app edgardnogueira/swagger-to-http generate -f /app/swagger.json -o /app/http-requests
```

The Docker image is suitable for CI/CD pipelines or environments where installing binaries isn't preferred.

### Building from Source

To build from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/edgardnogueira/swagger-to-http.git
   cd swagger-to-http
   ```

2. Build the binary:
   ```bash
   go build -o swagger-to-http ./cmd/swagger-to-http
   ```

3. Install the binary:
   ```bash
   sudo mv swagger-to-http /usr/local/bin/
   ```

## Verifying the Installation

After installation, verify that the tool is correctly installed:

```bash
swagger-to-http --version
```

You should see output displaying the current version of the tool.

## Next Steps

Now that you have installed `swagger-to-http`, proceed to the [Usage Guide](usage.md) to learn how to use the tool.

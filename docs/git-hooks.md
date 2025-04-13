# Git Hooks Integration

The swagger-to-http tool includes a powerful Git hooks integration system that automates the process of validating Swagger/OpenAPI files and generating HTTP files as part of your Git workflow.

## Features

- **Pre-commit hook**: Validates Swagger/OpenAPI files and automatically generates corresponding HTTP files before each commit
- **Post-merge hook**: Updates HTTP files when Swagger/OpenAPI files have changed after a merge or pull
- **Selective updates**: Only regenerates HTTP files for endpoints that have changed
- **Cross-platform support**: Works on both Linux/macOS (Bash) and Windows (PowerShell)
- **Node.js/Husky integration**: Seamless integration with Node.js projects using Husky
- **Configurable behavior**: Customize the behavior of the hooks through configuration files

## Installation

### Linux/macOS

```bash
# Navigate to your project directory
cd your-project-directory

# Make the install script executable
chmod +x hooks/install.sh

# Run the installation script
./hooks/install.sh
```

### Windows

```powershell
# Navigate to your project directory
cd your-project-directory

# Run the installation script (may need to adjust execution policy)
PowerShell -ExecutionPolicy Bypass -File hooks/install.ps1
```

## Configuration

After installation, a configuration file is created at `.swagger-to-http/hooks.config` (Bash) or `.swagger-to-http/hooks.config.ps1` (PowerShell). You can edit this file to customize the behavior of the hooks.

### Configuration Options

| Option | Description | Default |
|--------|-------------|--------|
| `HOOKS_ENABLED` | Enable or disable all hooks | `true` |
| `SWAGGER_FILE_PATTERNS` | Patterns for identifying Swagger/OpenAPI files | `**/swagger.json **/swagger.yaml **/openapi.json **/openapi.yaml` |
| `HTTP_OUTPUT_DIR` | Output directory for generated HTTP files | `http` |
| `VALIDATE_SWAGGER` | Whether to validate Swagger/OpenAPI files before generating HTTP files | `true` |
| `SELECTIVE_UPDATES` | Whether to regenerate all HTTP files on changes or only affected ones | `true` |

## Temporary Disabling Hooks

You can temporarily disable the hooks by setting the `SWAGGER_TO_HTTP_DISABLE_HOOKS` environment variable:

### Linux/macOS

```bash
# Disable hooks for a single command
SWAGGER_TO_HTTP_DISABLE_HOOKS=true git commit -m "Skip hooks for this commit"

# Disable hooks for the current session
export SWAGGER_TO_HTTP_DISABLE_HOOKS=true
```

### Windows

```powershell
# Disable hooks for a single command
$env:SWAGGER_TO_HTTP_DISABLE_HOOKS="true"; git commit -m "Skip hooks for this commit"

# Disable hooks for the current session
$env:SWAGGER_TO_HTTP_DISABLE_HOOKS="true"
```

## Integration with Node.js and Husky

If you have a Node.js project using Husky for Git hooks management, the installation script automatically configures Husky to use the swagger-to-http hooks.

If you want to manually add the hooks to an existing Husky setup:

```bash
# For pre-commit hook
npx husky add .husky/pre-commit "hooks/pre-commit.sh"

# For post-merge hook
npx husky add .husky/post-merge "hooks/post-merge.sh"
```

## Troubleshooting

### Hooks Not Running

1. Check if Git hooks are enabled in your Git configuration:
   ```bash
   git config core.hooksPath
   ```
   This should point to `.git/hooks` or your Husky hooks directory.

2. Verify the hook files have execute permissions (Linux/macOS):
   ```bash
   chmod +x .git/hooks/pre-commit .git/hooks/post-merge
   ```

3. Check if the hooks are disabled in the configuration file `.swagger-to-http/hooks.config`.

### Command Not Found Errors

Ensure the `swagger-to-http` binary is in your system PATH or in one of the following locations:
- `./bin/swagger-to-http`
- `./swagger-to-http`

## How It Works

### Pre-commit Hook

The pre-commit hook runs whenever you make a commit and performs the following actions:

1. Identifies Swagger/OpenAPI files that are being committed
2. Validates those files (if enabled)
3. Generates or updates corresponding HTTP files
4. Automatically stages the generated HTTP files for commit

### Post-merge Hook

The post-merge hook runs after a merge or pull operation and performs the following actions:

1. Identifies Swagger/OpenAPI files that were changed in the merge
2. Generates or updates corresponding HTTP files
3. Notifies you to review and commit the changes if needed

### Selective Updates

When selective updates are enabled, the tool only regenerates HTTP files for endpoints that have changed, which improves performance and reduces unnecessary file changes. The tool compares the previous and current versions of the Swagger/OpenAPI file to determine what has changed.

## Advanced Usage

### Custom Hook Installation

If you need to integrate the hooks into an existing Git hooks system, you can manually install them:

1. Copy the hook scripts to your hooks directory
2. Modify your existing hooks to call the swagger-to-http hooks

### CI/CD Integration

In CI/CD environments, you may want to disable the interactive features of the hooks. Set the `SWAGGER_TO_HTTP_DISABLE_HOOKS` environment variable to `true` in your CI/CD environment to skip the hooks during automated builds and tests.

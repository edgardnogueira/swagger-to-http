# PowerShell post-merge hook for swagger-to-http
# Updates HTTP files when Swagger/OpenAPI files have changed after a merge/pull

$ErrorActionPreference = "Stop"

$EchoPrefix = "[swagger-to-http]"
$ErrorPrefix = "[swagger-to-http ERROR]"
$WarningPrefix = "[swagger-to-http WARNING]"

# Check if the environment variable to disable hooks is set
if ($env:SWAGGER_TO_HTTP_DISABLE_HOOKS -eq "true") {
    Write-Host "$WarningPrefix Git hooks are disabled via SWAGGER_TO_HTTP_DISABLE_HOOKS" -ForegroundColor Yellow
    exit 0
}

# Load configuration
$ConfigFile = ".swagger-to-http/hooks.config.ps1"

if (Test-Path $ConfigFile) {
    . .\$ConfigFile
} else {
    Write-Host "$WarningPrefix Configuration file not found at $ConfigFile" -ForegroundColor Yellow
    Write-Host "$WarningPrefix Using default settings" -ForegroundColor Yellow
    
    $HOOKS_ENABLED = $true
    $SWAGGER_FILE_PATTERNS = @("**/swagger.json", "**/swagger.yaml", "**/openapi.json", "**/openapi.yaml")
    $HTTP_OUTPUT_DIR = "http"
    $VALIDATE_SWAGGER = $true
    $SELECTIVE_UPDATES = $true
}

# Check if hooks are enabled in config
if ($HOOKS_ENABLED -ne $true) {
    Write-Host "$WarningPrefix Git hooks are disabled in $ConfigFile" -ForegroundColor Yellow
    exit 0
}

# Get the swagger-to-http binary path
$SWAGGER_TO_HTTP_BIN = "swagger-to-http"

# Check if swagger-to-http is in PATH
try {
    $null = Get-Command $SWAGGER_TO_HTTP_BIN -ErrorAction Stop
} catch {
    # Try to find it in the project
    if (Test-Path "./bin/swagger-to-http.exe") {
        $SWAGGER_TO_HTTP_BIN = "./bin/swagger-to-http.exe"
    } elseif (Test-Path "./swagger-to-http.exe") {
        $SWAGGER_TO_HTTP_BIN = "./swagger-to-http.exe"
    } else {
        Write-Host "$ErrorPrefix swagger-to-http binary not found" -ForegroundColor Red
        Write-Host "$ErrorPrefix Please make sure it's installed and in your PATH" -ForegroundColor Red
        exit 1
    }
}

# Check if this was a pull/merge that updated files
function Get-UpdatedSwaggerFiles {
    $files = git diff-tree -r --name-only --no-commit-id ORIG_HEAD HEAD
    $swaggerFiles = @()
    
    foreach ($file in $files) {
        if ($file -match "(swagger|openapi)\.(json|yaml|yml)$") {
            $swaggerFiles += $file
        }
    }
    
    return $swaggerFiles
}

$UPDATED_SWAGGER_FILES = Get-UpdatedSwaggerFiles

# If no Swagger files were updated, exit early
if ($UPDATED_SWAGGER_FILES.Count -eq 0) {
    Write-Host "$EchoPrefix No Swagger/OpenAPI files changed, skipping update" -ForegroundColor Cyan
    exit 0
}

Write-Host "$EchoPrefix Swagger/OpenAPI files changed in pull/merge:" -ForegroundColor Cyan
$UPDATED_SWAGGER_FILES | ForEach-Object { Write-Host "  $_" }

# Generate HTTP files for the updated Swagger files
Write-Host "$EchoPrefix Updating HTTP files..." -ForegroundColor Cyan

foreach ($FILE in $UPDATED_SWAGGER_FILES) {
    Write-Host "$EchoPrefix Generating HTTP files for $FILE..." -ForegroundColor Cyan
    
    # Check if file exists (it might have been deleted)
    if (-not (Test-Path $FILE)) {
        Write-Host "$WarningPrefix File $FILE no longer exists, skipping" -ForegroundColor Yellow
        continue
    }
    
    # Determine output directory
    $FILE_DIR = Split-Path -Path $FILE -Parent
    if (-not $FILE_DIR) { $FILE_DIR = "." }
    
    $OUTPUT_DIR = Join-Path -Path $HTTP_OUTPUT_DIR -ChildPath $FILE_DIR
    
    # Ensure output directory exists
    if (-not (Test-Path $OUTPUT_DIR)) {
        New-Item -ItemType Directory -Path $OUTPUT_DIR -Force | Out-Null
    }
    
    # Generate HTTP files
    if ($SELECTIVE_UPDATES -eq $true) {
        # Use selective update mode
        & $SWAGGER_TO_HTTP_BIN generate $FILE --output $OUTPUT_DIR --selective
    } else {
        # Regenerate all files
        & $SWAGGER_TO_HTTP_BIN generate $FILE --output $OUTPUT_DIR
    }
    
    # If the command failed, inform but don't abort
    if ($LASTEXITCODE -ne 0) {
        Write-Host "$ErrorPrefix Error generating HTTP files for $FILE" -ForegroundColor Red
        continue
    }
    
    Write-Host "$EchoPrefix HTTP files for $FILE updated" -ForegroundColor Cyan
}

# Notify the user that they may need to commit the changes
Write-Host "$EchoPrefix HTTP files have been updated based on pulled changes" -ForegroundColor Cyan
Write-Host "$EchoPrefix Please review the changes and commit them if needed" -ForegroundColor Cyan

exit 0

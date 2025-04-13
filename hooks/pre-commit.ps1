# PowerShell pre-commit hook for swagger-to-http
# Validates Swagger/OpenAPI files and updates HTTP files if needed

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

# Find staged Swagger files
function Get-StagedSwaggerFiles {
    $files = & git diff --cached --name-only --diff-filter=ACMR
    $swaggerFiles = @()
    
    foreach ($file in $files) {
        if ($file -match "(swagger|openapi)\.(json|yaml|yml)$") {
            $swaggerFiles += $file
        }
    }
    
    return $swaggerFiles
}

$STAGED_SWAGGER_FILES = Get-StagedSwaggerFiles

# If no Swagger files are staged, exit early
if ($STAGED_SWAGGER_FILES.Count -eq 0) {
    Write-Host "$EchoPrefix No Swagger/OpenAPI files staged, skipping validation" -ForegroundColor Cyan
    exit 0
}

Write-Host "$EchoPrefix Found staged Swagger/OpenAPI files:" -ForegroundColor Cyan
$STAGED_SWAGGER_FILES | ForEach-Object { Write-Host "  $_" }

# Validate Swagger files if enabled
if ($VALIDATE_SWAGGER -eq $true) {
    Write-Host "$EchoPrefix Validating Swagger/OpenAPI files..." -ForegroundColor Cyan
    
    $EXIT_CODE = 0
    
    # Validate each file
    foreach ($FILE in $STAGED_SWAGGER_FILES) {
        Write-Host "$EchoPrefix Validating $FILE..." -ForegroundColor Cyan
        
        # Check if file exists (it might have been deleted)
        if (-not (Test-Path $FILE)) {
            Write-Host "$WarningPrefix File $FILE no longer exists, skipping" -ForegroundColor Yellow
            continue
        }
        
        # Use swagger-to-http to validate the file
        try {
            & $SWAGGER_TO_HTTP_BIN validate $FILE 2>$null
            if ($LASTEXITCODE -ne 0) {
                Write-Host "$ErrorPrefix $FILE is not a valid Swagger/OpenAPI file" -ForegroundColor Red
                $EXIT_CODE = 1
            } else {
                Write-Host "$EchoPrefix $FILE is valid" -ForegroundColor Cyan
            }
        } catch {
            Write-Host "$ErrorPrefix Error validating $FILE: $_" -ForegroundColor Red
            $EXIT_CODE = 1
        }
    }
    
    # If any file is invalid, abort the commit
    if ($EXIT_CODE -ne 0) {
        Write-Host "$ErrorPrefix Aborting commit due to invalid Swagger/OpenAPI files" -ForegroundColor Red
        exit $EXIT_CODE
    }
}

# Generate HTTP files
Write-Host "$EchoPrefix Generating HTTP files..." -ForegroundColor Cyan

# Generate HTTP files for each Swagger file
foreach ($FILE in $STAGED_SWAGGER_FILES) {
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
    
    # Stage the generated files
    $generatedFiles = Get-ChildItem -Path "$OUTPUT_DIR/*.http" -File
    foreach ($generatedFile in $generatedFiles) {
        & git add $generatedFile.FullName
    }
    
    Write-Host "$EchoPrefix HTTP files for $FILE generated and staged" -ForegroundColor Cyan
}

Write-Host "$EchoPrefix Pre-commit hook completed successfully" -ForegroundColor Cyan
exit 0

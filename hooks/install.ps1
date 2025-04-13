# PowerShell script to install Git hooks for swagger-to-http

$ErrorActionPreference = "Stop"

$EchoPrefix = "[swagger-to-http hooks]"

Write-Host "$EchoPrefix Installing Git hooks..." -ForegroundColor Cyan

# Create .git/hooks directory if it doesn't exist
if (-not (Test-Path ".git/hooks")) {
    New-Item -ItemType Directory -Path ".git/hooks" -Force | Out-Null
}

# Install pre-commit hook
Copy-Item -Path "hooks/pre-commit.ps1" -Destination ".git/hooks/pre-commit.ps1" -Force

# Create pre-commit batch file to invoke PowerShell script
@"
@echo off
powershell.exe -ExecutionPolicy Bypass -NoProfile -File "%~dp0pre-commit.ps1"
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
"@ | Out-File -FilePath ".git/hooks/pre-commit" -Encoding ascii -Force

Write-Host "$EchoPrefix Installed pre-commit hook" -ForegroundColor Cyan

# Install post-merge hook
Copy-Item -Path "hooks/post-merge.ps1" -Destination ".git/hooks/post-merge.ps1" -Force

# Create post-merge batch file to invoke PowerShell script
@"
@echo off
powershell.exe -ExecutionPolicy Bypass -NoProfile -File "%~dp0post-merge.ps1"
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
"@ | Out-File -FilePath ".git/hooks/post-merge" -Encoding ascii -Force

Write-Host "$EchoPrefix Installed post-merge hook" -ForegroundColor Cyan

# Add husky configuration if it exists in the project
if ((Get-Command -Name npm -ErrorAction SilentlyContinue) -and (Test-Path "package.json")) {
    Write-Host "$EchoPrefix Detected Node.js project, configuring Husky..." -ForegroundColor Cyan
    
    # Check if husky is already installed
    $packageJson = Get-Content "package.json" -Raw | ConvertFrom-Json
    $huskyInstalled = $false
    
    if ($packageJson.devDependencies -and $packageJson.devDependencies.husky) {
        $huskyInstalled = $true
    }
    
    if (-not $huskyInstalled) {
        # Install husky
        npm install --save-dev husky
        npx husky install
    }
    
    # Add husky hooks
    npx husky set .husky/pre-commit "hooks/pre-commit.ps1"
    npx husky set .husky/post-merge "hooks/post-merge.ps1"
    
    Write-Host "$EchoPrefix Husky configuration complete" -ForegroundColor Cyan
}

# Create config directory for hook settings
if (-not (Test-Path ".swagger-to-http")) {
    New-Item -ItemType Directory -Path ".swagger-to-http" -Force | Out-Null
}

# Create default configuration if it doesn't exist
if (-not (Test-Path ".swagger-to-http/hooks.config.ps1")) {
    @"
# swagger-to-http Git hooks configuration for PowerShell

# Set to \$false to disable hooks temporarily
\$HOOKS_ENABLED = \$true

# Swagger/OpenAPI file patterns (array)
\$SWAGGER_FILE_PATTERNS = @("**/swagger.json", "**/swagger.yaml", "**/openapi.json", "**/openapi.yaml")

# Output directory for HTTP files
\$HTTP_OUTPUT_DIR = "http"

# Whether to validate Swagger/OpenAPI files before generating HTTP files
\$VALIDATE_SWAGGER = \$true

# Whether to regenerate all HTTP files on changes or only affected ones
\$SELECTIVE_UPDATES = \$true
"@ | Out-File -FilePath ".swagger-to-http/hooks.config.ps1" -Encoding utf8 -Force
    
    Write-Host "$EchoPrefix Created default configuration in .swagger-to-http/hooks.config.ps1" -ForegroundColor Cyan
}

Write-Host "$EchoPrefix Git hooks installation complete!" -ForegroundColor Cyan
Write-Host "$EchoPrefix You can customize hook behavior by editing .swagger-to-http/hooks.config.ps1" -ForegroundColor Cyan

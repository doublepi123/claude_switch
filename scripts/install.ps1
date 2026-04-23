param(
    [string]$InstallDir = "$HOME\AppData\Local\Programs\claude-switch\bin"
)

$ErrorActionPreference = 'Stop'

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Split-Path -Parent $scriptDir
$binaryPath = Join-Path $InstallDir 'cs.exe'

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go is required but was not found in PATH.'
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Write-Host "Building claude-switch..."
$env:GOCACHE = if ($env:GOCACHE) { $env:GOCACHE } else { Join-Path $repoRoot '.gocache' }
go build -o $binaryPath $repoRoot

Write-Host "Installed to: $binaryPath"

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (-not $userPath) {
    $userPath = ''
}

$pathEntries = $userPath -split ';' | Where-Object { $_ -ne '' }
if ($pathEntries -notcontains $InstallDir) {
    Write-Host "Add this directory to PATH if needed: $InstallDir"
}

$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host '==> building frontend' -ForegroundColor Cyan
Push-Location (Join-Path $root 'frontend')
if (-not (Test-Path 'node_modules')) { npm install --no-audit --no-fund }
npm run build
Pop-Location

Write-Host '==> syncing dist into backend' -ForegroundColor Cyan
$dst = Join-Path $root 'backend\web\dist'
if (Test-Path $dst) { Get-ChildItem $dst | ForEach-Object { Remove-Item -Recurse -Force $_.FullName } }
else { New-Item -ItemType Directory -Path $dst | Out-Null }
Copy-Item -Recurse -Force (Join-Path $root 'frontend\dist\*') $dst

Write-Host '==> building backend' -ForegroundColor Cyan
$env:Path = [System.Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [System.Environment]::GetEnvironmentVariable('Path','User')
Push-Location (Join-Path $root 'backend')
go build -ldflags="-s -w" -o (Join-Path $root 'port-monitor.exe') .
Pop-Location

Write-Host "==> done: $(Join-Path $root 'port-monitor.exe')" -ForegroundColor Green

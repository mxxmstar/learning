# 启动/停止 MySQL Docker 服务的 PowerShell 脚本

param(
    [switch]$stop = $false
)

Write-Host "MySQL Docker Service Controller" -ForegroundColor Cyan

if ($stop) {
    Write-Host "Stopping MySQL Docker service..." -ForegroundColor Yellow
    try {
        docker compose -f script/mysql/docker-compose.yaml down
        if ($LASTEXITCODE -eq 0) {
            Write-Host "MySQL Docker service stopped successfully!" -ForegroundColor Green
        } else {
            Write-Host "Failed to stop MySQL Docker service" -ForegroundColor Red
            exit $LASTEXITCODE
        }
    } catch {
        Write-Host "Error occurred while stopping MySQL Docker service: $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Starting MySQL Docker service..." -ForegroundColor Yellow
    try {
        docker compose -f script/mysql/docker-compose.yaml up -d
        if ($LASTEXITCODE -eq 0) {
            Write-Host "MySQL Docker service started successfully!" -ForegroundColor Green
            Write-Host "MySQL is running on port 3308" -ForegroundColor Green
        } else {
            Write-Host "Failed to start MySQL Docker service" -ForegroundColor Red
            exit $LASTEXITCODE
        }
    } catch {
        Write-Host "Error occurred while starting MySQL Docker service: $_" -ForegroundColor Red
        exit 1
    }
}
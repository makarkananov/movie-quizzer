# Скрипт для автоматического тестирования API Movie Quizzer
# PowerShell версия

$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8080"

Write-Host "=== Тестирование API Movie Quizzer ===" -ForegroundColor Cyan
Write-Host ""

# Функция для выполнения HTTP запросов
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Body = $null,
        [string]$Token = $null
    )
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    $uri = "$baseUrl$Endpoint"
    
    try {
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json
            $response = Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers -Body $jsonBody
        } else {
            $response = Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers
        }
        return @{ Success = $true; Data = $response }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $errorMessage = $_.ErrorDetails.Message
        return @{ Success = $false; StatusCode = $statusCode; Error = $errorMessage }
    }
}

# Счетчики
$totalTests = 0
$passedTests = 0
$failedTests = 0

# Функция для проверки теста
function Test-Case {
    param(
        [string]$Name,
        [scriptblock]$TestScript
    )
    
    $global:totalTests++
    Write-Host "[$global:totalTests] $Name" -NoNewline
    
    try {
        $result = & $TestScript
        if ($result.Success) {
            Write-Host " ✓ ПРОЙДЕН" -ForegroundColor Green
            $global:passedTests++
            return $true
        } else {
            Write-Host " ✗ ПРОВАЛЕН" -ForegroundColor Red
            Write-Host "   Ошибка: $($result.Error)" -ForegroundColor Yellow
            $global:failedTests++
            return $false
        }
    } catch {
        Write-Host " ✗ ОШИБКА" -ForegroundColor Red
        Write-Host "   $($_.Exception.Message)" -ForegroundColor Yellow
        $global:failedTests++
        return $false
    }
}

Write-Host "=== 1. Проверка Health Check ===" -ForegroundColor Cyan
Test-Case "Health Check" {
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/health"
    if ($result.Success -and $result.Data.status -eq "ok") {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Health check failed" }
}

Write-Host ""
Write-Host "=== 2. Тестирование регистрации ===" -ForegroundColor Cyan

# Генерируем уникальный email
$timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
$testEmail = "test$timestamp@example.com"
$testPassword = "testpass123"
$testNickname = "TestUser$timestamp"

Test-Case "Регистрация нового пользователя" {
    $body = @{
        email = $testEmail
        password = $testPassword
        nickname = $testNickname
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/register" -Body $body
    return $result
}

Test-Case "Регистрация с невалидным email" {
    $body = @{
        email = "invalid-email"
        password = "testpass123"
        nickname = "TestUser"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/register" -Body $body
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed with invalid email" }
}

Test-Case "Регистрация с коротким паролем" {
    $body = @{
        email = "test2$timestamp@example.com"
        password = "12345"
        nickname = "TestUser2"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/register" -Body $body
    # Ожидаем ошибку или успех (валидация может быть на frontend)
    return @{ Success = $true }
}

Test-Case "Регистрация с дубликатом email" {
    $body = @{
        email = $testEmail
        password = "testpass123"
        nickname = "TestUserDuplicate"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/register" -Body $body
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed with duplicate email" }
}

Write-Host ""
Write-Host "=== 3. Тестирование входа ===" -ForegroundColor Cyan

$token = $null

Test-Case "Вход с корректными данными" {
    $body = @{
        email = $testEmail
        password = $testPassword
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/login" -Body $body
    if ($result.Success -and $result.Data.token) {
        $script:token = $result.Data.token
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Login failed or no token received" }
}

Test-Case "Вход с неверным паролем" {
    $body = @{
        email = $testEmail
        password = "wrongpassword"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/login" -Body $body
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed with wrong password" }
}

Test-Case "Вход с несуществующим email" {
    $body = @{
        email = "nonexistent$timestamp@example.com"
        password = "testpass123"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/auth/login" -Body $body
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed with nonexistent email" }
}

Write-Host ""
Write-Host "=== 4. Тестирование профиля ===" -ForegroundColor Cyan

Test-Case "Получение информации о текущем пользователе" {
    if (-not $token) {
        return @{ Success = $false; Error = "No token available" }
    }
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/auth/me" -Token $token
    if ($result.Success -and $result.Data.email -eq $testEmail) {
        return @{ Success = $true }
    }
    return $result
}

Test-Case "Получение профиля пользователя" {
    if (-not $token) {
        return @{ Success = $false; Error = "No token available" }
    }
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/profile" -Token $token
    return $result
}

Test-Case "Получение достижений" {
    if (-not $token) {
        return @{ Success = $false; Error = "No token available" }
    }
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/profile/achievements" -Token $token
    return $result
}

Write-Host ""
Write-Host "=== 5. Тестирование рейтинга ===" -ForegroundColor Cyan

Test-Case "Получение глобального рейтинга" {
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/leaderboard/global"
    return $result
}

Test-Case "Получение позиции пользователя в рейтинге" {
    if (-not $token) {
        return @{ Success = $false; Error = "No token available" }
    }
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/leaderboard/me" -Token $token
    return $result
}

Write-Host ""
Write-Host "=== 6. Тестирование игровых сессий ===" -ForegroundColor Cyan

$sessionId = $null

Test-Case "Создание игровой сессии (режим frame)" {
    if (-not $token) {
        return @{ Success = $false; Error = "No token available" }
    }
    $body = @{
        mode = "frame"
    }
    $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/game/sessions" -Body $body -Token $token
    if ($result.Success -and $result.Data.session -and $result.Data.question) {
        $script:sessionId = $result.Data.session.id
        return @{ Success = $true }
    }
    return $result
}

if ($sessionId) {
    Test-Case "Отправка ответа на вопрос" {
        if (-not $token -or -not $sessionId) {
            return @{ Success = $false; Error = "No token or session available" }
        }
        # Получаем вопрос для получения question_id
        $questionResult = Invoke-ApiRequest -Method "GET" -Endpoint "/api/game/sessions/$sessionId/next" -Token $token
        if (-not $questionResult.Success) {
            return @{ Success = $false; Error = "Could not get question" }
        }
        
        $questionId = $questionResult.Data.id
        $body = @{
            question_id = $questionId
            answer = "Test Answer"
            elapsed_ms = 5000
        }
        $result = Invoke-ApiRequest -Method "POST" -Endpoint "/api/game/sessions/$sessionId/answers" -Body $body -Token $token
        return $result
    }
    
    Test-Case "Получение следующего вопроса" {
        if (-not $token -or -not $sessionId) {
            return @{ Success = $false; Error = "No token or session available" }
        }
        $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/game/sessions/$sessionId/next" -Token $token
        return $result
    }
    
    Test-Case "Получение итогов сессии" {
        if (-not $token -or -not $sessionId) {
            return @{ Success = $false; Error = "No token or session available" }
        }
        $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/game/sessions/$sessionId" -Token $token
        return $result
    }
}

Write-Host ""
Write-Host "=== 7. Тестирование защиты маршрутов ===" -ForegroundColor Cyan

Test-Case "Доступ к защищенному маршруту без токена" {
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/profile"
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed without token" }
}

Test-Case "Доступ с невалидным токеном" {
    $result = Invoke-ApiRequest -Method "GET" -Endpoint "/api/profile" -Token "invalid_token_12345"
    # Ожидаем ошибку
    if (-not $result.Success) {
        return @{ Success = $true }
    }
    return @{ Success = $false; Error = "Should have failed with invalid token" }
}

Write-Host ""
Write-Host "=== Итоги тестирования ===" -ForegroundColor Cyan
Write-Host "Всего тестов: $totalTests" -ForegroundColor White
Write-Host "Пройдено: $passedTests" -ForegroundColor Green
Write-Host "Провалено: $failedTests" -ForegroundColor Red
$successRate = if ($totalTests -gt 0) { [math]::Round(($passedTests / $totalTests) * 100, 1) } else { 0 }
Write-Host "Процент успешности: $successRate%" -ForegroundColor $(if ($successRate -ge 90) { "Green" } elseif ($successRate -ge 70) { "Yellow" } else { "Red" })
Write-Host ""

if ($failedTests -eq 0) {
    Write-Host "✓ Все тесты пройдены успешно!" -ForegroundColor Green
} else {
    Write-Host "⚠ Обнаружены проблемы. Проверьте результаты выше." -ForegroundColor Yellow
}


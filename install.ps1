param(
    [switch]$DryRun,
    [string]$BaseUrl = "",
    [string]$InstallDir = "",
    [string]$Versao = "latest",
    [string]$HomeDir = ""
)

$ErrorActionPreference = 'Stop'

function Resolver-Arquitetura {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq 'ARM64') {
        return 'arm64'
    }
    if ($arch -eq 'AMD64' -or [Environment]::Is64BitOperatingSystem) {
        return 'amd64'
    }
    throw "Arquitetura nao suportada: $arch"
}

function Validar-Url {
    param([string]$Url)
    if ($Url -match '^(https://|http://127\.0\.0\.1|http://localhost|http://\[::1\])') {
        return
    }
    throw "Download exige HTTPS (ou loopback para testes): $Url"
}

function Extrair-HashEsperado {
    param(
        [string]$SumsPath,
        [string]$ArquivoAlvo
    )
    foreach ($linha in (Get-Content -Path $SumsPath)) {
        $limpa = $linha.Trim()
        if ([string]::IsNullOrWhiteSpace($limpa)) {
            continue
        }
        if ($limpa -match '^([0-9a-fA-F]{64})\s+(.+)$') {
            if ($Matches[2].Trim() -eq $ArquivoAlvo) {
                return $Matches[1].ToLowerInvariant()
            }
        }
    }
    return ""
}

if ([string]::IsNullOrWhiteSpace($Versao)) {
    $Versao = "latest"
}

if ([string]::IsNullOrWhiteSpace($BaseUrl)) {
    if ($Versao -eq "latest") {
        $BaseUrl = "https://github.com/heverton-dev/LegoLoop/releases/latest/download"
    } else {
        $BaseUrl = "https://github.com/heverton-dev/LegoLoop/releases/download/$Versao"
    }
}
$BaseUrl = $BaseUrl.TrimEnd('/')
Validar-Url -Url $BaseUrl

if ([string]::IsNullOrWhiteSpace($HomeDir)) {
    $HomeDir = $HOME
}

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    if (-not [string]::IsNullOrWhiteSpace($env:LOCALAPPDATA)) {
        $InstallDir = Join-Path $env:LOCALAPPDATA "LegoLoop\bin"
    } else {
        $InstallDir = Join-Path $HomeDir ".local\bin"
    }
}

$arch = Resolver-Arquitetura
$arquivo = "legoloop_${Versao}_windows_${arch}.zip"
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("legoloop-install-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

try {
    $sumsPath = Join-Path $tmpDir "SHA256SUMS"
    $pkgPath = Join-Path $tmpDir $arquivo

    Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/SHA256SUMS" -OutFile $sumsPath
    Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/$arquivo" -OutFile $pkgPath

    $esperado = Extrair-HashEsperado -SumsPath $sumsPath -ArquivoAlvo $arquivo
    if ([string]::IsNullOrWhiteSpace($esperado)) {
        Write-Error "Checksum SHA-256 ausente para $arquivo em SHA256SUMS"
        exit 1
    }

    $atual = (Get-FileHash -Path $pkgPath -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($atual -ne $esperado) {
        Write-Error "SHA-256 divergente para $arquivo (esperado $esperado, obtido $atual)"
        exit 1
    }

    if ($DryRun) {
        Write-Output "DRY-RUN: SHA-256 conferido ($atual); instalaria legoloop.exe em $InstallDir e executaria legoloop.exe instalar"
        exit 0
    }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Expand-Archive -Path $pkgPath -DestinationPath $InstallDir -Force

    $binDest = Join-Path $InstallDir "legoloop.exe"
    if (-not (Test-Path -Path $binDest)) {
        Write-Error "Binario legoloop.exe nao encontrado apos extracao em $InstallDir"
        exit 1
    }

    & $binDest instalar --apply --home $HomeDir
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
} finally {
    if (Test-Path -Path $tmpDir) {
        Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

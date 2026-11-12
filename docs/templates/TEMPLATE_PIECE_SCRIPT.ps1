#Requires -Version 5.1
<#
.SYNOPSIS
    {{NOME_DO_SCRIPT}} - peça Script do LegoLoop (template V2.2.0).

.DESCRIPTION
    Transformação atômica, determinística e idempotente com validação de argumentos,
    confinamento do alvo à raiz do projeto, escrita via arquivo temporário, suporte a
    saída estruturada em JSON e códigos de saída padronizados.
    Taxonomia: Scripts de terceiros preservam seu nome semântico próprio; o prefixo
    legoloop-* é reservado às funções nativas da plataforma.
    O exemplo implementado normaliza um arquivo de texto: remove espaços no fim das
    linhas, converte para LF e garante uma única quebra de linha final.
    Substitua apenas a função ConvertTo-NormalizedText pela sua transformação real.

.PARAMETER TargetFile
    Arquivo a processar. Relativo à raiz do projeto ou absoluto dentro dela.

.PARAMETER ProjectRoot
    Raiz que confina o alvo. Padrão: diretório atual.

.PARAMETER DryRun
    Calcula e reporta o que mudaria, sem gravar nada no disco.

.PARAMETER Json
    Emite a telemetria factual e o resultado estruturado em JSON no stdout.

.EXAMPLE
    .\meu_script.ps1 -TargetFile "docs\specs\SPEC.md" -DryRun -Json

.NOTES
    os_support: windows
    execution_modes: cli, agente, estudio
    harness_support: any
    Invariantes: é proibido deixar TODO ou stub no corpo (Zero Stubs); qualquer falha termina
    com exit code diferente de zero e mensagem no stderr; o resultado nunca depende de estado oculto.
    Códigos de saída (exit code, mesma convenção da CLI legoloop):
      0  sucesso (inclui "nada a fazer" e dry-run)
      1  falha de execução (E/S, permissão, erro inesperado)
      2  uso ou entrada inválida (alvo ausente, fora da raiz ou não é arquivo)
    Salve este arquivo como UTF-8 com BOM: o Windows PowerShell 5.1 lê UTF-8 sem BOM como ANSI.
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory = $true, HelpMessage = "Caminho do arquivo alvo")]
    [ValidateNotNullOrEmpty()]
    [string]$TargetFile,

    [Parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]$ProjectRoot = (Get-Location).Path,

    [Parameter(Mandatory = $false)]
    [switch]$DryRun,

    [Parameter(Mandatory = $false)]
    [switch]$Json
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$Utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# Exceção própria para uso/entrada inválida: exceções do .NET (inclusive ArgumentException)
# nunca são confundidas com erro de uso e caem em "falha de execução" (exit 1).
class UsageError : System.Exception {
    UsageError([string]$Message) : base($Message) {}
}

function Write-Info([string]$Message) { 
    if (-not $Json) { Write-Host $Message } 
}

function Write-Err([string]$Message) { 
    [Console]::Error.WriteLine($Message) 
}

# Resolve o alvo e garante que ele fica dentro da raiz (bloqueia "..", absolutos externos).
function Resolve-InsideRoot {
    param([string]$Path, [string]$Root)
    $rootFull = [System.IO.Path]::GetFullPath($Root).TrimEnd('\', '/')
    $candidate = if ([System.IO.Path]::IsPathRooted($Path)) { $Path } else { Join-Path $rootFull $Path }
    $full = [System.IO.Path]::GetFullPath($candidate)
    $prefix = $rootFull + [System.IO.Path]::DirectorySeparatorChar
    if (-not $full.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase)) {
        throw [UsageError]::new("Alvo fora da raiz do projeto: $Path")
    }
    return $full
}

# Transformação real e pura: texto entra, texto sai. Idempotente por construção.
function ConvertTo-NormalizedText {
    param([string]$Text)
    $lines = @($Text -split "\r?\n" | ForEach-Object { $_.TrimEnd() })
    return (($lines -join "`n").TrimEnd("`n")) + "`n"
}

function Invoke-Main {
    $sw = [System.Diagnostics.Stopwatch]::StartNew()
    $path = Resolve-InsideRoot -Path $TargetFile -Root $ProjectRoot
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
        throw [UsageError]::new("Arquivo alvo não encontrado: $path")
    }

    $original = [System.IO.File]::ReadAllText($path)
    $result = ConvertTo-NormalizedText -Text $original
    $isChanged = ($result -cne $original)

    if (-not $isChanged) {
        Write-Info "Nada a fazer: o arquivo já está normalizado ($path)."
        if ($Json) {
            $telemetry = @{
                status = "noop"
                target = $path
                changed = $false
                duration_ms = $sw.ElapsedMilliseconds
                exit_code = 0
            }
            Write-Output ($telemetry | ConvertTo-Json -Compress)
        }
        return 0
    }

    if ($DryRun) {
        Write-Info "[DRY-RUN] O arquivo seria normalizado: $path. Nenhuma alteração gravada."
        if ($Json) {
            $telemetry = @{
                status = "dry_run"
                target = $path
                changed = $true
                duration_ms = $sw.ElapsedMilliseconds
                exit_code = 0
            }
            Write-Output ($telemetry | ConvertTo-Json -Compress)
        }
        return 0
    }

    if ((Get-Item -LiteralPath $path).IsReadOnly) {
        throw [System.UnauthorizedAccessException]::new("Arquivo somente leitura, nada foi gravado: $path")
    }

    $tmp = "$path.legoloop-tmp"
    try {
        [System.IO.File]::WriteAllText($tmp, $result, $Utf8NoBom)
        # Substituição atômica; falha se o destino não for gravável.
        [System.IO.File]::Replace($tmp, $path, [NullString]::Value)
    }
    finally {
        if (Test-Path -LiteralPath $tmp) { Remove-Item -LiteralPath $tmp -Force }
    }

    Write-Info "Sucesso: arquivo normalizado ($path)."
    if ($Json) {
        $telemetry = @{
            status = "applied"
            target = $path
            changed = $true
            duration_ms = $sw.ElapsedMilliseconds
            exit_code = 0
        }
        Write-Output ($telemetry | ConvertTo-Json -Compress)
    }
    return 0
}

try {
    exit ([int](Invoke-Main))
}
catch [UsageError] {
    Write-Err "ERRO DE USO: $($_.Exception.Message)"
    exit 2
}
catch {
    Write-Err "FALHA DE EXECUÇÃO: $($_.Exception.Message)"
    exit 1
}
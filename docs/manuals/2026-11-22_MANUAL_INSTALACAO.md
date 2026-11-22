# Manual de Instalação do LegoLoop

> **Data:** 2026-11-22  
> **Iniciativa:** I01 — Distribuição e Forja (`AC-E09.1`, `AC-E09.2`, `AC-E09.4`)  
> **Módulo Go:** `github.com/heverton-dev/LegoLoop`  

---

## 1. Visão Geral

Todos os canais de instalação colocam o binário `legoloop` disponível na máquina e executam o comando `legoloop instalar`, que aciona a Forja (`internal/forja`) com backup automático, verificação de integridade e geração dos relatórios `RELATORIO.json` e `RELATORIO.md`.

Os scripts oficiais baixam o pacote e o arquivo `SHA256SUMS` somente por HTTPS, conferem o hash SHA-256 antes de extrair qualquer arquivo e abortam imediatamente com código de saída diferente de zero se houver divergência.

---

## 2. Instalação via Script POSIX (`install.sh` — Linux e macOS)

### Simulação sem gravação (`--dry-run`)

```sh
curl -fsSL https://raw.githubusercontent.com/heverton-dev/LegoLoop/main/install.sh | sh -s -- --dry-run
```

### Instalação padrão

```sh
curl -fsSL https://raw.githubusercontent.com/heverton-dev/LegoLoop/main/install.sh | sh
```

### Parâmetros suportados

| Parâmetro | Descrição | Padrão |
| :--- | :--- | :--- |
| `--dry-run` | Baixa e valida o SHA-256 sem gravar binário nem aplicar a Forja | desativado |
| `--versao V` | Tag do release a instalar (ex.: `v1.0.0`) | `latest` |
| `--base-url URL` | URL base dos artefatos de release e `SHA256SUMS` | GitHub Releases oficial |
| `--install-dir DIR` | Diretório de destino do binário `legoloop` | `$HOME/.local/bin` |
| `--home DIR` | Diretório HOME repassado ao `legoloop instalar --apply --home` | `$HOME` |

Exemplo com versão fixa e diretório personalizado:

```sh
sh install.sh --versao v1.0.0 --install-dir "$HOME/.local/bin" --home "$HOME"
```

---

## 3. Instalação via PowerShell (`install.ps1` — Windows)

### Simulação sem gravação (`-DryRun`)

```powershell
irm https://raw.githubusercontent.com/heverton-dev/LegoLoop/main/install.ps1 -OutFile install.ps1
.\install.ps1 -DryRun
```

### Instalação padrão

```powershell
.\install.ps1
```

### Parâmetros suportados

| Parâmetro | Descrição | Padrão |
| :--- | :--- | :--- |
| `-DryRun` | Baixa e valida o SHA-256 via `Get-FileHash -Algorithm SHA256` sem gravar no disco | desativado |
| `-Versao V` | Tag do release a instalar (ex.: `v1.0.0`) | `latest` |
| `-BaseUrl URL` | URL base dos artefatos de release e `SHA256SUMS` | GitHub Releases oficial |
| `-InstallDir DIR` | Diretório de destino do binário `legoloop.exe` | `%LOCALAPPDATA%\LegoLoop\bin` |
| `-HomeDir DIR` | Diretório HOME repassado ao `legoloop.exe instalar --apply --home` | `$HOME` |

Exemplo com versão fixa e diretório personalizado:

```powershell
.\install.ps1 -Versao v1.0.0 -InstallDir "$env:LOCALAPPDATA\LegoLoop\bin" -HomeDir "$HOME"
```

---

## 4. Instalação via `git clone` e `go install ./cmd/legoloop`

Requer Go instalado (conforme a versão declarada no `go.mod`).

```sh
git clone https://github.com/heverton-dev/LegoLoop.git
cd LegoLoop
go install ./cmd/legoloop
legoloop instalar --dry-run
```

Para aplicar a Forja no ambiente local após conferir a simulação:

```sh
legoloop instalar --apply
```

---

## 5. Instalação Remota via `go install` (sem clone)

Após a publicação da tag no repositório remoto, instale diretamente pelo caminho canônico do módulo Go:

```sh
go install github.com/heverton-dev/LegoLoop/cmd/legoloop@<tag>
legoloop instalar --dry-run
```

Substitua `<tag>` pela versão publicada (por exemplo, `v1.0.0`) ou `latest` para a versão estável mais recente:

```sh
go install github.com/heverton-dev/LegoLoop/cmd/legoloop@latest
legoloop instalar --apply
```

---

## 6. Verificação Pós-Instalação

1. Confira a versão instalada:
   ```sh
   legoloop --version
   ```
2. Verifique o estado do ambiente e das ferramentas opcionais (como Python):
   ```sh
   legoloop doctor
   ```
3. Consulte o relatório gerado pela Forja em `~/.legoloop/forja/<AAAA-MM-DD>/RELATORIO.md`.
4. Caso precise reverter uma aplicação da Forja, use o backup automático:
   ```sh
   legoloop forjar --desfazer <AAAA-MM-DD>
   ```

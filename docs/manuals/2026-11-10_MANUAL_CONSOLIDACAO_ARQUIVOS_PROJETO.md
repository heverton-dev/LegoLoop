# 📘 MANUAL OPERACIONAL // CONSOLIDAÇÃO DE CÓDIGO FONTE EM ARQUIVO ÚNICO

> **Versão:** 1.0.0  
> **Data:** 2026-11-10  
> **Classificação:** Engenharia Operacional / Auditoria / Extração de Contexto  
> **Padrão:** ISO Governança LegoLoop & ASD-STE100  

---

## 1. VISÃO GERAL & PROPÓSITO

Este manual documenta e disseca os comandos de terminal utilizados no ecossistema **LegoLoop** para varrer, filtrar e concatenar todo o código-fonte de um repositório em um único arquivo de texto plano (`.txt`).

### Principais Casos de Uso:
1. **Auditoria Integral e Análise de Conformidade:** Inspecionar padrões arquiteturais sem dispersão em dezenas de arquivos.
2. **Injeção de Contexto em LLMs:** Enviar o repositório completo para modelos de grande contexto (Gemini Pro, Claude Sonnet/Opus) em uma única requisição.
3. **Backup Rápido e Snapshot de Estado:** Registrar a fotografia exata do código pré-refatoração sem sobrecarga de metadados do Git.

---

## 2. COMANDOS CANÔNICOS DE CONSOLIDAÇÃO

### 2.1. Comando 1 — Consolidação Integral Sem Truncamento (`AUDITORIA_FONTE.txt`)
Utilizado para gerar o arquivo completo com **100% das linhas** dos arquivos selecionados.

```powershell
Get-ChildItem -Recurse -File -Include *.go,*.md,*.json,*.sh,*.typst,go.mod | Where-Object { $_.FullName -notmatch '\\(\.git|\.codebase-memory|dist|node_modules)\\' } | ForEach-Object { "`n`n================ FILE: $($_.FullName) ================`n"; Get-Content $_.FullName -Raw } | Set-Content -Encoding UTF8 AUDITORIA_FONTE.txt
```

#### Anatomia do Pipeline:
| Etapa / Parâmetro | Função |
| :--- | :--- |
| `Get-ChildItem -Recurse -File` | Localiza recursivamente apenas arquivos reais (ignora pastas). |
| `-Include *.go,*.md,*.json,*.sh,*.typst,go.mod` | Filtra estritamente as extensões do ecossistema Go/Typst/Docs. |
| `Where-Object { $_.FullName -notmatch ... }` | Exclui pastas pesadas ou geradas (`.git`, `.codebase-memory`, `dist`, `node_modules`). |
| `ForEach-Object { ... }` | Itera cada arquivo inserindo um divisor visual `==== FILE: <caminho> ====` antes do conteúdo. |
| `Get-Content $_.FullName -Raw` | Lê o arquivo inteiro como uma string contínua mantendo todas as quebras de linha. |
| `Set-Content -Encoding UTF8 AUDITORIA_FONTE.txt` | Salva o resultado consolidado em codificação UTF-8 sem perdas de acentuação. |

---

### 2.2. Comando 2 — Amostragem com Teto de 500 Linhas (`dump_repositorio.txt`)
Utilizado quando o projeto é muito extenso e deseja-se limitar o número máximo de linhas por arquivo para economizar tokens ou evitar travamento de editores.

```powershell
Get-ChildItem -Recurse -File -Include *.py,*.md,*.yaml,*.yml,*.json,*.go,*.sh,*.ps1,*.toml | Where-Object { $_.FullName -notmatch '\\(\.git|node_modules|__pycache__|venv|\.venv)\\' -and $_.Name -ne 'dump_repositorio.txt' } | ForEach-Object { "=== ARQUIVO: $($_.FullName.Replace($PWD.Path, '')) ==="; Get-Content $_.FullName -TotalCount 500 -ErrorAction SilentlyContinue; "" } | Out-File -FilePath dump_repositorio.txt -Encoding utf8
```

#### Anatomia das Diferenças:
| Parâmetro | Comportamento |
| :--- | :--- |
| `Replace($PWD.Path, '')` | Torna o caminho relativo à raiz do projeto, economizando caracteres. |
| `-TotalCount 500` | Lê apenas até as primeiras 500 linhas de cada arquivo. |
| `-ErrorAction SilentlyContinue` | Ignora silenciosamente arquivos bloqueados ou inacessíveis. |
| `$_ -ne 'dump_repositorio.txt'` | Evita recursão infinita caso o arquivo de saída já exista na pasta. |

---

### 2.3. Comando 3 — Equivalente em Shell Linux / Bash / WSL
Para ambientes Linux, macOS ou contêineres Docker:

```bash
find . -maxdepth 5 -not -path '*/.*' -not -path './node_modules*' -not -path './dist*' -type f \( -name "*.go" -o -name "*.md" -o -name "*.json" -o -name "*.sh" -o -name "*.typst" \) -exec sh -c 'for f; do printf "\n\n================ FILE: %s ================\n\n" "$f"; cat "$f"; done' _ {} + > AUDITORIA_FONTE.txt
```

---

## 3. FERRAMENTAS DEDICADAS DO ECOSSISTEMA (CLI)

Se preferir utilizar ferramentas de empacotamento modernas de mercado que já geram contagem de tokens e estrutura em XML/Markdown:

### Repomix (antigo Repopack)
Instalação e execução instantânea via Node:
```powershell
npx repomix --style xml --output repomix-output.txt
```

### GitIngest (Web ou Python)
Gera uma representação focada no formato que prompts de LLM melhor interpretam:
```powershell
pip install gitingest
gitingest . -o gitingest_dump.txt
```

---

## 4. DIRETRIZES DE SEGURANÇA E GOVERNANÇA

> [!CAUTION]
> **Vazamento de Segredos e Chaves de API:**
> Arquivos `.env`, `credentials.json`, `.pem` ou tokens do `9router` **NUNCA** devem ser incluídos no filtro.
> Mantenha a cláusula de exclusão `Where-Object` sempre ativa para prevenir que segredos locais sejam despejados no `.txt`.

> [!IMPORTANT]
> **Adicione ao `.gitignore`:**
> Certifique-se de que arquivos de dump (`AUDITORIA_FONTE.txt`, `dump_repositorio.txt`, `repomix-output.txt`) constem no `.gitignore` para não inflar desnecessariamente o histórico Git do repositório.

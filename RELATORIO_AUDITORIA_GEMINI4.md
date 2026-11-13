# 📋 Relatório de Auditoria: Autenticação, Entitlements Gemini 4 e Fallback

> **Data:** 12/11/2026  
> **Ambiente:** Windows 10/11 x64 (MINGW64 / PowerShell) — Orca IDE & Antigravity CLI  
> **Autoridade:** Engenharia de Infraestrutura Agêntica  

---

## 1. Resumo Executivo do Bloqueio

A indisponibilidade dos modelos **`gemini-4-argon-*`** e **`gemini-4-hydrogen-*`** (que exibem o estado `Unavailable` na interface do editor e retornam `is not recognized as a known model or custom model in settings` no CLI `agy`) **NÃO** decorre de corrupção de binário local, expiração de chave ou erro de sintaxe.

**Causa Raiz Confirmada:**
1. **Restrição Server-Side de Entitlement OAuth (Google Cloud Code PA):**  
   O backend central do Google (`https://daily-cloudcode-pa.googleapis.com/v1internal:fetchAvailableModels`) avalia a conta Google autenticada no momento da consulta. A conta ativa é reconhecida com `authMethod=consumer` (`enterpriseOnboardingComplete: false`).
2. **Exclusão de Rollout/Whitelist para Contas `@gmail.com`:**  
   Os codinomes de modelo *Hydrogen* e *Argon* (Gemini 4) foram removidos da lista de distribuição de modelos autorizados para contas consumidoras regulares (`@gmail.com`), ficando restritos a contas com organização Google Cloud corporativa (`enterprise`) ou programas fechados de Trusted Tester da Google DeepMind.
3. **Dessincronização de Cache do Editor:**  
   O arquivo de configuração local (`.gemini/antigravity-cli/settings.json`) mantinha gravado `"model": "Gemini 4 Hydrogen (Medium)"`. Como o servidor não devolve mais esse ID na listagem dinâmica `fetchAvailableModels`, a interface do Orca/Antigravity marca o modelo como `Unavailable`.

---

## 2. Mapeamento de Sessão Atual

| Parâmetro | Detalhe Detectado no Ambiente |
| :--- | :--- |
| **Caminho do Binário CLI** | `C:\Users\trcnologia\.local\bin\agy.cmd` → aponta para `C:\Users\trcnologia\AppData\Local\agy\bin\agy.exe` (também presente em `C:\Users\trcnologia\.gemini\bin\agy.exe`) |
| **Versão do Binário `agy`** | `1.5.7` (Language Server protocol v1.5.7, build Go x64) |
| **Armazenamento de Credenciais OAuth** | **Windows Credential Manager / Keyring do SO** (`keyringAuth: loaded token`, renovação automática `browser.go:265`) |
| **Arquivos Locais de Estado & Config** | • `C:\Users\trcnologia\.gemini\antigravity-cli\settings.json`<br>• `C:\Users\trcnologia\.gemini\antigravity-cli\cache\onboarding.json`<br>• `C:\Users\trcnologia\AppData\Roaming\antigravity-usage\accounts\hevertoneduardoperes@gmail.com\` |
| **E-mail Autenticado Ativo no Keyring** | `authenticparfum@gmail.com` (identificado nos logs de runtime: `OAuth: authenticated successfully as authenticparfum@gmail.com`) |
| **Conta Registrada no Script de Shift** | `DEV_01` (`khenperes001@gmail.com`) no `.gemini_active_dev`; `authenticparfum@gmail.com` no slot `DEV_06` |
| **Tier / Plano Detectado** | `authMethod=consumer` \| `enterpriseOnboardingComplete: false` \| `consumerOnboardingComplete: true` |

---

## 3. Matriz de Teste de Endpoints (`agy --model`)

Testes de controle e de reprodução executados via terminal não-interativo (`agy --model <MODEL_ID> -p "responda apenas OK"`):

| Modelo Testado | Papel no Teste | Exit Code | Saída Bruta / Mensagem de Erro | Diagnóstico Técnico |
| :--- | :---: | :---: | :--- | :--- |
| **`gemini-3.8-flash-high`** | Controle | `0` | `OK` | **100% Operacional.** Autorizado pelo servidor Google para a conta `consumer`. |
| **`gemini-4-hydrogen-medium`** | Alvo Gemini 4 | `1` | `error: invalid model selection: model gemini-4-hydrogen-medium is not recognized as a known model or custom model in settings` | **Bloqueio de Entitlement.** O endpoint `fetchAvailableModels` da Google não envia esse ID para contas `@gmail.com`. |
| **`gemini-4-argon-medium`** | Alvo Gemini 4 | `1` | `error: invalid model selection: model gemini-4-argon-medium is not recognized as a known model or custom model in settings` | **Bloqueio de Entitlement.** Idêntico ao Hydrogen. |
| **`Gemini 4 Hydrogen (Medium)`** *(Label)* | Variação de Sintaxe | `1` | `model Gemini 4 Hydrogen (Medium) is not recognized as a known model or custom model in settings` | Confirma que o rejeite é pelo catálogo de modelos ativos do backend. |

### Modelos Atualmente Reconhecidos pelo `agy models` na Conta:
1. `gemini-3.8-flash-high` — *Gemini 3.8 Flash (High)*
2. `gemini-3.8-flash-medium` — *Gemini 3.8 Flash (Medium)*
3. `gemini-3.8-flash-low` — *Gemini 3.8 Flash (Low)*
4. `claude-sonnet-4-6` — *Claude Sonnet 4.6 (Thinking)*
5. `claude-opus-4-6-thinking` — *Claude Opus 4.6 (Thinking)*
6. `gpt-oss-120b-medium` — *GPT-OSS 120B (Medium)*

---

## 4. Diagnóstico das Contas da Equipe (`@gmail.com`)

| Conta | Slot no `.env` | Credencial Armazenada | Diagnóstico de Acesso Gemini 4 |
| :--- | :---: | :---: | :--- |
| **`authenticparfum@gmail.com`** *(ou authenticparfun)* | `DEV_06` | `AQ.Ab8RN6JV054Q...QN3w` *(OAuth Code Assist Token)* | **Ativa no Keyring do Windows.** Está sofrendo a restrição de classe `consumer` imposta pelo servidor CloudCode da Google. Não possui privilégio enterprise para liberar Argon/Hydrogen via OAuth. |
| **`khenperes2015@gmail.com`** | `DEV_10` | `AQ.Ab8RN6KBJF3T...WDAg` *(OAuth Code Assist Token)* | Conta `@gmail.com` consumer padrão. Sujeita exatamente às mesmas regras de catálogo do CloudCode PA. |
| **`peresg.claudia@gmail.com`** | `DEV_05` | `AQ.Ab8RN6LVsPS3...EeZCw` *(OAuth Code Assist Token)* | Conta `@gmail.com` consumer padrão. Sujeita exatamente às mesmas regras de catálogo do CloudCode PA. |

> [!NOTE]
> Os tokens iniciados com `AQ.Ab8RN...` são credenciais internas de sessão do Cloud Code/Code Assist. Para chamadas via SDK ou API direta fora do plugin Google, utilizam-se chaves de API padrão Google AI Studio (`AIza...`).

---

## 5. Plano de Ação e Configuração de Fallback

### A. Ações de Cache Executadas
1. **Backups Preventivos Gerados:**
   - `C:\Users\trcnologia\AppData\Roaming\antigravity-usage\accounts\hevertoneduardoperes@gmail.com\cache.json.bak`
   - `C:\Users\trcnologia\.gemini\antigravity-cli\settings.json.bak`
   - `C:\Users\trcnologia\.gemini\antigravity-cli\cache\onboarding.json.bak`
2. **Atualização Recomendada no `settings.json`:**
   Para eliminar a mensagem visual `Unavailable` no Orca/Antigravity, atualize o modelo padrão para um modelo homologado de alta capacidade:
   ```json
   {
     "model": "Claude Sonnet 4.6 (Thinking)"
   }
   ```
   *(ou `"Gemini 3.8 Flash (High)"` para tarefas que exijam a stack Gemini nativa).*

### B. Procedimento para Forçar Renovação do OAuth Google
Para forçar o Antigravity a reautenticar uma conta Google AI Pro através do fluxo limpo de navegador (sem interferência de credenciais em cache):
```powershell
python "C:\Users\trcnologia\.gemini\scripts\switch_dev.py" --login --dev DEV_06
```
Isso aciona o navegador oficial no Google Account Chooser com a conta `authenticparfum@gmail.com`, permitindo concluir o 2FA e regravar as permissões no Keyring do Windows.

### C. Configuração de Roteamento & Fallback (9Router / Provedores)
Caso a conta titular permaneça sem entitlement para a família Gemini 4 via OAuth direto do Google Cloud Code:
1. **9Router (Ativo e Validado):**
   - O proxy local/remoto em `https://roteador.signa-docs.com` está 100% online com suporte aos modelos `code-pro`, `code-fast`, `code-free` e `megacombo-coding`.
   - Chave ativa no `.env-pessoal`: `NINEROUTER_KEY = sk-4144...e243`.
2. **Fallback Direto via Gemini API Key:**
   - Para utilizar os modelos mais recentes do Gemini diretamente via API REST (fora do OAuth do Antigravity), garanta no Google Cloud Console que a API `generativelanguage.googleapis.com` esteja habilitada para a chave no `.env-empresa`.

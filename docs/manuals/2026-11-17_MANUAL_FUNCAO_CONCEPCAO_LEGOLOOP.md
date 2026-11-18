# Manual Operacional: `legoloop-concepcao` (v1.0.0)

## 1. Propósito e Visão Arquitetural

O processo `legoloop-concepcao` é o processo canônico soberano do ecossistema LegoLoop responsável pela etapa de **Concepção e Alinhamento de Intenção**.
Ele conduz o alinhamento da intenção do usuário por meio de entrevista socrática pelo motor oficial em Go (`legoloop fabricar entrevista`), garantindo que nenhuma linha de especificação, planejamento ou código seja redigida antes de atingir um limiar estrito de confiança ($\ge 95\%$) ou premonição com três perguntas e aprovação humana formal no Stop-Gate 01.

O motor do processo não altera o código Go e opera em contexto isolado e estéril (`orchestrix_hermetic`), produzindo rastreabilidade cronológica em `docs/intents/` e no cofre canônico de sessões.

## 2. Contratos de Entrada e Saída (I/O)

- **Entradas (Consumes):**
  - Ideia do usuário em linguagem natural.
  - Template canônico v2 da intenção: `docs/templates/TEMPLATE_INTENT.md`.
- **Saídas Mecânicas (Produces):**
  - `docs/intents/INTERVIEW.md`: Transcrição factual da sessão socrática (`intv_*`) arquivada no cofre de sessões.
  - `docs/intents/ESPELHO_CANONICO.md`: Síntese estruturada contendo os 6 campos inegociáveis (`resultado_esperado`, `quem_se_beneficia`, `por_que_agora`, `criterio_de_sucesso`, `restricoes_inegociaveis`, `fora_de_escopo`).
  - `docs/intents/INTENT.md`: Documento INTENT canônico materializado e promovido ao alias canônico, validado por `file_exists('docs/intents/INTENT.md')`.
  - `docs/intents/INTENT_APROVADA.md`: Registro formal da decisão de aprovação humana no Stop-Gate 01.

## 3. Topologia e Nós do DAG (`processos/legoloop-concepcao/pipeline.dag.json`)

O processo executa sequencialmente 4 nós auditados mecanicamente:

1. **`entrevista_socratica` (Tipo: `script`)**
   - Dispara `legoloop fabricar entrevista iniciar "<ideia>"`.
   - Gera a sessão socrática interativa e itera turno a turno.
   - Gate: `exit_code == 0`.
2. **`interview_gravada` (Tipo: `unit_no`)**
   - Consolida e persiste o Espelho Canônico e o registro de sessão em disco.
   - Gate: `file_exists('docs/intents/ESPELHO_CANONICO.md')`.
3. **`materializar_intent` (Tipo: `script`)**
   - Executa `legoloop fabricar entrevista intent <sessao>`.
   - Promove o artefato formatado para `docs/intents/INTENT.md`.
   - Gate: `file_exists('docs/intents/INTENT.md')`.
4. **`stop_gate_01` (Tipo: `stop_gate`)**
   - Pausa obrigatória para validação humana.
   - Registra a decisão formal em `docs/intents/INTENT_APROVADA.md`.
   - Gate: `file_exists('docs/intents/INTENT_APROVADA.md')`.

## 4. Invariantes de Governança e Regras Negativas

- **Uma pergunta por turno:** Proibido disparar baterias de perguntas simultâneas.
- **Marcadores obrigatórios:** Todo turno deve conter `HIPÓTESE`, `CONFIANÇA`, `PERGUNTA` e `PALPITE`.
- **Isenção de STE-PT na entrevista:** A condução socrática deve utilizar linguagem natural acolhedora e investigativa.
- **Stop-Gate 01 Humano intransponível:** Nenhuma etapa subsequente (`SPEC`, `PLAN` ou execução) pode iniciar sem aprovação humana formal.
- **Intocabilidade do Motor Go:** As implementações em `internal/interview` e `internal/cli/fabricar_entrevista.go` não recebem remendos em tempo de execução.

## 5. Comandos de Operação e Verificação

```bash
# Validar integridade e sintaxe do DAG
legoloop validate processos/legoloop-concepcao/pipeline.dag.json

# Iniciar entrevista interativa oficial
legoloop fabricar entrevista iniciar "Minha nova iniciativa"

# Responder ao turno ativo
legoloop fabricar entrevista responder <sessao_id> "Minha resposta"

# Materializar a INTENT a partir da sessão
legoloop fabricar entrevista intent <sessao_id>

# Executar o processo completo via CLI
legoloop exec legoloop-concepcao
```

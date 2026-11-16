# Auditoria de Arquitetura — Verbos de Processo LegoLoop

**Data:** 2026-11-16 · **Auditor:** audit-verbos · **Escopo:** fabricar, matriz, listar, refinar, remover, compor, fechar
**Premissa auditada:** os verbos deveriam criar/montar qualquer peça ou processo (agnosticidade universal).
**Conclusão:** a premissa falha. Os verbos montam por moldes fixos (vocabulário fechado embutido, I/O de artefato derivado do molde, cadeias rígidas). O molde decide antes da peça.

> [!IMPORTANT] Achado do coordenador validado: `fabricar` tem molde fixo de etapas. A reprodução runtime do dry-run (5 LACUNA) **não foi refeita nesta auditoria** (exige LLM real); a validação é mecânica por código, com prova em arquivo:linha abaixo. A cadeia causal que produz 5 LACUNA e gates genéricos está integralmente comprovada em código.

---

## 1. Matriz de Violações [verbo × violação × prova × gravidade]

| Verbo | Violação da agnosticidade | Prova (arquivo:linha) | Gravidade |
| :--- | :--- | :--- | :--- |
| fabricar | Molde fixo de etapas: vocabulário fechado embutido com 9 etapas canônicas; o cérebro só propõe termos da lista | `internal/matriz/vocab.default.json:3-13`; `internal/fabricar/natureza.go:24-25` | Crítica |
| fabricar | I/O de etapa derivado do molde, não da peça: artefato sempre `out/NN_<etapa>.md` (ou `.csv` se natureza=dados e etapa ∈ lista fixa) | `internal/fabricar/natureza.go:184` (`etapasDeDados`), `natureza.go:188-201` (`ArtefatoDe`) | Crítica |
| fabricar | Gate genérico por tipo de artefato, não por contrato da peça | `internal/fabricar/natureza.go:252-255` (`etapaDe` → `GateParaArtefato`); `internal/processo/gates.go:36-48` | Crítica |
| fabricar | Cadeia fixa de 5 camadas (Clarificação→Varredura→Fusão→Tríade→Template) para qualquer pedido | `internal/fabricar/engine.go:44` (`const Layers = 5`) | Alta |
| fabricar | Fallback fixo de peças (`aidd-spec`, `aidd-tickets`) quando nenhuma peça é apta | `internal/fabricar/merger.go:32-34` (`DefaultSteps`) | Alta |
| fabricar | Peça única aceita só 6 tipos fixos e ferramentas fixas por tipo | `internal/fabricar/peca_unica.go:13` (`TiposPeca`), `peca_unica.go:20-24` (`ferramentasDoWorker`) | Média |
| fabricar | Cascata de lacuna mascara a lacuna: inventa peça `curadoria-externa-<etapa>` e grava como se fosse peça | `internal/fabricar/cascata.go:33-35` (`sugerirFonteExterna`) | Alta |
| matriz | Aptidão medida só por etapas/casos do vocabulário fechado; perfil da peça limitado a 2 etapas, 2 naturezas, 3 casos | `internal/matriz/catalogador.go:65-67` (`maxCasos/maxEtapas/maxNaturezas`); `internal/matriz/rank.go:12-14` (`Needs`) | Alta |
| matriz | Tipos de peça fixos (skill\|mcp\|script\|hook\|gate\|loop); peça sem perfil nunca é escolhida | `internal/matriz/registro.go:18`; `internal/fabricar/lacunas.go:30` (`motivoDesconhecida`) | Alta |
| listar | Lista só processos do cofre; peça solta não é cidadã de primeira classe | `internal/cli/listar.go:17-43`; `internal/listar/listar.go:33-41` (`Item`) | Baixa |
| refinar | Ao adicionar etapa, I/O fixo: artefato sempre `out/<id>.md` markdown; não se muda artefato, gate nem natureza da etapa | `internal/refinar/adicionar.go:37`; `internal/cli/refinar.go:41-43` (só `--peca/--caso/--titulo`) | Alta |
| remover | Remove só processos `legoloop-*` do cofre; peças fora do molde não são gerenciáveis | `internal/cli/remover.go:20` (`prefixo`), `remover.go:26,104` | Média |
| compor | Composição aceita só processo `legoloop-<nome>@<versão>`; peça nunca compõe etapa por este verbo | `internal/compor/compor.go:53-57` | Média |
| fechar | Cadeia fixa de 9 etapas hardcoded, igual para qualquer natureza — inclui passo Go-cêntrico ("Recompilação do Binário") inócuo para mídia/dados | `internal/fechar/fechar.go:74-101` (`executarEtapas`) | Alta |
| (núcleo) | Modelo de etapa fixo: Consome lista de caminhos + Produz exatamente 1 artefato de 5 tipos fixos; não modela estado, efeitos ou I/O não-arquivo | `internal/processo/processo.go:32,48-49` | Crítica (causa raiz) |

---

## 2. Validação do achado: 5 LACUNA e gates genéricos

Cadeia causal comprovada em código:

1. O cérebro recebe vocabulário fechado e **não pode inventar etapas** (`natureza.go:24-25`); termos fora da lista são descartados (`natureza.go:122-146`).
2. Etapas aceitas viram passos; etapa sem peça apta vira lacuna (`internal/fabricar/lacunas.go:41-46`, `addLacuna`).
3. O dry-run reporta cada lacuna como `LACUNA | <motivo>` (`internal/cli/fabricar_report.go:86`).
4. Cada etapa ganha artefato do molde `out/NN_<etapa>.<ext>` (`natureza.go:188-201`) e gate derivado do **tipo do artefato**, não da peça (`natureza.go:252-255` → `processo/gates.go:36-48`).

Com 5+ etapas pedidas e cofre sem peças aptas, a saída tem 5 LACUNA e gates genéricos por construção. Não verificado em runtime nesta auditoria (exige LLM real); a lógica acima é determinística e coberta por testes (`internal/fabricar/calibragem_test.go:36-47,84-99`).

---

## 3. Causa raiz do molde fixo

**Não existe contrato de peça.** A peça não declara o que consome nem o que produz. Sem esse contrato:

- a **etapa** (não a peça) define o I/O — e a etapa nasce do vocabulário fechado embutido (`vocab.default.json`);
- o **gate** é deduzido da extensão do arquivo do molde, não da promessa da peça;
- a **aptidão** vira casamento de palavras-chave do vocabulário (`rank.go`), não compatibilidade de I/O;
- os verbos derivados (refinar, fechar, compor) herdam o molde: refinar não consegue mudar I/O; fechar tem cadeia própria fixa; compor só compõe processo.

O `Promete` (`processo.go:33`) é o único vestígio de contrato — booleano global, sem tipo de saída declarado pela peça.

---

## 4. Proposta: contrato de peça agnóstica (I/O por natureza)

### 4.1 Princípio

A peça declara seu I/O. O processo monta nós a partir do contrato declarado. O molde vira default, nunca restrição. Três naturezas de nó:

| Natureza do nó | O que é | I/O declarado por |
| :--- | :--- | :--- |
| **peça** | unidade executável (skill, script, hook, gate, mcp) | manifesto da própria peça |
| **subprocesso** | peça que embrulha um processo (equivalente ao `Compoe` atual) | contrato do processo embrulhado |
| **processo** | composição de nós com dependências | união dos contratos dos nós + artefatos intermediários |

### 4.2 Esboço do manifesto (`contrato.json` da peça)

```json
{
  "nome": "coletor-csv",
  "tipo": "script",
  "naturezas": ["dados"],
  "entra": [{ "tipo": "url", "opcional": false }],
  "sai": [{ "tipo": "dados", "caminho_padrao": "out/coleta.csv", "promete": true }],
  "efeitos": ["arquivo", "rede"],
  "gate_sugerido": "line_count('out/coleta.csv') >= 2"
}
```

### 4.3 Regras de derivação

1. **Etapa = nó do contrato.** `Produz`/`Consome` saem de `sai`/`entra` da peça; sem contrato, cai no molde atual (retrocompatível).
2. **Gate do contrato, não da extensão.** `gate_sugerido` da peça vence; `GateParaArtefato` vira fallback.
3. **Aptidão por compatibilidade.** A matriz casa `sai` de peça com `entra` do nó; vocabulário vira meta-dado de busca, não cadeado.
4. **Vocabulário aberto por acréscimo.** Etapas/naturezas novas entram por manifesto, sem recompilar (`vocab.default.json` vira semente, não moldura).
5. **Refinar/compor uniformes.** `refinar --sai/--entra` edita contratos; `compor` aceita nó de qualquer natureza (peça | subprocesso | processo).
6. **Fechar parametrizável.** As 9 etapas viram DAG por natureza: natureza Go mantém build; natureza mídia/dados troca por etapas pertinentes.

---

## 5. Resumo executivo

A agnosticidade falha nos 7 verbos. Gravidade crítica em fabricar e no núcleo (`processo.Etapa`): I/O e gate derivam de molde fixo, não de contrato de peça. Causa raiz única: ausência de contrato de I/O declarado na peça. A proposta corrige na raiz: manifesto de contrato (entra/sai/efeitos) com derivação de etapa, gate e aptidão a partir dele, e molde apenas como default retrocompatível.

**Não alterado:** nenhum código. **Não executado:** commit/push (proibidos pela missão).

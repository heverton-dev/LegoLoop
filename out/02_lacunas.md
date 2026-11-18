# 02 — Lista de Lacunas Técnicas do Alvo

> **Alvo:** cmd/legoloop e internal/cli (comando fabricar)
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Classificação:** Risco e Oportunidades de Rigor Arquitetural

---

## 1. Inventário de Lacunas Identificadas

A auditoria identificou três lacunas de projeto e governança no comando fabricar:

### Lacuna 01: Despacho Não Uniforme em cmd/legoloop/main.go
- **Evidência:** cmd/legoloop/main.go:16-22.
- **Descrição:** O main() possui tratamento especial hardcoded para triagem e migrar-processos, enquanto fabricar e os demais verbos passam pelo roteador genérico cli.Run().
- **Impacto:** Quebra a simetria de despacho de comandos da CLI. Casos de borda de contexto ou parsing de flags acabam duplicados ou inconsistentes entre comandos de primeira classe.
- **Severidade:** Baixa (organizacional).

### Lacuna 02: Validação Rígida de Categorias no CLI FlagSet
- **Evidência:** internal/cli/fabricar.go:55.
- **Descrição:** A flag --category possui valor padrão build e documentação restrita a build, reverse, web ou mobile. Contudo, o motor e o ecossistema evoluíram para contemplar novas naturezas como áudio, vídeo, livro e automação.
- **Impacto:** O usuário da CLI fica limitado a categorias pré-definidas ou recebe avisos ao tentar fabricar processos de naturezas emergentes.
- **Severidade:** Média (extensibilidade).

### Lacuna 03: Separação de I/O na Apresentação de Erros STE
- **Evidência:** internal/cli/fabricar.go:186-205 e fabricar_report.go:15-25.
- **Descrição:** Quando invocado com --json, alguns fluxos de erro ou avisos de lacuna podem misturar saída legível ao stderr e stdout, dependendo do estágio em que a falha ocorre (ex: aviso de cérebro ou timeout).
- **Impacto:** Automações externas que consomem a CLI como subprocesso e esperam JSON puro no stdout necessitam filtrar ou tratar mensagens residuais.
- **Severidade:** Baixa (robustez de automação).

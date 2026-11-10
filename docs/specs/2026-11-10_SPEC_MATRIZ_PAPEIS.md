# 📋 SPEC: Matriz Canônica de Papéis e Casamento Determinístico de Peças

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-10  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (Fusão Addy Osmani + Matt Pocock + AIDD-Spec)  
> **Iniciativa:** Fábrica Determinística LegoLoop  
> **Referência da Intenção:** [docs/intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md](../intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md)  
> **Alvo Primário:** internal/taxonomy, internal/matriz, internal/generator, cmd/legoloop  
> **Consome:** INTENT aprovada · **Produz:** PLAN e ROUTE  
> **Template:** TEMPLATE_SPEC.md v2.0.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)

> [!IMPORTANT]
> A IA lista abaixo todas as premissas técnicas adotadas para evitar decisões silenciosas. Alerte no Stop-Gate 02 caso alguma premissa esteja incorreta:
1. Runtime Go 1.23 conforme o go.mod, tipos nativos e zero dependências externas novas.
2. Armazenamento em arquivos JSON locais no SO (docs/vocab/, matriz.json, pipeline.dag.json), sem banco de dados externo.
3. SSOT tripartite: enums e structs em Go (internal/taxonomy), dicionário JSON canônico e schemas JSON.
4. Reuso integral dos pacotes internal/taxonomy e internal/matriz, expandindo estruturas existentes com total retrocompatibilidade via migração automática.

---

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)

O Processo e os agentes devem respeitar estritamente o sistema de fronteiras em 3 níveis:
* **Sempre Fazer:**
  * Desenvolver via TDD estrito (teste unitário falhando antes do código de produção).
  * Manter funções pequenas, enxutas e de responsabilidade única (< 35 linhas).
  * Validar rigorosamente todos os inputs e decodificações de JSON, incluindo caminhos (sem .., sem caminho absoluto fora da raiz).
  * Tratar caminhos de forma compatível com Windows e POSIX (ilepath.Join / ilepath.ToSlash).
  * Recompilar o binário ao final de qualquer ticket que altere Go (AGENTS.md, regra 5).
* **Perguntar Antes:**
  * Adicionar qualquer nova dependência externa no go.mod.
  * Alterar schemas persistidos no disco que possam quebrar retrocompatibilidade.
  * Alterar contratos públicos de CLI, de API do Studio, do DAG do Processo ou de ROUTE.json.
  * Modificar regras globais no AGENTS.md.
* **Nunca Fazer:**
  * Criar funções vazias com TODO, pass, // implementar depois ou mocks simulados em código de produção (Lei #1: Zero Stubs). Dublês de teste ficam em _test.go, apenas nas bordas (E/S, LLM, relógio).
  * Commitar arquivos temporários, segredos ou variáveis de ambiente no git.
  * Ignorar testes falhando ou quebrar a suíte global de testes.
  * Usar lert(), confirm() ou emojis em frontend (quando houver UI).
  * Permitir papéis arbitrários fora dos 18 crachás canônicos homologados.

---

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)

A solução implementa o catálogo de papéis canônicos no subsistema internal/taxonomy, atualiza os metadados das peças em internal/matriz, estende a CLI legoloop matriz ver e aprimora o algoritmo de seleção de ferramentas em internal/generator.

### Diagrama de Fronteiras (Mermaid)

`mermaid
flowchart TD
    subgraph Entrada [Camada de Entrada]
        In1[Nó do DAG com Papel]
        In2[Comando CLI matriz ver --papel]
    end

    subgraph Nucleo [Núcleo de Governança e Casamento]
        Tax[internal/taxonomy (18 Papéis / 7 Guildas)]
        Mat[internal/matriz (Peças com Papel Primário e Secundários)]
        Gen[internal/generator (Resolvedor de 2 Níveis)]
        Tax --> Mat
        Mat --> Gen
    end

    subgraph Saida [Camada de Saída]
        Out1[Ferramenta Ouro Casada sem Alucinação]
        Out2[Tabela Filtrada no Terminal / Erro Fatal se Inválido]
    end

    In1 --> Gen
    In2 --> Mat
    Gen --> Out1
    Mat --> Out2
`

### Costuras de Teste (Test Seams — Pocock)
* **Costura Primária (Highest Seam):** 	axonomy.ValidarPapel(slug string) error e matriz.FiltrarPorPapel(papel string, guilda string) ([]Peca, error) e generator.CasarFerramenta(papelRequerido string) (Peca, error).
* **Justificativa:** Funções puras e de entrada de módulo que operam sobre estruturas de dados em memória e arquivos de fixture sem necessidade de simulações externas complexas.
* **Costuras Secundárias:** 	axonomy.ListarGuildas() []Guilda e matriz.MigrarLegado(p *Peca) error.
* **Dublês Permitidos:** Sistema de arquivos temporário via 	.TempDir(). Zero mocks de lógica interna.

---

## 4. Contratos de Dados & Schemas Formais

`go
package taxonomy

import (
	errors
	strings
)

// Erros sentinela do subsistema de papéis e taxonomia.
var (
	ErrPapelInvalido             = errors.New(taxonomy: papel operacional desconhecido ou fora do catalogo)
	ErrGuildaInvalida            = errors.New(taxonomy: guilda operacional desconhecida)
	ErrPapelPrimarioObrigatorio  = errors.New(matriz: peca deve possuir papel_primario valido)
	ErrNenhumaFerramentaEncontrada = errors.New(generator: nenhuma ferramenta compativel encontrada no cofre)
)

// Guilda representa uma divisao operacional da fabrica.
type Guilda string

const (
	GuildaGovernanca Guilda = governanca
	GuildaReversa    Guilda = reversa
	GuildaConstrucao Guilda = construcao
	GuildaInfra      Guilda = infra
	GuildaAudio      Guilda = audiovisual
	GuildaQualidade  Guilda = qualidade
	GuildaPublicacao Guilda = publicacao
)

// Papel representa um cracha canonico de execucao.
type Papel struct {
	Slug        string json:slug
	Nome        string json:nome
	Guilda      Guilda json:guilda
	Descricao   string json:descricao
	Ferramentas string json:ferramentas_tipicas
}

// CatalogoCanonico contem os 18 crachas oficiais das 7 guildas.
var CatalogoCanonico = map[string]Papel{
	// 1. Governanca
	patrocinador:      {Slug: patrocinador, Nome: Patrocinador, Guilda: GuildaGovernanca, Descricao: Aprovacao de regras e Stop-Gates, Ferramentas: Chat UI},
	arquiteto:         {Slug: arquiteto, Nome: Arquiteto, Guilda: GuildaGovernanca, Descricao: Requisitos EARS e schemas, Ferramentas: Templates v2},
	planejador:        {Slug: planejador, Nome: Planejador, Guilda: GuildaGovernanca, Descricao: Divisao em tickets atomicos, Ferramentas: PLAN v2},
	orquestrador:      {Slug: orquestrador, Nome: Orquestrador, Guilda: GuildaGovernanca, Descricao: Despacho do DAG e isolamento esteril, Ferramentas: LegoLoop CLI},
	// 2. Engenharia Reversa
	dissecador-web:    {Slug: dissecador-web, Nome: Dissecador Web, Guilda: GuildaReversa, Descricao: Extracao de tokens e rotas SaaS, Ferramentas: REA Web},
	dissecador-mobile: {Slug: dissecador-mobile, Nome: Dissecador Mobile, Guilda: GuildaReversa, Descricao: Descompilacao de APKs, Ferramentas: Jeb / Apktool},
	dissecador-midia:  {Slug: dissecador-midia, Nome: Dissecador Midia, Guilda: GuildaReversa, Descricao: Forense de video e transcricao, Ferramentas: Video-REA},
	// 3. Construcao
	engenheiro-frontend: {Slug: engenheiro-frontend, Nome: Engenheiro Frontend, Guilda: GuildaConstrucao, Descricao: Interfaces Web e Next.js, Ferramentas: React / Tailwind},
	engenheiro-mobile:   {Slug: engenheiro-mobile, Nome: Engenheiro Mobile, Guilda: GuildaConstrucao, Descricao: Aplicativos moveis, Ferramentas: Android / iOS},
	engenheiro-backend:  {Slug: engenheiro-backend, Nome: Engenheiro Backend, Guilda: GuildaConstrucao, Descricao: APIs e logica de dominio, Ferramentas: Go / REST},
	engenheiro-tdd:      {Slug: engenheiro-tdd, Nome: Engenheiro TDD, Guilda: GuildaConstrucao, Descricao: Red-Green-Refactor e testes, Ferramentas: Go Test / Vitest},
	// 4. Infra e Operacoes
	engenheiro-cloud:  {Slug: engenheiro-cloud, Nome: Engenheiro Cloud, Guilda: GuildaInfra, Descricao: VPS, Edge e Docker, Ferramentas: Cloudflare / Ops},
	engenheiro-secops: {Slug: engenheiro-secops, Nome: Engenheiro SecOps, Guilda: GuildaInfra, Descricao: Path-jail e segredos, Ferramentas: Security Gates},
	// 5. Audiovisual
	roteirista-copy:   {Slug: roteirista-copy, Nome: Roteirista Copy, Guilda: GuildaAudio, Descricao: Roteirizacao estruturada, Ferramentas: Craft Story},
	ingestor-midia:    {Slug: ingestor-midia, Nome: Ingestor Midia, Guilda: GuildaAudio, Descricao: Download e FFMPEG, Ferramentas: yt-dlp / Whisper},
	animador-motion:   {Slug: animador-motion, Nome: Animador Motion, Guilda: GuildaAudio, Descricao: Motion procedural em codigo, Ferramentas: Remotion / Three.js},
	pos-produtor:      {Slug: pos-produtor, Nome: Pos Produtor, Guilda: GuildaAudio, Descricao: Cortes de silencio e RMS, Ferramentas: Auto-cutter},
	// 6. Qualidade
	auditor-qualidade: {Slug: auditor-qualidade, Nome: Auditor Qualidade, Guilda: GuildaQualidade, Descricao: Linters e mutacao de teste, Ferramentas: Taxlinter / Funclen},
	// 7. Publicacao
	redator-tecnico:   {Slug: redator-tecnico, Nome: Redator Tecnico, Guilda: GuildaPublicacao, Descricao: Manuais ASD-STE100, Ferramentas: Manual Creator},
	compilador-livro:  {Slug: compilador-livro, Nome: Compilador Livro, Guilda: GuildaPublicacao, Descricao: Livros didaticos e PDF, Ferramentas: LegoPress},
}

// ValidarPapel verifica se o slug pertence a SSOT oficial.
func ValidarPapel(slug string) error {
	s := strings.TrimSpace(slug)
	if _, ok := CatalogoCanonico[s]; !ok {
		return ErrPapelInvalido
	}
	return nil
}
`

`go
package matriz

import taxonomy

// Peca estendida com metadados de papeis.
type Peca struct {
	ID                string   json:id
	Nome              string   json:nome
	Natureza          string   json:natureza
	Tier              string   json:tier
	Score             float64  json:score
	PapelPrimario     string   json:papel_primario
	PapeisSecundarios []string json:papeis_secundarios,omitempty
}

// Validate aplica o disjuntor estrito Poka-Yoke.
func (p Peca) Validate() error {
	if err := taxonomy.ValidarPapel(p.PapelPrimario); err != nil {
		return ErrPapelPrimarioObrigatorio
	}
	for _, sec := range p.PapeisSecundarios {
		if err := taxonomy.ValidarPapel(sec); err != nil {
			return taxonomy.ErrPapelInvalido
		}
	}
	return nil
}
`

---

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — Catálogo de Papéis e Validação Estrita (Taxonomia)

* [ ] AC-E01.1: **Validação com Sucesso de Papel Homologado**
  * **Dado que:** O identificador engenheiro-frontend é submetido à validação
  * **Quando:** A função 	axonomy.ValidarPapel(engenheiro-frontend) é invocada
  * **Então:** Retorna 
il e confirma pertencimento à guilda construcao.
* [ ] AC-E01.2: **Disjuntor Poka-Yoke para Papel Desconhecido**
  * **Dado que:** O identificador ninja-react é submetido à validação
  * **Quando:** A função 	axonomy.ValidarPapel(ninja-react) é invocada
  * **Então:** Retorna erro sentinela 	axonomy.ErrPapelInvalido e falha mecanicamente.

### Épico 02 — Contrato da Peça na Matriz e Filtros na CLI

* [ ] AC-E02.1: **Rejeição de Peça com Papel Primário Inválido**
  * **Dado que:** Uma peça com PapelPrimario = termo-invalido é validada
  * **Quando:** O método peca.Validate() é executado
  * **Então:** Retorna erro sentinela e impede inserção na Matriz.
* [ ] AC-E02.2: **Filtro por Papel na CLI**
  * **Dado que:** A Matriz possui 20 peças com diferentes papéis
  * **Quando:** O comando legoloop matriz ver --papel engenheiro-tdd é executado
  * **Então:** Retorna apenas peças cujo papel primário ou secundário seja engenheiro-tdd.
* [ ] AC-E02.3: **Filtro por Guilda na CLI**
  * **Dado que:** O usuário deseja listar ferramentas da guilda audiovisual
  * **Quando:** O comando legoloop matriz ver --guilda audiovisual é executado
  * **Então:** Retorna todas as peças associadas aos 4 crachás da guilda audiovisual.

### Épico 03 — Casamento Determinístico no Gerador (Fabricar)

* [ ] AC-E03.1: **Priorização de Papel Primário sobre Secundário**
  * **Dado que:** Existem duas peças: Peça A (primário engenheiro-tdd, score 0.90) e Peça B (secundário engenheiro-tdd, score 0.95)
  * **Quando:** O gerador resolve uma ferramenta para o crachá engenheiro-tdd
  * **Então:** Seleciona Peça A por compatibilidade primária determinística.
* [ ] AC-E03.2: **Alerta Explícito de Lacuna de Ferramenta**
  * **Dado que:** Um nó do processo exige o papel dissecador-mobile e nenhuma peça no cofre atende a este crachá
  * **Quando:** O gerador resolve a ferramenta
  * **Então:** Retorna ErrNenhumaFerramentaEncontrada e interrompe sem inventar ferramenta simulada.

---

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)

### Comandos de Verificação Mecânica
* **Build:** go build -o bin/legoloop.exe ./cmd/legoloop
* **Test:** go test -count=1 ./...
* **Estática e formatação:** go vet ./... e gofmt -l . (saída vazia).
* **Taxonomia e Vocabulário:** go run ./tools/taxlinter internal cmd loops skills docs/templates
* **Dependências intactas:** git diff --exit-code go.mod go.sum

### Critério de Done Global
A especificação técnica é considerada **100% CONCLUÍDA** quando:
- [ ] 100% dos critérios binários AC-E{N}.{M} estiverem marcados como atendidos ([x]).
- [ ] O comando de **Test** passar com exit code 0.
- [ ] go vet ./... e gofmt -l . não reportarem nada.
- [ ] O linter taxonômico passar com exit code 0.
- [ ] O comando de **Build** recompilar com sucesso o binário in/legoloop.exe.
- [ ] Zero arquivos temporários ou órfãos residuais.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**  
*As premissas assumidas, as fronteiras (boundaries), os contratos e os critérios binários acima atendem à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar o PLAN (TEMPLATE_PLAN.md) e a ROUTE (TEMPLATE_ROUTE.json).*

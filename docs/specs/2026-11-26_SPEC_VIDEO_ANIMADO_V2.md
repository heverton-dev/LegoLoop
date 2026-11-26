---
id: SPEC-LEGOLOOP-VIDEO-ANIMADO-002
status: APPROVED
date: 2026-11-26
schema_version: 2
upstream: [EMENDA-INTENT-VIDEO-ANIMADO-20261126]
substitui_em_parte: [SPEC-LEGOLOOP-DIORAMA-RENDER-001]
covers: [R7, R8, R9, R10, R11, R12]
---

# 📋 SPEC v2.0: macro `legoloop-video-animado` (inventário por estilo, contratos novos, motor portado, renomeação)

> **Status:** APPROVED (Stop-Gate 02 aprovado em 2026-11-26, resposta "[1]").
> **Data:** 2026-11-26. **Versão:** 2.0.0.
> **Referência da Intenção:** [docs/intents/2026-11-26_EMENDA_INTENT_VIDEO_ANIMADO.md](../intents/2026-11-26_EMENDA_INTENT_VIDEO_ANIMADO.md) (aprovada no Stop-Gate 01)  
> **Relação com a SPEC v1.4.0** (`docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md`): os épicos E01 a E06 continuam válidos para voz, Whisper, mix, SFX, render e cortes. Mudam apenas: E01 (pré-voo ganha estilo e motor), E03 (a timeline passa a vir do motor portado, via ponte do E10) e E06 (nomes). Os 29 critérios `AC-E01.1` a `AC-E06.2` ficam, exceto os do E03 que o E10 substitui.
> **Base de fatos:** `docs/progress/2026-11-25_plan-legoloop-diorama-render/ANALISE_IMPACTO_INVENTARIO_ESTILOS.md` e `achados.md` (A7 a A12).

---

## 1. Premissas Técnicas

1. Decisões D1 a D6 da Emenda valem como premissas (nome `video-animado`, inventário por estilo, `modo_producao` e `estilo`, motor portado, palito primeiro, licença pela opção A).
2. Go 1.21+ só com biblioteca padrão; `go.mod` e `go.sum` intactos.
3. O motor roda em Python + Playwright + Chromium + FFmpeg (já comprovado nesta máquina em 2026-11-26: 4 cenas do First Brands em 31,7 s e a cena 1 do demo em 15 s, 0 crédito).
4. Enquanto a licença da RobbiSkill não for liberada, o motor fica em `motor-local/` (ignorado pelo git). Go o localiza pela variável `LEGOLOOP_MOTOR_DIR` ou por `motor-local/`.
5. O inventário guarda o índice em JSON versionado (`inventario/video/`) e os arquivos pesados em `projetos/inventario/` (fora do git), ligados por hash SHA-256.
6. Cada item de inventário pertence a um só estilo. Só SFX, trilha e voz ficam em `comum/`.
7. Flow continua só na aba do usuário (porta 9222), imagem a 0 crédito, gasto só com cotação e "sim".
8. O processo não depende de agente nem de sistema operacional: caminhos com `filepath`, comandos externos injetáveis nos testes.

## 2. Fronteiras de Ação

### Sempre fazer
* TDD Red-Green com testes herméticos (motor, Flow e Whisper injetáveis).
* Gravar JSON de forma atômica (arquivo temporário e renomeação).
* Registrar cada achado em `achados.md` e cada decisão em arquivo.

### Perguntar antes
* Qualquer gasto de crédito do Flow.
* Aprovar um estilo (`rascunho` para `validado`) e aprovar itens novos do inventário.
* Tirar o motor de `motor-local/` para `internal/`.

### Nunca fazer
* Publicar o repositório com o motor dentro enquanto a licença estiver pendente.
* Aprovar item em estilo `rascunho`.
* Apagar versões de render ou itens do inventário (só marcar como obsoleto).
* Reescrever documentos históricos (só nota de redirecionamento no topo).

## 3. Arquitetura e Costuras de Teste

```
inventario/video/                      (índice versionado)
  comum/{sfx,trilhas,vozes}.json
  estilos/<estilo>/
    estilo.json                        (referências, renderizador, prompt mestre, paleta, luz, status, versão, herda)
    personagens.json  skins.json  trajes.json  penteados.json  adornos.json
    aderecos.json  poses.json  humores.json  cenarios.json  objetos.json
projetos/inventario/<estilo>/...       (arquivos pesados, fora do git)
motor-local/                           (motor portado, fora do git enquanto a licença estiver pendente)
internal/video/{roteiro,ativos,render,inventario,distribuicao}/
```

Costuras de teste: `Motor` (interface que executa o render e o `check_timeline`), `Flow` (interface de geração e download), `Relogio`, `SistemaDeArquivos`. Nenhum teste chama Chromium, Flow ou rede.

## 4. Contratos de Dados

```go
package inventario

var (
	ErrEstiloDuplicado    = errors.New("video-inventario: estilo ja existe")
	ErrEstiloInexistente  = errors.New("video-inventario: estilo inexistente")
	ErrEstiloNaoValidado  = errors.New("video-inventario: estilo nao validado nao aceita item aprovado")
	ErrTipoInvalido       = errors.New("video-inventario: tipo de item invalido")
	ErrHashDivergente     = errors.New("video-inventario: hash do arquivo diverge")
	ErrHerancaCiclica     = errors.New("video-inventario: heranca de estilos em ciclo")
	ErrParteNaoSuportada  = errors.New("video-inventario: parte nao suportada pelo renderizador")
)

type Estilo struct {
	SchemaVersion   int          `json:"schema_version"`
	ID              string       `json:"id"`               // kebab-case
	Nome            string       `json:"nome"`
	Status          string       `json:"status"`           // "rascunho" | "validado"
	Versao          string       `json:"versao"`           // semver
	Herda           string       `json:"herda,omitempty"`  // id de outro estilo
	Renderizador    string       `json:"renderizador"`     // ex.: "molde-pilula"
	PartesSuportadas []string    `json:"partes_suportadas"` // penteados, adornos, trajes que o renderizador desenha
	PromptMestre    string       `json:"prompt_mestre"`    // prompt de estilo dos cenários do Flow
	Paleta          []string     `json:"paleta"`
	Luz             *Luz         `json:"luz,omitempty"`    // rim, cor, quantidade
	Referencias     []Referencia `json:"referencias"`      // imagem ou vídeo, com origem
}

type Referencia struct {
	Tipo   string `json:"tipo"`   // "imagem" | "video" | "link"
	Origem string `json:"origem"` // caminho ou URL
	Nota   string `json:"nota,omitempty"`
}

type Item struct {
	SchemaVersion int               `json:"schema_version"`
	ID            string            `json:"id"`
	Tipo          string            `json:"tipo"`           // pele|traje|penteado|adorno|aderecos|pose|humor|personagem|cenario|objeto
	Estilo        string            `json:"estilo"`
	EstiloVersao  string            `json:"estilo_versao"`
	Tags          []string          `json:"tags"`           // tema, epoca, canal
	Parametros    map[string]any    `json:"parametros"`     // cores, forma, bloco cast
	Arquivo       string            `json:"arquivo,omitempty"`
	Hash          string            `json:"hash,omitempty"` // SHA-256 do arquivo
	Origem        string            `json:"origem"`         // projeto e data
	Aprovado      bool              `json:"aprovado"`
	Herdavel      bool              `json:"herdavel,omitempty"`
	Obsoleto      bool              `json:"obsoleto,omitempty"`
}
```

`projeto.json` v2 acrescenta `modo_producao` (`codigo|flow|hibrido`) e `estilo` (id). O campo antigo `estilo_visual` continua aceito por uma versão, com aviso, e é lido como `modo_producao`.
`cenas.json` v2: o elemento móvel de personagem ganha `personagem_id`, `pose`, `humor`, `olhar` e `fala` (booleano); `sprite_id` fica só para renderizadores de sprite.
`elenco.json`: lista de personagens, cada um com um bloco `cast` (`skin`, `suit`, `shirt`, `tie`, `hair`, `hair_col`, `glasses`, `hat`, etc.) formado só por itens do inventário do estilo.

## 5. Critérios Binários de Aceite por Épico

### Épico 07: Inventário e Estilos (R7)

* [ ] `AC-E07.1`: **Criar estilo**
  * **Dado que:** `estilo.json` válido (`id` kebab-case, `renderizador`, `status: rascunho`, `prompt_mestre`, paleta e 1 ou mais referências com origem).
  * **Quando:** `NovoEstilo` é chamado.
  * **Então:** grava `inventario/video/estilos/<id>/estilo.json` de forma atômica; um `id` repetido retorna `ErrEstiloDuplicado` sem alterar o arquivo.

* [ ] `AC-E07.2`: **Item pertence a um estilo existente**
  * **Dado que:** item com `estilo` que não existe, ou `tipo` fora da lista.
  * **Quando:** `GravarItem` é chamado.
  * **Então:** retorna `ErrEstiloInexistente` ou `ErrTipoInvalido` e nada é gravado.

* [ ] `AC-E07.3`: **Estilo em rascunho não aceita item aprovado**
  * **Dado que:** estilo com `status: rascunho`.
  * **Quando:** `GravarItem` recebe `aprovado: true`.
  * **Então:** retorna `ErrEstiloNaoValidado`; com `aprovado: false` grava normalmente.

* [ ] `AC-E07.4`: **Busca por estilo, tipo e tags**
  * **Dado que:** 5 itens de 2 estilos e tipos variados.
  * **Quando:** `Buscar(estilo, tipo, tags)` é chamado.
  * **Então:** devolve só os itens aprovados e não obsoletos que casam com todas as tags, ordenados por `id`; sem resultado devolve lista vazia e `nil`.

* [ ] `AC-E07.5`: **Arquivo pesado confere por hash**
  * **Dado que:** item com `arquivo` e `hash`.
  * **Quando:** `Verificar` roda com o arquivo ausente ou alterado.
  * **Então:** retorna `ErrHashDivergente` apontando o item.

* [ ] `AC-E07.6`: **Área comum só aceita itens sem estilo**
  * **Dado que:** itens `sfx`, `trilha` e `voz`.
  * **Quando:** gravados em `comum/`.
  * **Então:** são aceitos; um tipo de estilo (por exemplo `penteado`) em `comum/` retorna `ErrTipoInvalido`.

* [ ] `AC-E07.7`: **Herança de estilo**
  * **Dado que:** estilo B com `herda: A` e itens de A marcados `herdavel: true`.
  * **Quando:** `Buscar` roda em B.
  * **Então:** inclui os itens herdáveis de A; uma herança em ciclo retorna `ErrHerancaCiclica`.

* [ ] `AC-E07.8`: **Validação de estilo**
  * **Dado que:** estilo `rascunho` com referências e 1 ou mais personagens de teste aprovados pelo usuário (Stop-Gate humano gravado em arquivo).
  * **Quando:** `ValidarEstilo` roda.
  * **Então:** muda `status` para `validado` e grava a data e o arquivo de aprovação; sem o arquivo de aprovação, retorna erro tipado e o estilo continua `rascunho`.

* [ ] `AC-E07.9`: **Gênero e faixa etária do personagem**
  * **Dado que:** item do tipo `personagem` com `genero` (`masculino`, `feminino` ou `neutro`) e `faixa_etaria` (`crianca`, `adulto` ou `idoso`).
  * **Quando:** o item é validado e buscado por `BuscarPersonagens`.
  * **Então:** valor fora do vocabulário retorna erro tipado; personagem sem os dois campos é recusado; o filtro por gênero e faixa devolve só os itens que casam; o personagem de teste `crianca` declara `cast.proporcao` com `cabeca` e `corpo`.

### Épico 08: Contratos Novos (R8)

* [ ] `AC-E08.1`: **Pré-voo lê modo e estilo**
  * **Dado que:** `projeto.json` v2 com `modo_producao` e `estilo` `validado`.
  * **Quando:** o pré-voo roda.
  * **Então:** aprova; com estilo `rascunho` ou inexistente, bloqueia com erro tipado; com `estilo_visual` legado, aprova e emite `AVISO_ESTILO_VISUAL_LEGADO`.

* [ ] `AC-E08.2`: **`cenas.json` v2 aponta para o elenco e para o estilo**
  * **Dado que:** elemento de personagem com `personagem_id`, `pose` e `humor`.
  * **Quando:** `ValidarStoryboard` roda.
  * **Então:** exige `personagem_id` presente em `elenco.json` e `pose` e `humor` existentes nos inventários do estilo; violação retorna erro tipado por elemento.

* [ ] `AC-E08.3`: **`elenco.json` só usa itens do inventário e partes suportadas**
  * **Dado que:** bloco `cast` com `hair` ou `hat` fora de `partes_suportadas` do renderizador, ou item inexistente.
  * **Quando:** `ValidarElenco` roda.
  * **Então:** retorna `ErrParteNaoSuportada` ou `ErrEstiloInexistente`, listando o campo.

### Épico 09: Peça 2 Reaberta (R9)

* [ ] `AC-E09.1`: **Etapa 3 vetorial monta o elenco do inventário, sem Flow**
  * **Dado que:** estilo vetorial `validado` e inventário com itens aprovados para as tags do roteiro.
  * **Quando:** a Etapa 3 roda.
  * **Então:** grava `elenco.json` com `cast`, não chama a interface `Flow` e registra 0 crédito.

* [ ] `AC-E09.2`: **Item novo entra como não aprovado**
  * **Dado que:** o roteiro pede uma parte que o inventário não tem (por exemplo `adorno` colar de osso).
  * **Quando:** a Etapa 3 cria o item.
  * **Então:** grava com `aprovado: false`, avisa a lacuna ao usuário e só passa a `true` depois do Stop-Gate humano.

* [ ] `AC-E09.3`: **Caminho de recorte continua**
  * **Dado que:** estilo com renderizador de recorte.
  * **Quando:** a Etapa 3 roda.
  * **Então:** executa fichas e `rembg` como hoje, e todos os testes atuais de `elenco_test.go`, `fichas_test.go` e `drivers_test.go` seguem verdes.

* [ ] `AC-E09.4`: **Fundos consultam o inventário antes de gerar**
  * **Dado que:** inventário do estilo com um cenário aprovado de tags `banco-central`, `noite`.
  * **Quando:** a Etapa 4 precisa desse cenário.
  * **Então:** reutiliza o item (copia o arquivo, confere o hash) e não enfileira geração; fundo novo gerado entra no inventário com prompt, tags e hash.

* [ ] `AC-E09.5`: **Entrega suspeita do Flow não entra**
  * **Dado que:** job de imagem concluído sem a confirmação de renomeação no Flow, ou arquivo com hash igual ao de outro item de tipo diferente, ou fora do padrão de estilo (checagem de cor média e dimensão).
  * **Quando:** o download é registrado.
  * **Então:** marca o item como `suspeito`, não o aprova e exige confirmação humana por folha de contato (achado A11).

### Épico 10: Motor Portado (R10)

* [ ] `AC-E10.1`: **Motor fora do git enquanto a licença estiver pendente**
  * **Dado que:** `motor-local/LICENCA.json` com `status: pendente`.
  * **Quando:** o teste de governança roda `git check-ignore motor-local`.
  * **Então:** a pasta está ignorada; com `status: liberada` ou `reescrito`, o teste permite mover para `internal/video/motor/`.

* [ ] `AC-E10.2`: **Pré-voo do motor**
  * **Dado que:** motor, `python`, Playwright, Chromium e `ffmpeg` configurados.
  * **Quando:** o pré-voo roda.
  * **Então:** aprova; a ausência de qualquer um retorna `ErrMotorAusente` com a instrução de instalação.

* [ ] `AC-E10.3`: **Ponte `cenas.json` para `timeline.json`**
  * **Dado que:** `cenas.json` v2, `elenco.json`, `timings.json` e fundos do inventário da cena 1 do demo.
  * **Quando:** a ponte gera `timeline.json`.
  * **Então:** o `check_timeline.py` do motor retorna `CHECK OK`, cada corte tem `bg` existente e a duração total é igual à da voz (± 1 quadro).

* [ ] `AC-E10.4`: **Render de referência**
  * **Dado que:** a timeline da cena 1.
  * **Quando:** o motor renderiza em 9:16 a 30 quadros por segundo.
  * **Então:** o MP4 tem 1080x1920, trilha de áudio, duração igual à voz (± 1 quadro) e zero quadros em branco.

* [ ] `AC-E10.5`: **Determinismo**
  * **Dado que:** a mesma timeline.
  * **Quando:** o render roda duas vezes.
  * **Então:** o cache por cena tem os mesmos hashes e os arquivos de cena são idênticos.

* [ ] `AC-E10.6`: **Luz de contorno por estilo**
  * **Dado que:** `estilo.json` com `luz.rim`.
  * **Quando:** a ponte monta a timeline.
  * **Então:** cada corte recebe `light.rim` do estilo, e um estilo sem `luz` gera corte sem contorno.

* [ ] `AC-E10.7`: **Boca acompanha a fala**
  * **Dado que:** personagem com `fala: true` durante uma linha de voz.
  * **Quando:** o corte é renderizado.
  * **Então:** os quadros do trecho falado mostram 2 ou mais estados distintos de boca (comparados por diferença de imagem na região da boca) e o trecho sem fala mostra 1.

* [ ] `AC-E10.8`: **Proporção de criança**
  * **Dado que:** personagem com `cast.proporcao` (`cabeca` maior que 1 e `corpo` menor que 1).
  * **Quando:** o corte é renderizado.
  * **Então:** a cabeça do quadro tem raio maior e o corpo altura menor que as do adulto do mesmo estilo (medidos na imagem); sem `proporcao` o desenho não muda.

### Épico 11: Renomeação com Alias (R11)

* [ ] `AC-E11.1`: **Comandos novos e antigos**
  * **Dado que:** `video-roteiro`, `video-ativos` e `video-render` registrados.
  * **Quando:** `diorama-roteiro` ou `diorama-ativos` é executado.
  * **Então:** executa o mesmo código e imprime aviso de depreciação em `stderr`, com exit code igual ao do comando novo.

* [ ] `AC-E11.2`: **Pacotes movidos sem regressão**
  * **Dado que:** `internal/diorama/*` movido para `internal/video/*`.
  * **Quando:** `go build ./...`, `go vet ./...` e `go test -count=1 ./...` rodam.
  * **Então:** exit code 0 e `grep "internal/diorama"` não acha import em código Go.

* [ ] `AC-E11.3`: **Registro e skills migrados**
  * **Dado que:** `skills/legoloop-video-*` e `processos/legoloop-video-*` criados e os nomes antigos redirecionando.
  * **Quando:** `legoloop espinha verificar` e `legoloop fechar` rodam.
  * **Então:** exit code 0, e `docs/PROCESSOS_CANONIZADOS.md` registra a migração com data.

* [ ] `AC-E11.4`: **Histórico preservado**
  * **Dado que:** documentos arquivados de 2026-11-24 e 2026-11-25.
  * **Quando:** a nota de redirecionamento é aplicada.
  * **Então:** `git diff` mostra mudanças só nas linhas iniciais de cada arquivo.

### Épico 12: Prova de Reprodução (R12)

* [ ] `AC-E12.1`: **Cena 1 pelo processo**
  * **Dado que:** estilo molde-de-pílula `validado` e inventário com os 9 fundos.
  * **Quando:** o processo (não um script solto) renderiza a cena 1 do `diorama-demo`.
  * **Então:** gera o MP4 com voz, e o usuário aprova a revisão visual (Stop-Gate humano gravado em arquivo).

* [ ] `AC-E12.2`: **First Brands reproduz**
  * **Dado que:** o roteiro do First Brands e o inventário do estilo.
  * **Quando:** o processo gera a timeline completa.
  * **Então:** `check_timeline.py` dá `CHECK OK`, mediana entre 2,5 e 4 s e cortes por minuto entre 12 e 21.

* [ ] `AC-E12.3`: **Segundo vídeo usa e alimenta o inventário**
  * **Dado que:** vídeo de tema diferente.
  * **Quando:** o elenco é criado.
  * **Então:** reutiliza 3 ou mais itens do primeiro vídeo e grava 3 ou mais itens novos com `origem` do segundo.

### Critérios Não-Funcionais

* [ ] `AC-NF.1`: cobertura `>= 95%` em cada pacote novo (`inventario`, e as partes novas de `ativos` e `render`).
* [ ] `AC-NF.2`: nenhuma dependência nova em `go.mod`, `go.sum` intactos; funções dentro do limite de `tools/funclen`.
* [ ] `AC-NF.3`: os testes passam em Windows, Linux e macOS (caminhos com `filepath`, sem shell específico).

## 6. Análise de Modos de Falha (FMEA)

Escala 1 a 10. RPN = Severidade × Ocorrência × Detecção. Acima de 100 exige critério binário dedicado.

| # | Modo de falha | Efeito | S | O | D | RPN | Mitigação | AC |
| :--- | :--- | :--- | :-: | :-: | :-: | :-: | :--- | :--- |
| F14 | Licença do motor negada ou indefinida | Não dá para publicar o repositório | 9 | 5 | 3 | 135 | Motor em `motor-local/` fora do git; plano B de reescrita limpa de `puppet.js` e `sets.js` | `AC-E10.1` |
| F15 | Flow entrega arquivo errado (achado A11) | Fundo ou ficha errada no vídeo | 7 | 6 | 3 | 126 | Marcar `suspeito` e exigir folha de contato | `AC-E09.5` |
| F16 | Estilo muda depois de itens aprovados | Itens ficam fora do padrão | 6 | 4 | 3 | 72 | `estilo_versao` no item; mudança maior marca `revisar` | `AC-E07.2` |
| F17 | Inventário enche de duplicatas | Consulta lenta e confusa | 5 | 7 | 4 | 140 | Hash por arquivo e busca antes de criar | `AC-E07.5`, `AC-E09.4` |
| F18 | Renomeação quebra processos canonizados | Fechamento reprova | 8 | 5 | 2 | 80 | Alias, `espinha verificar` e migração registrada | `AC-E11.3` |
| F19 | Item com forma que o renderizador não desenha | Quadro sem a peça ou erro de JS | 6 | 6 | 2 | 72 | `partes_suportadas` no estilo | `AC-E08.3` |
| F20 | Estilo não validado usado em produção | Vídeo fora do padrão aprovado | 7 | 4 | 2 | 56 | Pré-voo bloqueia | `AC-E08.1` |
| F21 | Python, Playwright ou Chromium ausentes | Render não roda no agente ou SO | 6 | 5 | 2 | 60 | Pré-voo com instrução de instalação | `AC-E10.2` |

## 7. Rastreabilidade e Comandos

Requisitos: R7 inventário e estilos, R8 contratos, R9 Peça 2, R10 motor, R11 renomeação, R12 prova. Cada `AC-E07.*` a `AC-E12.*` aponta para um único requisito pela numeração do épico.

### Comandos de verificação mecânica
- **Build:** `go build -o bin/legoloop.exe ./cmd/legoloop`
- **Testes e cobertura:** `go test -count=1 -cover ./internal/video/...` e `go test -count=1 ./...`
- **Estática:** `go vet ./...`, `go run ./tools/funclen internal`, `go run ./tools/taxlinter internal cmd loops skills docs/templates`, `gofmt -l .`
- **Dependências:** `git diff --exit-code -- go.mod go.sum`
- **Espinha e fechamento:** `bin/legoloop.exe espinha verificar` e `bin/legoloop.exe fechar "legoloop-video-animado"`

### Ordem de tickets (um por vez, com autorização)
1. Pacote e esquema de inventário (E07).
2. Estilo molde-de-pílula com referências e validação humana (E07.8).
3. Contratos novos e pré-voo (E08).
4. Peça 2 reaberta com `cast`, fundos e entrega suspeita (E09).
5. Motor em `motor-local/`, ponte e render de referência (E10).
6. Voz, mix e cortes sobre o motor (E01 a E06 da v1.4.0, ajustados).
7. Renomeação com alias (E11).
8. Estilo palito (D5) e provas E12.

### Critério de Done global
- [ ] Todos os critérios `AC-E07.1` a `AC-E12.3`, `AC-NF.1` a `AC-NF.3` e os `AC-E01.*`, `AC-E02.*`, `AC-E04.*`, `AC-E05.*`, `AC-E06.*` mantidos passam.
- [ ] Cobertura `>= 95%` nos pacotes novos e todos os portões mecânicos com exit code 0.
- [ ] `legoloop fechar "legoloop-video-animado"` com 5 de 5.

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
*A SPEC `v2.0.0-rascunho` (6 épicos novos, 30 critérios binários novos, FMEA com 8 modos de falha novos) atende às suas expectativas para avançarmos ao planejamento (`PLAN`, `ROUTE` e `DAG`)?*

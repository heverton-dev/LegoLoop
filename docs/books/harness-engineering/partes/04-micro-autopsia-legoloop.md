# Capítulo 4: Autópsia do LegoLoop — Onde Estão os Erros no Nosso Próprio Repositório

Vamos olhar no espelho com honestidade cirúrgica. O LegoLoop é um projeto impressionante, com ideias arquiteturais brilhantes (DAGs isomórficos, Tríade de Interfaces, compilação de livros via LegoPress e governança ISO de artefatos). 

Mas se ele é tão bem pensado, por que a experiência acabou de colapsar nesta sessão, gerando um dossiê com campos nulos e arquivos espalhados na raiz?

Vamos dissecar a anatomia dessa falha em 4 dimensões:

---

## 1. Onde Erramos: O Código Go do `video_rea.go`

Abra o arquivo `internal/cli/video_rea.go` nas linhas 70 a 84:
```go
func runIngest(args []string, d Deps) int {
    fs := newFlags("video-rea ingest", d)
    rawURL := fs.String("url", "", "URL do video")
    out := fs.String("out", "", "arquivo de saida")
    // ...
    vInput := videorea.VideoInput{URL: *rawURL}
    return saveJSONOutput(vInput, *out, d)
}
```

Olhe o que o comando `ingest` faz de verdade:
* Ele pega a string da URL;
* Cria um struct `VideoInput{URL: url}`;
* Grava no disco e encerra com **Exit Code 0**.

Ele **não** rodou o `yt-dlp`. Ele **não** baixou a descrição do vídeo. Ele **não** inspecionou o comentário fixado do autor. Ele **não** verificou se havia links do Google Drive ou de GitHub.

**Esse foi o primeiro erro fatal de Harness:** O código Go prometeu que fez um "ingest", mas só salvou uma string em um arquivo JSON. O resto do trabalho de investigação foi jogado nas costas do LLM via terminal. E o LLM, operando sob sua natureza preguiçosa, fez o mínimo possível para fechar o turno.

---

## 2. Por que Erramos: A Ausência de Asserções no Go

Olhe agora a função `executeForensics`:
```go
dossier, err := videorea.ExtractForensicsWithBrain(ctx, vInput, transcript, desc, b)
if err != nil {
    return 1
}
return saveJSONOutput(dossier, out, d)
```

O Go confiou que o struct `dossier` retornado pelo modelo era bom o suficiente simplesmente porque `err == nil`. Ele não checou:
* `if len(dossier.ToolsIdentified) == 0`
* `if len(dossier.ExplicitPrompts) == 0`
* `if len(dossier.OCRFindings) == 0`

Se o código Go tivesse uma linha simples:
```go
if len(dossier.ExplicitPrompts) == 0 {
    fmt.Fprintln(d.Stderr, "ERRO MECÂNICO: Nenhum prompt foi identificado no vídeo. O dossiê é inválido.")
    return 1
}
```
O comando teria quebrado no terminal com **Exit Code 1**. O modelo não teria conseguido avançar para a curadoria. Ele teria sido obrigado a procurar onde estavam os prompts antes de falar qualquer coisa para você no chat.

---

## 3. Quando Erramos: A Dependência da Atenção Humana

Erramos no momento em que desenhamos o pipeline achando que o operador humano serviria como um validador de saída de JSON. 

O operador humano deve aprovar **decisões de negócio** (*"Gostei do tema ET de Varginha"*, *"Aprovo o orçamento R$ 0,00"*). O operador humano **nunca** deve ser acionado para verificar se um campo do JSON está nulo ou se o arquivo foi parar na pasta errada. Essa verificação é papel da máquina.

Quando você teve que abrir o `dossier.json`, ver que estava nulo e xingar no chat, o Harness já havia falhado há 3 turnos atrás.

---

## 4. Como Corrigir no LegoLoop Definitivamente

A reforma necessária no LegoLoop não exige reescrever o projeto inteiro. Ela exige instalar **4 travas mecânicas no Go**:

1. **Path Jail no `bundle.go`:** O motor de bundle deve ser a única entidade que cria caminhos de arquivos. Se um arquivo for gravado fora de `docs/videorea/<slug>/`, o Go rejeita com erro.
2. **Ingest Mecânico Real:** O `legoloop video-rea ingest` deve invocar `yt-dlp` internamente com `--dump-json`, extrair descrição, comentários e já baixar os links anexados automaticamente antes de passar a bola para a IA.
3. **Validação Estrita de Dossiê:** Criar o método `ValidateStrict()` que recusa salvar qualquer dossiê com campos nulos ou listas vazias.
4. **Substituir Regras de Prosa por Gates de Binário:** Em vez de ter 14 regras no `AGENTS.md` esperando que o modelo lembre de todas, criar um comando de validação rápida `legoloop audit-gate` que roda antes de cada commit e reprova o build se qualquer regra mecânica for violada.

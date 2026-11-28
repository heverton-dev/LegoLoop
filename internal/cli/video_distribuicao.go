package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/video/distribuicao"
)

const usoVideoDistribuicao = `Uso: legoloop video-distribuicao <subcomando> [opções]

Subcomandos (todos aceitam --projeto DIR):
  prevoo                                   confere pasta de saída, ffmpeg, python com Pillow e as variáveis do Postiz
  revisao folhas --video REL --versao N    folhas de contato e pacote do revisor sem contexto
  revisao parecer --arquivo ARQ            grava review/vN_parecer.json (não sobrescreve)
  revisao liberar --versao N               exit 0 só se o parecer for ENVIAR
  vitrine pontuar --palavra P TÍTULO...    nota cada título
  vitrine gravar --arquivo ARQ --titulo N  valida o pacote e grava publicacao.md
  capa --arquivo ARQ                       gera a capa 1280x720 (fundo sem texto + texto por código)
  acervo adicionar --arquivo ARQ           adiciona item (regra de originalidade)
  acervo listar                            lista o acervo
  acervo autorizar --id ID [--revogar]     grava ou revoga a autorização de publicação
  agenda [--config ARQ] [--agora RFC3339] [--aplicar]   plano quente/perene; --aplicar grava no acervo
  publicar canais                          lista os canais do Postiz
  publicar enviar --id ID --arquivo ARQ [--executar]   dry-run por padrão; --executar envia de verdade
  comentarios rascunhar --entrada ARQ --id ID --texto T   grava rascunho de resposta
  comentarios decidir --id ID (--aprovar|--ignorar) [--editado T]
  comentarios listar --entrada ARQ         comentários sem resposta
  aprendizado --metricas ARQ               grava acervo/aprendizado.md
  canal textos                             gera colar.md a partir de canal/textos.json e redes.json
  canal imagens                            gera imagens de perfil e banners das redes em canal/imagens/
`

func cmdVideoDistribuicao(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprint(d.Stdout, usoVideoDistribuicao)
		return 0
	}
	mapa := map[string]func(context.Context, []string, Deps) int{
		"prevoo": distPrevoo, "revisao": distRevisao, "vitrine": distVitrine, "capa": distCapa,
		"acervo": distAcervo, "agenda": distAgenda, "publicar": distPublicar, "comentarios": distComentarios,
		"aprendizado": distAprendizado, "canal": distCanal,
	}
	f, ok := mapa[args[0]]
	if !ok {
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s", args[0], usoVideoDistribuicao)
		return 2
	}
	return f(ctx, args[1:], d)
}

// opcoesDist reune as flags comuns lidas por todos os subcomandos.
type opcoesDist struct {
	projeto, arquivo, id, video, entrada, texto, editado, titulo, palavra, config, agora, raiz, metricas string
	versao                                                                                               int
	passo                                                                                                float64
	aprovar, ignorar, revogar, aplicar, executar                                                         bool
	restantes                                                                                            []string
}

func lerOpcoesDist(nome string, args []string, d Deps) (opcoesDist, bool) {
	var o opcoesDist
	fs := flag.NewFlagSet("video-distribuicao "+nome, flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	for alvo, nomeFlag := range map[*string]string{&o.projeto: "projeto", &o.arquivo: "arquivo", &o.id: "id", &o.video: "video", &o.entrada: "entrada",
		&o.texto: "texto", &o.editado: "editado", &o.titulo: "titulo", &o.palavra: "palavra", &o.config: "config", &o.agora: "agora", &o.metricas: "metricas"} {
		fs.StringVar(alvo, nomeFlag, "", nomeFlag)
	}
	fs.StringVar(&o.raiz, "raiz", ".", "raiz onde procurar arquivos .env*")
	fs.IntVar(&o.versao, "versao", 0, "versão do render")
	fs.Float64Var(&o.passo, "passo", 2, "segundos entre quadros da folha de contato")
	for alvo, nomeFlag := range map[*bool]string{&o.aprovar: "aprovar", &o.ignorar: "ignorar", &o.revogar: "revogar", &o.aplicar: "aplicar", &o.executar: "executar"} {
		fs.BoolVar(alvo, nomeFlag, false, nomeFlag)
	}
	if err := fs.Parse(args); err != nil {
		return o, false
	}
	if o.projeto == "" {
		fmt.Fprintf(d.Stderr, "video-distribuicao %s: --projeto é obrigatório\n", nome)
		return o, false
	}
	o.restantes = fs.Args()
	return o, true
}

func distFalha(d Deps, nome string, err error) int {
	fmt.Fprintf(d.Stderr, "video-distribuicao %s: %v\n", nome, err)
	return 1
}

func distJSON(d Deps, v any) int {
	dados, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return distFalha(d, "json", err)
	}
	fmt.Fprintln(d.Stdout, string(dados))
	return 0
}

func lerArquivoJSON(caminho string, alvo any) error {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}
	return json.Unmarshal(dados, alvo)
}

// subDist separa o verbo do grupo (revisao, vitrine, acervo...) do restante dos argumentos.
func subDist(grupo string, args []string, d Deps) (string, []string, bool) {
	if len(args) == 0 {
		fmt.Fprintf(d.Stderr, "video-distribuicao %s: informe o subcomando\n", grupo)
		return "", nil, false
	}
	return args[0], args[1:], true
}

func distPrevoo(_ context.Context, args []string, d Deps) int {
	o, ok := lerOpcoesDist("prevoo", args, d)
	if !ok {
		return 2
	}
	if err := os.MkdirAll(filepath.Join(o.projeto, "acervo"), 0o755); err != nil {
		fmt.Fprintf(d.Stdout, "FALTA: pasta acervo/ (%v)\n", err)
		return 1
	}
	fmt.Fprintln(d.Stdout, "OK: pasta acervo/")
	for _, f := range []string{"ffmpeg", "ffprobe"} {
		distAviso(d, f, func() bool { _, err := exec.LookPath(f); return err == nil })
	}
	distAviso(d, "python com Pillow (capas)", func() bool { return exec.Command("python", "-c", "import PIL").Run() == nil })
	for _, v := range []string{"POSTIZ_API_KEY", "POSTIZ_URL"} {
		_, tem := distribuicao.SegredoDoAmbiente(v, o.raiz, os.Getenv)
		distAviso(d, v, func() bool { return tem })
	}
	fmt.Fprintln(d.Stdout, "pre-voo aprovado")
	return 0
}

func distAviso(d Deps, nome string, ok func() bool) {
	nivel := "OK"
	if !ok() {
		nivel = "AVISO"
	}
	fmt.Fprintf(d.Stdout, "%s: %s\n", nivel, nome)
}

func distRevisao(ctx context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("revisao", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("revisao "+sub, resto, d)
	if !ok {
		return 2
	}
	switch sub {
	case "folhas":
		return revisaoFolhas(ctx, o, d)
	case "parecer":
		return revisaoParecer(o, d)
	case "liberar":
		if err := distribuicao.Liberar(o.projeto, o.versao, time.Now()); err != nil {
			return distFalha(d, "revisao liberar", err)
		}
		fmt.Fprintf(d.Stdout, "v%d liberada para a vitrine\n", o.versao)
		return 0
	}
	fmt.Fprintf(d.Stderr, "revisao: subcomando desconhecido %q\n", sub)
	return 2
}

func revisaoFolhas(ctx context.Context, o opcoesDist, d Deps) int {
	folhas, err := distribuicao.GerarFolhas(ctx, distribuicao.ExecSistema{}, distribuicao.DuracaoFFprobe, o.projeto, o.video, o.passo)
	if err != nil {
		return distFalha(d, "revisao folhas", err)
	}
	caminho, err := distribuicao.GravarPacote(o.projeto, o.versao, folhas, o.passo)
	if err != nil {
		return distFalha(d, "revisao folhas", err)
	}
	fmt.Fprintf(d.Stdout, "%d folhas; pacote do revisor em %s\n", len(folhas), caminho)
	return 0
}

func revisaoParecer(o opcoesDist, d Deps) int {
	var p distribuicao.Parecer
	if err := lerArquivoJSON(o.arquivo, &p); err != nil {
		return distFalha(d, "revisao parecer", err)
	}
	if err := distribuicao.GravarParecer(o.projeto, p); err != nil {
		return distFalha(d, "revisao parecer", err)
	}
	fmt.Fprintf(d.Stdout, "parecer v%d gravado: %s\n", p.Versao, p.Veredito)
	return 0
}

func distVitrine(_ context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("vitrine", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("vitrine "+sub, resto, d)
	if !ok {
		return 2
	}
	switch sub {
	case "pontuar":
		return distJSON(d, distribuicao.PontuarTitulos(o.restantes, o.palavra))
	case "gravar":
		var v distribuicao.Vitrine
		if err := lerArquivoJSON(o.arquivo, &v); err != nil {
			return distFalha(d, "vitrine gravar", err)
		}
		n := 0
		if _, err := fmt.Sscan(o.titulo, &n); err != nil {
			return distFalha(d, "vitrine gravar", fmt.Errorf("--titulo deve ser o índice do título (0 a 4)"))
		}
		if err := distribuicao.GravarPublicacao(o.projeto, v, n); err != nil {
			return distFalha(d, "vitrine gravar", err)
		}
		fmt.Fprintln(d.Stdout, "publicacao.md gravado")
		return 0
	}
	fmt.Fprintf(d.Stderr, "vitrine: subcomando desconhecido %q\n", sub)
	return 2
}

func distCapa(ctx context.Context, args []string, d Deps) int {
	o, ok := lerOpcoesDist("capa", args, d)
	if !ok {
		return 2
	}
	var p distribuicao.PedidoCapa
	if err := lerArquivoJSON(o.arquivo, &p); err != nil {
		return distFalha(d, "capa", err)
	}
	res, err := distribuicao.GerarCapa(ctx, distribuicao.ExecSistema{}, o.projeto, p)
	if err != nil {
		return distFalha(d, "capa", err)
	}
	return distJSON(d, res)
}

func distAcervo(_ context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("acervo", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("acervo "+sub, resto, d)
	if !ok {
		return 2
	}
	a, err := distribuicao.CarregarAcervo(o.projeto)
	if err != nil {
		return distFalha(d, "acervo", err)
	}
	if sub == "listar" {
		return distJSON(d, a)
	}
	if err := distAcervoMudar(&a, sub, o); err != nil {
		return distFalha(d, "acervo "+sub, err)
	}
	if err := a.Gravar(o.projeto); err != nil {
		return distFalha(d, "acervo "+sub, err)
	}
	fmt.Fprintf(d.Stdout, "acervo atualizado: %d itens\n", len(a.Itens))
	return 0
}

func distAcervoMudar(a *distribuicao.Acervo, sub string, o opcoesDist) error {
	switch sub {
	case "adicionar":
		var it distribuicao.Item
		if err := lerArquivoJSON(o.arquivo, &it); err != nil {
			return err
		}
		return a.Adicionar(it)
	case "autorizar":
		return a.Autorizar(o.id, !o.revogar)
	}
	return fmt.Errorf("subcomando desconhecido %q", sub)
}

func distAgenda(_ context.Context, args []string, d Deps) int {
	o, ok := lerOpcoesDist("agenda", args, d)
	if !ok {
		return 2
	}
	cfg, agora, err := agendaEntradas(o)
	if err != nil {
		return distFalha(d, "agenda", err)
	}
	a, err := distribuicao.CarregarAcervo(o.projeto)
	if err != nil {
		return distFalha(d, "agenda", err)
	}
	reels := a.GerarReels()
	plano, err := distribuicao.Planejar(a, cfg, agora)
	if err == nil && o.aplicar {
		err = a.Aplicar(plano)
		if err == nil {
			err = a.Gravar(o.projeto)
		}
	}
	if err != nil {
		return distFalha(d, "agenda", err)
	}
	fmt.Fprintf(d.Stderr, "reels derivados: %d; aplicado: %v\n", reels, o.aplicar)
	return distJSON(d, plano)
}

func agendaEntradas(o opcoesDist) (distribuicao.ConfigAgenda, time.Time, error) {
	var cfg distribuicao.ConfigAgenda
	if o.config != "" {
		if err := lerArquivoJSON(o.config, &cfg); err != nil {
			return cfg, time.Time{}, err
		}
	}
	if o.agora == "" {
		return cfg, time.Now().UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, o.agora)
	return cfg, t, err
}

func distPublicar(ctx context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("publicar", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("publicar "+sub, resto, d)
	if !ok {
		return 2
	}
	if sub != "canais" && sub != "enviar" {
		fmt.Fprintf(d.Stderr, "publicar: subcomando desconhecido %q\n", sub)
		return 2
	}
	c, err := distribuicao.NovoCliente(o.raiz, os.Getenv)
	if sub == "enviar" && o.executar == false {
		c, err = distribuicao.Cliente{}, nil
	}
	if err != nil {
		return distFalha(d, "publicar", err)
	}
	if sub == "canais" {
		canais, err := c.Canais(ctx)
		if err != nil {
			return distFalha(d, "publicar canais", err)
		}
		return distJSON(d, canais)
	}
	return distEnviar(ctx, c, o, d)
}

func distEnviar(ctx context.Context, c distribuicao.Cliente, o opcoesDist, d Deps) int {
	var p distribuicao.PedidoEnvio
	if err := lerArquivoJSON(o.arquivo, &p); err != nil {
		return distFalha(d, "publicar enviar", err)
	}
	p.DryRun = !o.executar
	a, err := distribuicao.CarregarAcervo(o.projeto)
	if err != nil {
		return distFalha(d, "publicar enviar", err)
	}
	res, err := c.PublicarDoAcervo(ctx, &a, o.id, p)
	if err != nil {
		return distFalha(d, "publicar enviar", err)
	}
	if !p.DryRun {
		if err := distEnvioGravar(o, p, res, a); err != nil {
			return distFalha(d, "publicar enviar", err)
		}
	}
	return distJSON(d, res)
}

func distComentarios(_ context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("comentarios", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("comentarios "+sub, resto, d)
	if !ok {
		return 2
	}
	r, err := distribuicao.CarregarRespostas(o.projeto)
	if err != nil {
		return distFalha(d, "comentarios", err)
	}
	switch sub {
	case "listar":
		cs, err := distribuicao.CarregarComentarios(o.entrada)
		if err != nil {
			return distFalha(d, "comentarios listar", err)
		}
		return distJSON(d, distribuicao.Pendentes(cs, r))
	case "rascunhar", "decidir":
		if err := distComentarioMudar(&r, sub, o); err != nil {
			return distFalha(d, "comentarios "+sub, err)
		}
		if err := r.Gravar(o.projeto); err != nil {
			return distFalha(d, "comentarios "+sub, err)
		}
		fmt.Fprintf(d.Stdout, "respostas: %d, aprovadas: %d\n", len(r.Itens), len(r.Aprovadas()))
		return 0
	}
	fmt.Fprintf(d.Stderr, "comentarios: subcomando desconhecido %q\n", sub)
	return 2
}

func distComentarioMudar(r *distribuicao.Respostas, sub string, o opcoesDist) error {
	if sub == "decidir" {
		if o.aprovar == o.ignorar {
			return fmt.Errorf("informe --aprovar ou --ignorar")
		}
		return r.Decidir(o.id, o.aprovar, o.editado)
	}
	cs, err := distribuicao.CarregarComentarios(o.entrada)
	if err != nil {
		return err
	}
	for _, c := range cs {
		if c.ID == o.id {
			return r.Rascunhar(c, o.texto)
		}
	}
	return fmt.Errorf("comentário %q não está em %s", o.id, o.entrada)
}

func distAprendizado(_ context.Context, args []string, d Deps) int {
	o, ok := lerOpcoesDist("aprendizado", args, d)
	if !ok {
		return 2
	}
	var ms []distribuicao.Metrica
	if err := lerArquivoJSON(o.metricas, &ms); err != nil {
		return distFalha(d, "aprendizado", err)
	}
	caminho, err := distribuicao.GravarAprendizado(o.projeto, ms, time.Now().UTC())
	if err != nil {
		return distFalha(d, "aprendizado", err)
	}
	fmt.Fprintf(d.Stdout, "aprendizado gravado em %s\n", caminho)
	return 0
}

// distEnvioGravar grava o acervo e a linha de auditoria de um envio real.
func distEnvioGravar(o opcoesDist, p distribuicao.PedidoEnvio, res distribuicao.ResultadoEnvio, a distribuicao.Acervo) error {
	if err := a.Gravar(o.projeto); err != nil {
		return err
	}
	return distribuicao.RegistrarEnvio(o.projeto, distribuicao.RegistroEnvio{ItemID: o.id, Canal: p.Canal, Visibilidade: p.Visibilidade,
		Modo: fmt.Sprint(res.Resumo["modo"]), EnviadoEm: time.Now().UTC(), Resposta: string(res.Resposta)})
}

func distCanal(ctx context.Context, args []string, d Deps) int {
	sub, resto, ok := subDist("canal", args, d)
	if !ok {
		return 2
	}
	o, ok := lerOpcoesDist("canal "+sub, resto, d)
	if !ok {
		return 2
	}
	canalDir := filepath.Join(o.projeto, "canal")
	switch sub {
	case "textos":
		rel, err := distribuicao.GerarTextosCanal(canalDir)
		if err != nil {
			return distFalha(d, "canal textos", err)
		}
		fmt.Fprintf(d.Stdout, "textos do canal gerados: %s (excedidos: %d)\n", rel.Saida, rel.Excedidos)
		return 0
	case "imagens":
		destLogo := filepath.Join(o.projeto, "imagens", "logo.png")
		rel, err := distribuicao.GerarImagensCanal(ctx, distribuicao.ExecSistema{}, canalDir, destLogo)
		if err != nil {
			return distFalha(d, "canal imagens", err)
		}
		fmt.Fprintf(d.Stdout, "imagens do canal geradas: %d arquivos (logo: %s)\n", len(rel.Arquivos), rel.LogoSincronizado)
		return 0
	}
	fmt.Fprintf(d.Stderr, "canal: subcomando desconhecido %q\n", sub)
	return 2
}

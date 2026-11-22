package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/progresso"
)

// opcoesVerificarCLI reúne as flags de `progresso verificar`.
type opcoesVerificarCLI struct {
	arquivo, raiz           string
	maxIdade                time.Duration
	sondar, reconciliar, js bool
}

func cmdProgressoVerificar(ctx context.Context, args []string, d Deps) int {
	var o opcoesVerificarCLI
	fs := newFlags(`progresso verificar`, d)
	fs.StringVar(&o.arquivo, `arquivo`, ``, `arquivo JSON de progresso`)
	fs.StringVar(&o.raiz, `raiz`, `.`, `raiz onde rodam os comandos dos Tickets`)
	fs.DurationVar(&o.maxIdade, `max-idade`, 24*time.Hour, `idade máxima de um Ticket IN_PROGRESS antes de virar órfão (0 desliga)`)
	fs.BoolVar(&o.sondar, `sondar`, false, `roda o Gate dos Tickets não aprovados e aponta os que já passam`)
	fs.BoolVar(&o.reconciliar, `reconciliar`, false, `reexecuta o Gate dos órfãos: verde vira PASSED, vermelho volta a PENDING`)
	fs.BoolVar(&o.js, `json`, false, `saída em JSON`)
	if err := fs.Parse(args); err != nil || o.arquivo == `` {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	achados, err := verificarProgresso(ctx, o)
	if err != nil {
		return fail(d, err)
	}
	return relatarAchados(d, achados, o.js)
}

func verificarProgresso(ctx context.Context, o opcoesVerificarCLI) ([]progresso.Achado, error) {
	if o.reconciliar {
		if err := reconciliarProgresso(ctx, o); err != nil {
			return nil, err
		}
	}
	p, err := progresso.Carregar(o.arquivo)
	if err != nil {
		return nil, err
	}
	opt := progresso.OpcoesVerificar{Agora: time.Now().UTC(), IdadeMaxima: o.maxIdade}
	var achados []progresso.Achado
	if md, errMD := os.ReadFile(progresso.CaminhoEspelhoMD(o.arquivo)); errMD != nil {
		achados = append(achados, progresso.Achado{Regra: `espelho-ausente`, Motivo: errMD.Error()})
	} else {
		texto := string(md)
		opt.Espelho = &texto
	}
	achados = append(achados, progresso.Verificar(p, opt)...)
	if o.sondar {
		achados = append(achados, progresso.Sondar(ctx, p, o.raiz, nil)...)
	}
	return achados, nil
}

func reconciliarProgresso(ctx context.Context, o opcoesVerificarCLI) error {
	return progresso.Atualizar(o.arquivo, func(p *progresso.Progresso) error {
		agora := time.Now().UTC()
		ids := progresso.Orfaos(*p, agora, o.maxIdade)
		progresso.ReconciliarOrfaos(ctx, p, o.raiz, agora.Format(time.RFC3339), ids, nil)
		return nil
	})
}

func relatarAchados(d Deps, achados []progresso.Achado, asJSON bool) int {
	if asJSON {
		if code := printJSON(d, map[string]any{`ok`: len(achados) == 0, `achados`: achados}); code != 0 {
			return code
		}
	} else if len(achados) == 0 {
		fmt.Fprintln(d.Stdout, `progresso coerente`)
	}
	if len(achados) == 0 {
		return 0
	}
	if !asJSON {
		for _, a := range achados {
			fmt.Fprintf(d.Stdout, "%s %s %s\n", a.Regra, a.Ticket, a.Motivo)
		}
	}
	return 1
}

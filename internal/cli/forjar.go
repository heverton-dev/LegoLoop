package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/forja"
	"github.com/heverton-dev/LegoLoop/internal/skills"
)

type opcoesCmdForjar struct {
	dryRun    bool
	apply     bool
	desfazer  string
	perfil    string
	sem       listFlag
	home      string
	verificar bool
	corrigir  bool
	cofre     string
	orcamento string
	asJSON    bool
}

func cmdForjar(ctx context.Context, args []string, d Deps) int {
	op, code := parseFlagsForjar(args, d)
	if code != 0 {
		return code
	}
	home, err := resolverHomeInstalar(d, op.home)
	if err != nil {
		return fail(d, err)
	}
	if op.verificar {
		return executarVerificacaoForjar(d, home, op)
	}
	return executarMotorForjar(ctx, d, home, op)
}

func parseFlagsForjar(args []string, d Deps) (opcoesCmdForjar, int) {
	var op opcoesCmdForjar
	fs := newFlags("forjar", d)
	fs.BoolVar(&op.dryRun, "dry-run", false, "simula a Forja sem gravar no disco")
	fs.BoolVar(&op.apply, "apply", false, "aplica a Forja com backup")
	fs.StringVar(&op.desfazer, "desfazer", "", "data do backup a restaurar (AAAA-MM-DD)")
	fs.StringVar(&op.perfil, "perfil", "padrao", "perfil de contexto: minimo, padrao ou gateway")
	fs.Var(&op.sem, "sem", "harness excluido (pode repetir)")
	fs.StringVar(&op.home, "home", "", "diretorio HOME alvo")
	fs.BoolVar(&op.verificar, "verificar", false, "verifica deriva ou orcamento de tokens dos harnesses")
	fs.BoolVar(&op.corrigir, "corrigir", false, "corrige derivas movendo skills eager ao cofre com backup")
	fs.StringVar(&op.cofre, "cofre", "", "pasta de destino das skills no cofre")
	fs.StringVar(&op.orcamento, "orcamento", "", "caminho do JSON de orcamento de tokens")
	fs.BoolVar(&op.asJSON, "json", false, "emite o resultado em JSON")
	if err := fs.Parse(args); err != nil {
		return op, 2
	}
	if op.dryRun && op.apply {
		fmt.Fprintln(d.Stderr, "erro: use apenas --dry-run ou --apply")
		return op, 2
	}
	return op, 0
}

func executarVerificacaoForjar(d Deps, home string, op opcoesCmdForjar) int {
	if strings.TrimSpace(op.orcamento) != "" {
		return executarVerificacaoOrcamento(d, home, op)
	}
	if op.corrigir {
		return executarCorrecaoDeriva(d, home, op)
	}
	return executarDeteccaoDeriva(d, home, op)
}

func executarDeteccaoDeriva(d Deps, home string, op opcoesCmdForjar) int {
	derivas, err := forja.DetectarDeriva(home)
	if err != nil {
		return fail(d, err)
	}
	res := forja.ResultadoDeriva{Aprovado: len(derivas) == 0, Derivas: derivas}
	if op.asJSON {
		if code := printJSON(d, res); code != 0 {
			return code
		}
		if !res.Aprovado {
			return 1
		}
		return 0
	}
	if !res.Aprovado {
		fmt.Fprintln(d.Stderr, "deriva reprovada:")
		fmt.Fprintln(d.Stderr, forja.FormatarDerivas(derivas))
		return 1
	}
	fmt.Fprintln(d.Stdout, forja.FormatarDerivas(derivas))
	return 0
}

func executarCorrecaoDeriva(d Deps, home string, op opcoesCmdForjar) int {
	res, err := forja.CorrigirDeriva(home, op.cofre, "")
	if err != nil {
		return fail(d, err)
	}
	if op.asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "deriva corrigida: %d skill(s) movida(s) ao cofre (backup=%s)\n",
		len(res.Corrigidos), filepath.ToSlash(res.BackupRaiz))
	return 0
}

func executarVerificacaoOrcamento(d Deps, home string, op opcoesCmdForjar) int {
	arq := strings.TrimSpace(op.orcamento)
	orc, err := forja.CarregarOrcamento(arq)
	if err != nil {
		return fail(d, err)
	}
	medidos := medirTokensHome(home)
	res := forja.VerificarOrcamento(orc, medidos)
	if op.asJSON {
		if code := printJSON(d, res); code != 0 {
			return code
		}
		if !res.Aprovado {
			return 1
		}
		return 0
	}
	if !res.Aprovado {
		fmt.Fprintln(d.Stderr, "orcamento reprovado:")
		fmt.Fprintln(d.Stderr, forja.FormatarExcessos(res))
		return 1
	}
	fmt.Fprintln(d.Stdout, forja.FormatarExcessos(res))
	return 0
}

func medirTokensHome(home string) map[string]int {
	grupos, ordem := gruposSkillsCusto(home)
	medidos := make(map[string]int, len(ordem))
	for _, nome := range ordem {
		custo := skills.CustoGrupo(nome, grupos[nome])
		medidos[nome] = custo.Tokens
	}
	return medidos
}

func executarMotorForjar(ctx context.Context, d Deps, home string, op opcoesCmdForjar) int {
	rel, err := forja.Executar(ctx, forja.Opcao{
		Home:     home,
		Perfil:   op.perfil,
		Sem:      expandirListaSem(op.sem),
		Aplicar:  op.apply,
		Desfazer: op.desfazer,
	})
	if err != nil {
		return tratarErroInstalar(d, err)
	}
	if op.asJSON {
		return printJSON(d, rel)
	}
	modo := "DRY-RUN"
	if op.desfazer != "" {
		modo = "DESFAZER"
	} else if op.apply {
		modo = "APLICADO"
	}
	fmt.Fprintf(d.Stdout, "Forjar (%s): veredito=%s perfil=%s etapas=%d\n",
		modo, rel.Veredito, rel.Perfil, len(rel.Etapas))
	return 0
}

func expandirListaSem(itens []string) []string {
	var out []string
	for _, item := range itens {
		for _, parte := range strings.Split(item, ",") {
			limpo := strings.TrimSpace(parte)
			if limpo != "" {
				out = append(out, limpo)
			}
		}
	}
	return out
}

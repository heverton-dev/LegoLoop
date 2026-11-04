package cli

import (
	"errors"
	"fmt"

	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"matriz.perfil_confirmado": "Perfil de %s confirmado. A origem agora é confirmado.",
		"matriz.perfil_definido":   "Perfil de %s definido por você. A origem agora é confirmado.",
	} {
		ste.Default.Register(id, msg)
	}
}

// cmdMatrizPerfil despacha `matriz perfil confirmar|definir <peça>`.
func cmdMatrizPerfil(args []string, d Deps) int {
	if len(args) == 0 || (args[0] != "confirmar" && args[0] != "definir") {
		return matrizUso(d)
	}
	sub := args[0]
	alvo, rest := splitPositional(args[1:])
	fs := newFlags("matriz perfil "+sub, d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	var etapas, naturezas, casos listFlag
	if sub == "definir" {
		fs.Var(&etapas, "etapa", "etapa do vocabulário (repetível)")
		fs.Var(&naturezas, "natureza", "natureza do vocabulário (repetível)")
		fs.Var(&casos, "caso", "caso de uso do vocabulário (repetível)")
	}
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if alvo == "" {
		alvo = fs.Arg(0)
	}
	if alvo == "" {
		return matrizUso(d)
	}
	v, err := resolveVaultFlag(*vault)
	if err != nil {
		return fail(d, err)
	}
	if sub == "confirmar" {
		return perfilResultado(d, "matriz.perfil_confirmado", alvo, errOf(matriz.Confirmar(v, alvo)))
	}
	p := matriz.PerfilDeUso{Etapas: etapas, Naturezas: naturezas, CasosDeUso: casos}
	return perfilResultado(d, "matriz.perfil_definido", alvo, errOf(matriz.Definir(v, alvo, p)))
}

func errOf(_ matriz.Entrada, err error) error { return err }

// perfilResultado imprime o resultado: termo fora do vocabulário é erro de uso (exit 2); o resto é falha (exit 1).
func perfilResultado(d Deps, msg, alvo string, err error) int {
	if errors.Is(err, matriz.ErrVocabulario) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	if err != nil {
		return fail(d, err)
	}
	fmt.Fprintln(d.Stdout, ste.Msg(msg, alvo))
	return 0
}

package forja

import "strings"

// PerfisValidos devolve a lista canônica de perfis de contexto da Forja.
func PerfisValidos() []string {
	return []string{"minimo", "padrao", "gateway"}
}

// ValidarPerfil verifica se perfil pertence ao conjunto minimo, padrao ou gateway.
func ValidarPerfil(perfil string) error {
	switch perfil {
	case "minimo", "padrao", "gateway":
		return nil
	default:
		return ErrPerfilDesconhecido
	}
}

// NormalizarSem converte a lista --sem em um conjunto case-insensitive.
func NormalizarSem(sem []string) map[string]bool {
	out := make(map[string]bool, len(sem))
	for _, item := range sem {
		chave := strings.ToLower(strings.TrimSpace(item))
		if chave != "" {
			out[chave] = true
		}
	}
	return out
}

// EstaExcluido informa se harness foi marcado para exclusão via --sem.
func EstaExcluido(harness string, excluidos map[string]bool) bool {
	if len(excluidos) == 0 {
		return false
	}
	chave := strings.ToLower(strings.TrimSpace(harness))
	return excluidos[chave]
}

// FiltrarHarnesses separa harnesses ativos e ignorados preservando a ordem.
func FiltrarHarnesses(harnesses, sem []string) (ativos, ignorados []string) {
	excluidos := NormalizarSem(sem)
	for _, h := range harnesses {
		if EstaExcluido(h, excluidos) {
			ignorados = append(ignorados, h)
			continue
		}
		ativos = append(ativos, h)
	}
	return ativos, ignorados
}

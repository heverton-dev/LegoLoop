package fabricar

import "errors"

// ErrStopGatePending indica parada obrigatória no Stop-Gate após a cascata de
// 3 camadas (Cofre -> SO -> GitHub) antes da materialização quando --yes não está ativo.
var ErrStopGatePending = errors.New("fabricar: parada obrigatória no Stop-Gate após curadoria em 3 camadas (Cofre -> SO -> GitHub)")

// CuradoriaStopGateSummary consolida o resumo do processo e as sugestões do GitHub
// apresentados para aprovação humana no Stop-Gate antes da materialização.
type CuradoriaStopGateSummary struct {
	ProcessoResumo   string   `json:"processo_resumo"`
	Camadas          []string `json:"camadas"`
	SugestoesGitHub  []string `json:"sugestoes_github,omitempty"`
	AwaitingApproval bool     `json:"awaiting_approval"`
}

// AvaliarStopGateCuradoria formaliza a parada obrigatória no Stop-Gate após a
// cascata de 3 camadas (Cofre -> SO -> GitHub) antes da materialização.
func AvaliarStopGateCuradoria(req Request, resumo string, sugestoesGitHub []string) (CuradoriaStopGateSummary, error) {
	summary := CuradoriaStopGateSummary{
		ProcessoResumo:  resumo,
		Camadas:         []string{"Cofre", "SO", "GitHub"},
		SugestoesGitHub: append([]string(nil), sugestoesGitHub...),
	}
	if req.Interactive {
		summary.AwaitingApproval = true
		return summary, ErrStopGatePending
	}
	return summary, nil
}

// EtapaCuradoriaFronteira representa uma etapa do subprocesso obrigatório de
// Curadoria e Composição de Fronteira do verbo `legoloop fabricar`.
type EtapaCuradoriaFronteira struct {
	Ordem     int
	ID        string
	Nome      string
	Descricao string
}

// SubprocessoCuradoriaFronteira retorna as 6 etapas obrigatórias do subprocesso
// de Curadoria e Composição de Fronteira para `legoloop fabricar`.
func SubprocessoCuradoriaFronteira() []EtapaCuradoriaFronteira {
	return []EtapaCuradoriaFronteira{
		{1, "entrevista_escopo", "Entrevista e Escopo", "Definir objetivo, critérios binários de aceite, restrições e nível de qualidade (0 a 10)."},
		{2, "inventario_local", "Inventário Local", "Varrer SO e Cofre: Matriz de Peças, skills, MCPs e processos instalados."},
		{3, "curadoria_externa", "Curadoria Externa", "Buscar em GitHub/Fontes quando houver lacuna, importar via legoloop import e descartar repositórios temporários após canonizar."},
		{4, "elevacao_canonizacao", "Elevação e Canonização", "Executar legoloop upgrade --diagnose, corrigir I/O Contracts e Negative Guardrails e canonizar por SHA-256."},
		{5, "composicao_dag", "Composição do DAG Completo", "Compor nós reais com context_in, context_out, routing e gates determinísticos; nunca 1 nó genérico."},
		{6, "projecao_padrao", "Projeção Padrão em loops/<slug>/", "Gerar os 4 arquivos padrão (<slug>.loop.json, pipeline.dag.json, triagem.json, memoria_execucao.json), Sexteto em docs/ e validar com legoloop validate."},
	}
}

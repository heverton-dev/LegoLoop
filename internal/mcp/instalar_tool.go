package mcp

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/forja"
	"github.com/heverton-dev/LegoLoop/internal/instalar"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

type respostaInstalarMCP struct {
	RelatorioJSON string          `json:"relatorio_json"`
	RelatorioMD   string          `json:"relatorio_md,omitempty"`
	Veredito      string          `json:"veredito"`
	DryRun        bool            `json:"dry_run"`
	Aplicado      bool            `json:"aplicado"`
	AtalhoGrafo   string          `json:"atalho_grafo,omitempty"`
	Relatorio     forja.Relatorio `json:"relatorio"`
}

func (r *Router) callInstalar(id any, args map[string]any) JSONRPCResponse {
	op, err := montarOpcaoInstalar(args)
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	res, err := instalar.Executar(context.Background(), op)
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	caminhoJSON, err := resolverCaminhoRelatorioJSON(op.Home, res)
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	return r.jsonResult(id, respostaInstalarMCP{
		RelatorioJSON: caminhoJSON,
		RelatorioMD:   res.CaminhoMD,
		Veredito:      res.Relatorio.Veredito,
		DryRun:        !op.Aplicar,
		Aplicado:      res.Aplicado,
		AtalhoGrafo:   res.AtalhoGrafo,
		Relatorio:     res.Relatorio,
	})
}

func montarOpcaoInstalar(args map[string]any) (instalar.Opcao, error) {
	home, err := resolverHomeInstalar(args)
	if err != nil {
		return instalar.Opcao{}, err
	}
	return instalar.Opcao{
		Home:    home,
		Perfil:  strArg(args, "perfil", "padrao"),
		Sem:     stringListArg(args, "sem"),
		Aplicar: resolverAplicarInstalar(args),
		Data:    strArg(args, "data", ""),
	}, nil
}

func resolverHomeInstalar(args map[string]any) (string, error) {
	if h := strings.TrimSpace(strArg(args, "home", "")); h != "" {
		return h, nil
	}
	return os.UserHomeDir()
}

func resolverAplicarInstalar(args map[string]any) bool {
	if dryRun, ok := args["dry_run"].(bool); ok {
		return !dryRun
	}
	if aplicar, ok := args["aplicar"].(bool); ok {
		return aplicar
	}
	return false
}

func resolverCaminhoRelatorioJSON(home string, res instalar.Resultado) (string, error) {
	if strings.TrimSpace(res.CaminhoJSON) != "" {
		return res.CaminhoJSON, nil
	}
	dirRel, err := jail.DentroDaRaiz(home, filepath.Join(".legoloop", "forja", res.Relatorio.Data))
	if err != nil {
		return "", err
	}
	return filepath.Join(dirRel, "RELATORIO.json"), nil
}

func instalarToolDef() map[string]any {
	return map[string]any{
		"name":        "instalar",
		"description": "Executa a instalacao do LegoLoop com a Forja em modo simulacao (dry_run) ou aplicacao, devolvendo o caminho do RELATORIO.json e o veredito.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"dry_run": map[string]any{"type": "boolean", "description": "Quando true (padrao), simula sem gravar arquivos no disco"},
				"home":    map[string]any{"type": "string", "description": "Diretorio HOME alvo (padrao: HOME do usuario)"},
				"perfil":  map[string]any{"type": "string", "description": "Perfil de contexto: minimo, padrao ou gateway"},
				"sem":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Harnesses excluidos da execucao"},
			},
		},
	}
}

package videorea

import "fmt"

// FindingClass diz se o LegoLoop DEVE, PODE ou NAO PODE integrar o achado.
type FindingClass string

const (
	ClassDeve    FindingClass = "DEVE"
	ClassPode    FindingClass = "PODE"
	ClassNaoPode FindingClass = "NAO_PODE"
)

// Finding classifica um achado do dossiê com o porquê e o como.
type Finding struct {
	Kind  string       `json:"kind"`
	Name  string       `json:"name"`
	Class FindingClass `json:"class"`
	Why   string       `json:"why"`
	How   string       `json:"how"`
}

// classifyFindings classifica todos os achados do dossiê em DEVE, PODE ou NAO_PODE.
func classifyFindings(d *ForensicsDossier, subs []SubstitutedTool, adopted []AdoptedTool) []Finding {
	var out []Finding
	covered := map[string]bool{}
	for _, s := range subs {
		covered[s.Original.Name] = true
		out = append(out, classifyPaidTool(s))
	}
	for _, a := range adopted {
		covered[a.Adopted.Name] = true
		out = append(out, classifyFreeTool(a.Adopted.Name, a.Diagnosis))
	}
	for _, name := range append(append([]string{}, d.ToolsIdentified...), d.HiddenTools...) {
		if !covered[name] {
			covered[name] = true
			out = append(out, classifyUnknownTool(name))
		}
	}
	out = append(out, classifyIdeas(d)...)
	return out
}

func classifyPaidTool(s SubstitutedTool) Finding {
	return Finding{
		Kind: "ferramenta", Name: s.Original.Name, Class: ClassNaoPode,
		Why: fmt.Sprintf("Ferramenta proprietária paga (%s) cria custo fixo e dependência de fornecedor.", s.Original.CostTier),
		How: "Integrar o substituto soberano: " + s.Replacement.Name + ".",
	}
}

func classifyFreeTool(name string, diag ToolDiagnosis) Finding {
	f := Finding{Kind: "ferramenta", Name: name, Class: ClassPode, How: diag.FitSolution}
	if diag.FitStatus == FitOuro {
		f.Class = ClassDeve
		f.Why = "Ferramenta gratuita com encaixe ouro no host do operador."
		return f
	}
	f.Why = "Ferramenta gratuita, mas o encaixe no host é " + string(diag.FitStatus) + "."
	if diag.InstallCommand != "" {
		f.How = diag.FitSolution + " Instalar com: " + diag.InstallCommand
	}
	return f
}

func classifyUnknownTool(name string) Finding {
	return Finding{
		Kind: "ferramenta", Name: name, Class: ClassPode,
		Why: "Custo e licença desconhecidos: nenhuma matriz cobre esta ferramenta.",
		How: "Verificar a licença e o custo antes de integrar; só então abrir o ticket de adoção.",
	}
}

func classifyIdeas(d *ForensicsDossier) []Finding {
	var out []Finding
	for _, c := range d.Concepts {
		out = append(out, Finding{Kind: "conceito", Name: c, Class: ClassPode,
			Why: "Ideia citada no vídeo, sem prova de encaixe no LegoLoop.",
			How: "Registrar no INTENT de um processo novo por entrevista socrática."})
	}
	for _, p := range d.Processes {
		out = append(out, Finding{Kind: "processo", Name: p, Class: ClassPode,
			Why: "Processo descrito no vídeo, ainda sem DAG nem gate.",
			How: "Fabricar o processo com legoloop fabricar e validar com legoloop validate."})
	}
	for _, r := range d.Repositories {
		out = append(out, Finding{Kind: "repositorio", Name: r, Class: ClassPode,
			Why: "Código externo ainda sem auditoria de licença e de dependências.",
			How: "Auditar a licença e as dependências; copiar só o método, pois o projeto usa só a stdlib."})
	}
	return out
}

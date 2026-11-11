package fabricar

import (
	"fmt"
	"strings"

	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

// Plano é tudo que a gravação vai fazer, mostrado ao usuário numa vez só, antes de qualquer efeito.
type Plano struct {
	Arquivos    []string // caminhos que serão gravados
	Vinculos    []string // links que serão criados nos harnesses
	Sobrescreve []string // arquivos que já existem e serão trocados (o original vira .legoloop-bak)
	Removivel   []string // peças de mídia removível ou de rede que entram na fusão
	Detectados  []string // nomes dos harnesses detectados
	Ligados     []string // nomes dos harnesses que vão receber o link
	Lacunas     []string // etapas sem peça apta: com elas nada é gravado (só o dry-run chega ao cartão)
}

// ApproveFunc mostra o plano e devolve true para gravar. replan recalcula o plano com outra lista de harnesses
// (o "personalizar"); a última lista aplicada com replan é a que vale na gravação.
type ApproveFunc func(p Plano, replan func(ligados []string) Plano) (bool, error)

// aprovar monta o plano, mostra ao usuário e devolve os vínculos que valem. Sem Aprovar ou no dry-run segue direto.
func aprovar(d Deps, req Request, files []PlannedFile, m Merged, name string) ([]processo.Vinculo, []string, error) {
	skill := artifactsOf(files).SkillPath
	vs, warns := planLinks(req, d.Harnesses, skill, name)
	if d.Aprovar == nil || req.DryRun {
		return vs, warns, nil
	}
	build := func(vs []processo.Vinculo) Plano { return novoPlano(d, req, files, m, vs) }
	replan := func(ligados []string) Plano {
		r2 := req
		r2.LinkHarnesses, r2.NoLink = ligados, len(ligados) == 0
		vs, _ = planLinks(r2, d.Harnesses, skill, name)
		return build(vs)
	}
	ok, err := d.Aprovar(build(vs), replan)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, ErrAborted
	}
	return vs, warns, nil
}

func novoPlano(d Deps, req Request, files []PlannedFile, m Merged, vs []processo.Vinculo) Plano {
	p := Plano{Arquivos: filePaths(files), Vinculos: linkDests(vs)}
	for _, f := range files {
		if _, err := d.FS.Stat(f.Path); err == nil {
			p.Sobrescreve = append(p.Sobrescreve, f.Path)
		}
	}
	for _, pc := range m.Used {
		if pc.Untrusted {
			p.Removivel = append(p.Removivel, pc.Path)
		}
	}
	for _, h := range d.Harnesses {
		p.Detectados = append(p.Detectados, h.Nome)
	}
	for _, v := range vs {
		p.Ligados = append(p.Ligados, v.Harness)
	}
	return p
}

// cartaoDaLacuna destaca as etapas sem peça: o aviso diz o que fazer, e não há efeito, porque nada é gravado.
func cartaoDaLacuna(etapas []string) ste.Cartao {
	c := ste.Cartao{
		Acao:        "Parar: o processo tem lacuna",
		Aviso:       msgLacuna(etapas),
		SeDerErrado: "Nada foi gravado. Não há o que desfazer. Corrija a lacuna e rode o comando de novo.",
	}
	for _, e := range etapas {
		c.Passos = append(c.Passos, ste.Msg("fabricar.lacuna_etapa", e))
	}
	c.Passos = append(c.Passos, "O programa não grava nenhum arquivo.")
	return c
}

// Cartao traduz o plano para o cartão STE-PT: o aviso de risco vem antes dos passos.
func (p Plano) Cartao() ste.Cartao {
	if len(p.Lacunas) > 0 {
		return cartaoDaLacuna(p.Lacunas)
	}
	c := ste.Cartao{
		Acao:        "Gravar o processo novo",
		Passos:      []string{fmt.Sprintf("O programa grava %d arquivos: skill, documento do Studio, DAG e os documentos.", len(p.Arquivos))},
		SeDerErrado: "Se um passo falhar, o programa desfaz tudo. Nenhum arquivo fica pela metade. Para remover depois, apague a pasta do processo no cofre e os links.",
	}
	for _, a := range p.Arquivos {
		c.Efeitos = append(c.Efeitos, "grava "+a)
	}
	if len(p.Vinculos) > 0 {
		c.Passos = append(c.Passos, fmt.Sprintf("O programa liga a skill a %d harnesses: %s.", len(p.Ligados), strings.Join(p.Ligados, ", ")))
	} else {
		c.Passos = append(c.Passos, "O programa não liga a skill a nenhum harness.")
	}
	for _, v := range p.Vinculos {
		c.Efeitos = append(c.Efeitos, "cria o link "+v)
	}
	if len(p.Sobrescreve) > 0 {
		c.Aviso = fmt.Sprintf("%d arquivos já existem. O programa troca cada um e guarda o original como .legoloop-bak.", len(p.Sobrescreve))
	}
	if len(p.Removivel) > 0 {
		c.Passos = append(c.Passos, fmt.Sprintf("O processo usa %d peças de mídia removível ou de rede.", len(p.Removivel)))
	}
	if len(p.Detectados) > 0 {
		c.Parametros = []ste.Parametro{{Nome: "harnesses", Valor: strings.Join(p.Ligados, ", "),
			Ajuda: "Nomes separados por vírgula. Detectados: " + strings.Join(p.Detectados, ", ") + ". Escreva nenhum para não ligar."}}
	}
	return c
}

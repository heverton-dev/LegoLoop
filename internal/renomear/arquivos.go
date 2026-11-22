package renomear

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/jail"
	"github.com/heverton-dev/LegoLoop/internal/processo"
)

func planejarPastaProjeto(op Opcoes) []Alteracao {
	dirProc := filepath.Join(op.Raiz, caminhos.DirProcessos)
	origem := filepath.Join(dirProc, op.Antigo)
	destino := filepath.Join(dirProc, op.Novo)
	if _, err := os.Stat(origem); err != nil {
		return nil
	}
	alts := []Alteracao{{Tipo: "pasta_processo", Origem: origem, Destino: destino}}
	alts = append(alts, planejarInternos(origem, destino, op.Antigo, op.Novo)...)
	solto := filepath.Join(dirProc, op.Antigo+".processo.json")
	if _, err := os.Stat(solto); err == nil {
		alts = append(alts, Alteracao{
			Tipo:    "manifesto_processo",
			Origem:  solto,
			Destino: filepath.Join(dirProc, op.Novo+".processo.json"),
		})
	}
	return alts
}

func planejarInternos(origem, destino, antigo, novo string) []Alteracao {
	var alts []Alteracao
	for _, arq := range []string{caminhos.ArquivoProcessoDAG, caminhos.ArquivoDAG, "triagem.json", "anticorpos.json", "ANTICORPOS.md"} {
		p := filepath.Join(origem, arq)
		if _, err := os.Stat(p); err == nil {
			alts = append(alts, Alteracao{
				Tipo:    "arquivo_interno",
				Origem:  p,
				Destino: filepath.Join(destino, arq),
				Detalhe: "atualizar " + arq,
			})
		}
	}
	studio := filepath.Join(origem, antigo+caminhos.SufixoStudio)
	if _, err := os.Stat(studio); err == nil {
		alts = append(alts, Alteracao{
			Tipo:    "studio_doc",
			Origem:  studio,
			Destino: filepath.Join(destino, novo+caminhos.SufixoStudio),
		})
	}
	return alts
}

func aplicarPastaProjeto(op Opcoes, bk *backupCtx) error {
	dirProc := filepath.Join(op.Raiz, caminhos.DirProcessos)
	origem := filepath.Join(dirProc, op.Antigo)
	if _, err := os.Stat(origem); err != nil {
		return nil
	}
	if err := bk.salvar(origem, "renomear_pasta"); err != nil {
		return err
	}
	if err := jail.Renomear(dirProc, op.Antigo, op.Novo); err != nil {
		return err
	}
	destino := filepath.Join(dirProc, op.Novo)
	if err := atualizarArquivosInternos(destino, op.Antigo, op.Novo); err != nil {
		return err
	}
	return renomearArquivoSolto(dirProc, op.Antigo+".processo.json", op.Novo+".processo.json", op.Antigo, op.Novo, bk)
}

func atualizarArquivosInternos(dir, antigo, novo string) error {
	for _, arq := range []string{caminhos.ArquivoProcessoDAG, caminhos.ArquivoDAG} {
		if err := atualizarChaveJSON(filepath.Join(dir, arq), "name", antigo, novo); err != nil {
			return err
		}
	}
	if err := atualizarTriagem(filepath.Join(dir, "triagem.json"), antigo, novo); err != nil {
		return err
	}
	if err := atualizarChaveJSON(filepath.Join(dir, "anticorpos.json"), "processo", antigo, novo); err != nil {
		return err
	}
	if err := atualizarTextoSeExiste(filepath.Join(dir, "ANTICORPOS.md"), antigo, novo); err != nil {
		return err
	}
	return renomearStudioInterno(dir, antigo, novo)
}

func renomearStudioInterno(dir, antigo, novo string) error {
	orig := filepath.Join(dir, antigo+caminhos.SufixoStudio)
	if _, err := os.Stat(orig); err != nil {
		return nil
	}
	dest := filepath.Join(dir, novo+caminhos.SufixoStudio)
	if err := os.Rename(orig, dest); err != nil {
		return err
	}
	return atualizarChaveJSON(dest, "name", antigo, novo)
}

func atualizarChaveJSON(caminho, chave, antigo, novo string) error {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var obj map[string]any
	if err := json.Unmarshal(dados, &obj); err != nil {
		return err
	}
	if val, ok := obj[chave].(string); ok && val == antigo {
		obj[chave] = novo
	}
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(caminho, append(out, '\n'), 0o644)
}

func atualizarTriagem(caminho, antigo, novo string) error {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var obj map[string]any
	if err := json.Unmarshal(dados, &obj); err != nil {
		return err
	}
	if val, ok := obj["processo"].(string); ok && val == antigo {
		obj["processo"] = novo
	}
	atualizarRegrasTriagem(obj, antigo, novo)
	atualizarAmbiguidadeTriagem(obj, antigo, novo)
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(caminho, append(out, '\n'), 0o644)
}

func atualizarRegrasTriagem(obj map[string]any, antigo, novo string) {
	regras, ok := obj["regras_deterministicas"].([]any)
	if !ok {
		return
	}
	for _, r := range regras {
		m, ok := r.(map[string]any)
		if !ok {
			continue
		}
		if v, ok := m["padrao_slug"].(string); ok && v == antigo {
			m["padrao_slug"] = novo
		}
		if cmds, ok := m["comandos_cli"].([]any); ok {
			for i, c := range cmds {
				if s, ok := c.(string); ok {
					cmds[i] = substituirComandoTriagem(s, antigo, novo)
				}
			}
		}
	}
}

func substituirComandoTriagem(cmd, antigo, novo string) string {
	if cmd == "legoloop "+antigo {
		return "legoloop " + novo
	}
	if cmd == "/legoloop-"+antigo || cmd == "/"+antigo {
		return "/" + novo
	}
	return cmd
}

func atualizarAmbiguidadeTriagem(obj map[string]any, antigo, novo string) {
	crit, ok := obj["criterio_ambiguidade"].(map[string]any)
	if !ok {
		return
	}
	if desc, ok := crit["descricao"].(string); ok {
		crit["descricao"] = substituirReferencias(desc, antigo, novo)
	}
}

func atualizarTextoSeExiste(caminho, antigo, novo string) error {
	dados, err := os.ReadFile(caminho)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	novoTexto := substituirReferencias(string(dados), antigo, novo)
	return os.WriteFile(caminho, []byte(novoTexto), 0o644)
}

func planejarCofreEVinculos(op Opcoes) []Alteracao {
	var alts []Alteracao
	if op.Vault != "" {
		for _, par := range paresCofre(op.Vault, op.Antigo, op.Novo) {
			if _, err := os.Lstat(par[0]); err == nil {
				alts = append(alts, Alteracao{Tipo: "cofre", Origem: par[0], Destino: par[1]})
			}
		}
	}
	if op.Home != "" {
		hs := processo.DetectarHarnesses(processo.OSFS(), op.Home)
		for _, h := range hs {
			orig := filepath.Join(h.Skills, op.Antigo)
			if _, err := os.Lstat(orig); err == nil {
				alts = append(alts, Alteracao{
					Tipo:    "vinculo_harness",
					Origem:  orig,
					Destino: filepath.Join(h.Skills, op.Novo),
					Detalhe: h.Nome,
				})
			}
		}
	}
	return alts
}

func paresCofre(vault, antigo, novo string) [][2]string {
	return [][2]string{
		{filepath.Join(vault, caminhos.DirProcessos, antigo), filepath.Join(vault, caminhos.DirProcessos, novo)},
		{filepath.Join(vault, "skills", antigo), filepath.Join(vault, "skills", novo)},
		{processo.CaminhoNoCofre(vault, antigo), processo.CaminhoNoCofre(vault, novo)},
	}
}

func aplicarCofreEVinculos(op Opcoes, bk *backupCtx) error {
	if op.Vault != "" {
		if err := aplicarCofre(op, bk); err != nil {
			return err
		}
	}
	if op.Home != "" {
		return aplicarVinculos(op, bk)
	}
	return nil
}

func aplicarCofre(op Opcoes, bk *backupCtx) error {
	for _, par := range paresCofre(op.Vault, op.Antigo, op.Novo) {
		info, err := os.Lstat(par[0])
		if err != nil {
			continue
		}
		if err := bk.salvar(par[0], "renomear_cofre"); err != nil {
			return err
		}
		if err := os.Rename(par[0], par[1]); err != nil {
			return err
		}
		if info.IsDir() {
			if err := atualizarArquivosInternos(par[1], op.Antigo, op.Novo); err != nil {
				return err
			}
		} else if strings.HasSuffix(par[1], ".json") {
			if err := atualizarChaveJSON(par[1], "name", op.Antigo, op.Novo); err != nil {
				return err
			}
		}
	}
	return nil
}

func aplicarVinculos(op Opcoes, bk *backupCtx) error {
	hs := processo.DetectarHarnesses(processo.OSFS(), op.Home)
	for _, h := range hs {
		orig := filepath.Join(h.Skills, op.Antigo)
		info, err := os.Lstat(orig)
		if err != nil {
			continue
		}
		if err := bk.salvar(orig, "renomear_vinculo"); err != nil {
			return err
		}
		dest := filepath.Join(h.Skills, op.Novo)
		if info.Mode()&os.ModeSymlink != 0 {
			if err := reapontarSymlink(orig, dest, op); err != nil {
				return err
			}
			continue
		}
		if err := os.Rename(orig, dest); err != nil {
			return err
		}
	}
	return nil
}

func reapontarSymlink(orig, dest string, op Opcoes) error {
	alvo, err := os.Readlink(orig)
	if err != nil {
		return err
	}
	novoAlvo := alvo
	if filepath.Base(alvo) == op.Antigo {
		novoAlvo = filepath.Join(filepath.Dir(alvo), op.Novo)
	}
	if err := os.Remove(orig); err != nil {
		return err
	}
	return os.Symlink(novoAlvo, dest)
}

func renomearArquivoSolto(dir, nomeAntigo, nomeNovo, antigo, novo string, bk *backupCtx) error {
	orig := filepath.Join(dir, nomeAntigo)
	if _, err := os.Stat(orig); err != nil {
		return nil
	}
	if err := bk.salvar(orig, "renomear_manifesto"); err != nil {
		return err
	}
	dest := filepath.Join(dir, nomeNovo)
	if err := os.Rename(orig, dest); err != nil {
		return err
	}
	return atualizarChaveJSON(dest, "name", antigo, novo)
}

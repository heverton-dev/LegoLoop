package forja

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/jail"
	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// Executar executa o motor da Forja nos modos simulação, aplicar ou desfazer.
func Executar(ctx context.Context, op Opcao) (Relatorio, error) {
	if err := ctx.Err(); err != nil {
		return Relatorio{}, err
	}
	if err := op.Validate(); err != nil {
		return Relatorio{}, err
	}
	op = preencherPadroes(op)
	rel := Relatorio{
		Versao:   op.Versao,
		Data:     op.Data,
		SO:       op.SO,
		Perfil:   op.Perfil,
		Veredito: VereditoOK,
	}
	if strings.TrimSpace(op.Desfazer) != "" {
		return executarDesfazer(ctx, op, rel)
	}
	return executarHarnesses(ctx, op, rel)
}

func preencherPadroes(op Opcao) Opcao {
	if strings.TrimSpace(op.Data) == "" {
		op.Data = time.Now().UTC().Format("2006-01-02")
	}
	if strings.TrimSpace(op.SO) == "" {
		op.SO = runtime.GOOS
	}
	if strings.TrimSpace(op.Versao) == "" {
		op.Versao = VersaoContrato
	}
	if op.Harnesses == nil {
		op.Harnesses = detectarNomesHarnesses(op.Home)
	}
	return op
}

func detectarNomesHarnesses(home string) []string {
	encontrados := processo.DetectarHarnesses(processo.OSFS(), home)
	nomes := make([]string, 0, len(encontrados))
	for _, h := range encontrados {
		nomes = append(nomes, h.Nome)
	}
	return nomes
}

func executarHarnesses(ctx context.Context, op Opcao, rel Relatorio) (Relatorio, error) {
	excluidos := NormalizarSem(op.Sem)
	for _, h := range op.Harnesses {
		if err := ctx.Err(); err != nil {
			rel.Veredito = VereditoFalha
			return rel, err
		}
		etapa, err := processarHarness(ctx, op, h, excluidos)
		rel.Etapas = append(rel.Etapas, etapa)
		if err != nil {
			rel.Veredito = VereditoFalha
			return rel, err
		}
	}
	return rel, nil
}

func processarHarness(ctx context.Context, op Opcao, h string, excluidos map[string]bool) (EtapaResultado, error) {
	if EstaExcluido(h, excluidos) {
		return EtapaResultado{
			Nome:     "forjar:" + h,
			Harness:  h,
			Estado:   EstadoIgnorada,
			Mensagem: "excluído por --sem",
		}, nil
	}
	pendentes, err := planejarPendentes(ctx, op, h)
	if err != nil {
		return EtapaResultado{
			Nome:     "forjar:" + h,
			Harness:  h,
			Estado:   EstadoFalha,
			Mensagem: err.Error(),
		}, err
	}
	if !op.Aplicar {
		return EtapaResultado{
			Nome:      "forjar:" + h,
			Harness:   h,
			Estado:    EstadoPendente,
			Alterados: caminhosRelativos(op.Home, pendentes),
			Mensagem:  "simulação",
		}, nil
	}
	return aplicarEtapa(op, h, pendentes)
}

func planejarPendentes(ctx context.Context, op Opcao, h string) ([]AlteracaoPlanejada, error) {
	alts, err := gerarAlteracoes(ctx, op, h)
	if err != nil {
		return nil, err
	}
	var pendentes []AlteracaoPlanejada
	for _, alt := range alts {
		seguro, err := jail.DentroDaRaiz(op.Home, alt.Caminho)
		if err != nil {
			return nil, err
		}
		alt.Caminho = seguro
		if alt.Modo == 0 {
			alt.Modo = 0o644
		}
		if arquivoDiferente(alt.Caminho, alt.Conteudo) {
			pendentes = append(pendentes, alt)
		}
	}
	return pendentes, nil
}

func gerarAlteracoes(ctx context.Context, op Opcao, h string) ([]AlteracaoPlanejada, error) {
	if op.Geradores != nil {
		if gen, ok := op.Geradores[h]; ok && gen != nil {
			return gen(ctx, op, h)
		}
	}
	if err := jail.ValidarNomeProcesso(h); err != nil {
		return nil, err
	}
	destino, err := jail.DentroDaRaiz(op.Home, filepath.Join("."+h, "legoloop_forja.cfg"))
	if err != nil {
		return nil, err
	}
	conteudo := fmt.Sprintf("harness=%s\nperfil=%s\nso=%s\n", h, op.Perfil, op.SO)
	return []AlteracaoPlanejada{{Caminho: destino, Conteudo: []byte(conteudo), Modo: 0o644}}, nil
}

func arquivoDiferente(caminho string, alvo []byte) bool {
	atual, err := os.ReadFile(caminho)
	if err != nil {
		return true
	}
	return !bytes.Equal(atual, alvo)
}

func caminhosRelativos(home string, alts []AlteracaoPlanejada) []string {
	if len(alts) == 0 {
		return nil
	}
	out := make([]string, 0, len(alts))
	for _, alt := range alts {
		rel, err := filepath.Rel(home, alt.Caminho)
		if err != nil {
			out = append(out, filepath.ToSlash(alt.Caminho))
			continue
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out
}

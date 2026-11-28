package roteiro

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/jail"
	"github.com/heverton-dev/LegoLoop/internal/redc"
)

// PerfilEstiloDNA guarda os parametros estruturais extraidos da referencia sem plagio visual.
type PerfilEstiloDNA struct {
	SchemaVersion   int               `json:"schema_version"`
	CanalReferencia string            `json:"canal_referencia"`
	TomDeVoz        string            `json:"tom_de_voz"`
	Arquetipo       string            `json:"arquetipo"`
	RitmoCortes     string            `json:"ritmo_cortes"`
	MotionTokens    map[string]string `json:"motion_tokens,omitempty"`
	PadroesGancho   []string          `json:"padroes_gancho,omitempty"`
}

// AplicarChannelDNA aplica Brand Swap sobre o DNA do canal referencia, produzindo Marca propria e PerfilEstiloDNA.
func AplicarChannelDNA(dna redc.ChannelDNA, tokens redc.BrandTokens, base Marca) (Marca, PerfilEstiloDNA, error) {
	if strings.TrimSpace(dna.CanalSlug) == "" {
		return Marca{}, PerfilEstiloDNA{}, fmt.Errorf("%w: canal_slug obrigatorio no ChannelDNA", ErrSchemaInvalido)
	}
	if strings.TrimSpace(tokens.NomeMarca) == "" {
		return Marca{}, PerfilEstiloDNA{}, fmt.Errorf("%w: nome_marca obrigatorio em BrandTokens", ErrSchemaInvalido)
	}
	m := normalizarMarcaBase(base, tokens, dna.CanalSlug)
	p := construirPerfilDNA(dna, tokens.TomPreferido)
	return m, p, nil
}

func normalizarMarcaBase(base Marca, t redc.BrandTokens, slugRef string) Marca {
	m := base
	if m.SchemaVersion <= 0 {
		m.SchemaVersion = 1
	}
	if len(m.Formatos) == 0 {
		m.Formatos = []string{"16:9", "9:16", "1:1"}
	}
	if m.Voz.Voz == "" {
		m.Voz = VozMarca{Motor: "edge-tts", Voz: "pt-BR-ThalitaNeural", Velocidade: "-5%"}
	}
	m.Nome = strings.TrimSpace(t.NomeMarca)
	m.CanalReferencia = slugRef
	m.Cores = coresDeTokens(t)
	m.Fontes = fontesDeTokens(t)
	m.Logo = logoPadraoMarca(m.Logo)
	m.Vinhetas = vinhetasPadraoMarca(m.Nome, m.Vinhetas)
	if m.Legendas == nil {
		m.Legendas = &LegendasMarca{AtivoEm: []string{}, Fonte: "corpo", Cor: "#FFFFFF", Destaque: "foco", PalavrasPorBloco: 3, PosicaoYPct: 68, ArquivoSRT: true}
	}
	return m
}

func coresDeTokens(t redc.BrandTokens) CoresMarca {
	pri := valorOuPadrao(t.CorPrimaria, "#FFD24A")
	sec := valorOuPadrao(t.CorSecundaria, "#3DDC84")
	return CoresMarca{Foco: pri, Foco2: sec, Alerta: "#FF5A5A", Texto: "#FFFFFF", Contorno: "#111111", FundoCartao: "#000000"}
}

func fontesDeTokens(t redc.BrandTokens) FontesMarca {
	tit := valorOuPadrao(t.FonteTitulo, "Bebas Neue")
	cor := valorOuPadrao(t.FonteCorpo, "Inter")
	return FontesMarca{
		Titulo:     FonteMarca{Familia: tit, Peso: 400},
		Corpo:      FonteMarca{Familia: cor, Peso: 600},
		Numeros:    FonteMarca{Familia: "Anton", Peso: 400},
		Manuscrita: FonteMarca{Familia: "Caveat", Peso: 700},
	}
}

func logoPadraoMarca(atual LogoMarca) LogoMarca {
	if atual.Arquivo != "" {
		return atual
	}
	return LogoMarca{Arquivo: "marca/logo.png", Posicao: "canto_superior_direito", Opacidade: 0.85, LarguraPct: 8}
}

func vinhetasPadraoMarca(nome string, atual VinhetasMarca) VinhetasMarca {
	if atual.Abertura.Texto != "" {
		return atual
	}
	return VinhetasMarca{
		Abertura:   VinhetaCfg{Ativo: true, DuracaoS: 1.5, Texto: strings.ToUpper(nome)},
		Fechamento: VinhetaCfg{Ativo: true, DuracaoS: 3.0, Texto: "Inscreva-se"},
	}
}

func construirPerfilDNA(dna redc.ChannelDNA, tomPreferido string) PerfilEstiloDNA {
	tom := strings.TrimSpace(tomPreferido)
	if tom == "" {
		tom = dna.TomDeVoz
	}
	mt := map[string]string{}
	maps.Copy(mt, dna.MotionTokens)
	return PerfilEstiloDNA{
		SchemaVersion:   1,
		CanalReferencia: dna.CanalSlug,
		TomDeVoz:        tom,
		Arquetipo:       dna.Arquetipo,
		RitmoCortes:     dna.RitmoCortes,
		MotionTokens:    mt,
		PadroesGancho:   append([]string(nil), dna.TopGanchos...),
	}
}

func valorOuPadrao(v, padrao string) string {
	if strings.TrimSpace(v) == "" {
		return padrao
	}
	return strings.TrimSpace(v)
}

// GravarBrandSwapDNA aplica o ChannelDNA e persiste marca.json e perfil_estilo.json de forma atomica.
func GravarBrandSwapDNA(projDir string, dna redc.ChannelDNA, tokens redc.BrandTokens, base Marca) (Marca, PerfilEstiloDNA, error) {
	m, p, err := AplicarChannelDNA(dna, tokens, base)
	if err != nil {
		return Marca{}, PerfilEstiloDNA{}, err
	}
	if err := gravarJSONAtomicoRoteiro(projDir, "marca.json", m); err != nil {
		return Marca{}, PerfilEstiloDNA{}, err
	}
	if err := gravarJSONAtomicoRoteiro(projDir, "perfil_estilo.json", p); err != nil {
		return Marca{}, PerfilEstiloDNA{}, err
	}
	return m, p, nil
}

func gravarJSONAtomicoRoteiro(projDir, rel string, v any) error {
	destino, err := jail.DentroDaRaiz(projDir, rel)
	if err != nil {
		return err
	}
	dados, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := destino + ".tmp"
	if err := os.WriteFile(tmp, append(dados, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, destino)
}

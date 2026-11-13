package combo

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrComboInvalido          = errors.New("combo: nome de combo desconhecido")
	ErrTetoTokensExcedido     = errors.New("combo: teto de tokens excedido para o combo")
	ErrFerramentaNaoSuportada = errors.New("combo: modelo não suporta chamadas de função")
)

// MembroCombo define um modelo individual participante de um combo.
type MembroCombo struct {
	ID                string  `json:"id"`
	Provedor          string  `json:"provedor"`
	LatenciaMedia100k float64 `json:"latencia_media_100k_s"`
	SuportaFuncao     bool    `json:"suporta_funcao"`
	SuportaImagem     bool    `json:"suporta_imagem"`
}

// ComboHonesto define a especificação completa de um combo certificado.
type ComboHonesto struct {
	Nome             string        `json:"nome"`
	Finalidade       string        `json:"finalidade"`
	MaxTokens        int           `json:"max_tokens"`
	ExigeFerramentas bool          `json:"exige_ferramentas"`
	ExigeImagem      bool          `json:"exige_imagem"`
	Membros          []MembroCombo `json:"membros"`
}

// Validate verifica a integridade estrutural do combo.
func (c ComboHonesto) Validate() error {
	if strings.TrimSpace(c.Nome) == "" || len(c.Membros) == 0 || c.MaxTokens <= 0 {
		return ErrComboInvalido
	}
	return nil
}

// CatalogoCombos armazena a definicao canonica dos 6 combos honestos.
var CatalogoCombos = map[string]ComboHonesto{
	"t-agente-128k": {
		Nome:             "t-agente-128k",
		Finalidade:       "Agentes autonomos com suporte a ferramentas e chamadas de funcao",
		MaxTokens:        131072,
		ExigeFerramentas: true,
		ExigeImagem:      false,
		Membros: []MembroCombo{
			{ID: "mistral/ministral-3b-latest", Provedor: "mistral", LatenciaMedia100k: 2.4, SuportaFuncao: true, SuportaImagem: true},
			{ID: "gemini/gemini-3.5-flash-lite", Provedor: "gemini", LatenciaMedia100k: 2.8, SuportaFuncao: true, SuportaImagem: false},
			{ID: "nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", Provedor: "nvidia", LatenciaMedia100k: 3.1, SuportaFuncao: true, SuportaImagem: true},
			{ID: "openrouter/dots-studio/dots-3-note-preview:free", Provedor: "openrouter", LatenciaMedia100k: 5.2, SuportaFuncao: true, SuportaImagem: true},
			{ID: "cf/@cf/meta/llama-4-scout-17b-16e-instruct", Provedor: "cloudflare", LatenciaMedia100k: 5.9, SuportaFuncao: true, SuportaImagem: false},
			{ID: "mistral/ministral-3b-2512", Provedor: "mistral", LatenciaMedia100k: 7.1, SuportaFuncao: true, SuportaImagem: true},
		},
	},
	"t-chat-128k": {
		Nome:             "t-chat-128k",
		Finalidade:       "Chat geral de alto contexto sem exigencia de ferramentas",
		MaxTokens:        131072,
		ExigeFerramentas: false,
		ExigeImagem:      false,
		Membros: []MembroCombo{
			{ID: "mistral/codestral-latest", Provedor: "mistral", LatenciaMedia100k: 2.1, SuportaFuncao: false, SuportaImagem: false},
			{ID: "gemini/gemini-3.5-flash-lite", Provedor: "gemini", LatenciaMedia100k: 2.8, SuportaFuncao: true, SuportaImagem: false},
			{ID: "nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", Provedor: "nvidia", LatenciaMedia100k: 3.1, SuportaFuncao: true, SuportaImagem: true},
			{ID: "openrouter/cohere/north-mini-code:free", Provedor: "openrouter", LatenciaMedia100k: 4.9, SuportaFuncao: false, SuportaImagem: false},
			{ID: "cf/@cf/meta/llama-4-scout-17b-16e-instruct", Provedor: "cloudflare", LatenciaMedia100k: 5.9, SuportaFuncao: true, SuportaImagem: false},
			{ID: "mistral/ministral-3b-latest", Provedor: "mistral", LatenciaMedia100k: 2.4, SuportaFuncao: true, SuportaImagem: true},
		},
	},
	"t-raciocinio-128k": {
		Nome:             "t-raciocinio-128k",
		Finalidade:       "Raciocinio complexo e inferencia logica de alto contexto",
		MaxTokens:        131072,
		ExigeFerramentas: false,
		ExigeImagem:      false,
		Membros: []MembroCombo{
			{ID: "gemini/gemini-3.5-flash-lite", Provedor: "gemini", LatenciaMedia100k: 2.8, SuportaFuncao: true, SuportaImagem: false},
			{ID: "nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", Provedor: "nvidia", LatenciaMedia100k: 3.1, SuportaFuncao: true, SuportaImagem: true},
			{ID: "openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free", Provedor: "openrouter", LatenciaMedia100k: 6.3, SuportaFuncao: true, SuportaImagem: true},
			{ID: "cf/@cf/openai/gpt-oss-20b", Provedor: "cloudflare", LatenciaMedia100k: 9.6, SuportaFuncao: false, SuportaImagem: false},
			{ID: "gemini/gemini-3-flash-preview", Provedor: "gemini", LatenciaMedia100k: 4.1, SuportaFuncao: true, SuportaImagem: false},
			{ID: "cf/@cf/openai/gpt-oss-120b", Provedor: "cloudflare", LatenciaMedia100k: 9.8, SuportaFuncao: true, SuportaImagem: false},
		},
	},
	"t-visao-128k": {
		Nome:             "t-visao-128k",
		Finalidade:       "Processamento multimodal e analise visual",
		MaxTokens:        131072,
		ExigeFerramentas: false,
		ExigeImagem:      true,
		Membros: []MembroCombo{
			{ID: "mistral/ministral-3b-latest", Provedor: "mistral", LatenciaMedia100k: 2.4, SuportaFuncao: true, SuportaImagem: true},
			{ID: "nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", Provedor: "nvidia", LatenciaMedia100k: 3.1, SuportaFuncao: true, SuportaImagem: true},
			{ID: "openrouter/dots-studio/dots-3-note-preview:free", Provedor: "openrouter", LatenciaMedia100k: 5.2, SuportaFuncao: true, SuportaImagem: true},
			{ID: "mistral/ministral-3b-2512", Provedor: "mistral", LatenciaMedia100k: 7.1, SuportaFuncao: true, SuportaImagem: true},
			{ID: "nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b", Provedor: "nvidia", LatenciaMedia100k: 3.5, SuportaFuncao: true, SuportaImagem: true},
			{ID: "openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free", Provedor: "openrouter", LatenciaMedia100k: 6.3, SuportaFuncao: true, SuportaImagem: true},
		},
	},
	"t-codigo-128k": {
		Nome:             "t-codigo-128k",
		Finalidade:       "Geracao e edicao de codigo sem ferramentas",
		MaxTokens:        131072,
		ExigeFerramentas: false,
		ExigeImagem:      false,
		Membros: []MembroCombo{
			{ID: "mistral/codestral-latest", Provedor: "mistral", LatenciaMedia100k: 2.1, SuportaFuncao: false, SuportaImagem: false},
			{ID: "openrouter/cohere/north-mini-code:free", Provedor: "openrouter", LatenciaMedia100k: 4.9, SuportaFuncao: false, SuportaImagem: false},
			{ID: "mistral/codestral-2508", Provedor: "mistral", LatenciaMedia100k: 5.9, SuportaFuncao: false, SuportaImagem: false},
			{ID: "mistral/mistral-code-latest", Provedor: "mistral", LatenciaMedia100k: 0.0, SuportaFuncao: false, SuportaImagem: false},
		},
	},
	"t-micro-8k": {
		Nome:             "t-micro-8k",
		Finalidade:       "Requisicoes rapidas de baixo contexto",
		MaxTokens:        8192,
		ExigeFerramentas: false,
		ExigeImagem:      false,
		Membros: []MembroCombo{
			{ID: "groq/openai/gpt-oss-120b", Provedor: "groq", LatenciaMedia100k: 1.6, SuportaFuncao: false, SuportaImagem: false},
			{ID: "groq/openai/gpt-oss-20b", Provedor: "groq", LatenciaMedia100k: 2.0, SuportaFuncao: false, SuportaImagem: false},
		},
	},
}

// CatalogoHonesto retorna uma copia dos combos canônicos certificados.
func CatalogoHonesto() map[string]ComboHonesto {
	copia := make(map[string]ComboHonesto, len(CatalogoCombos))
	for k, v := range CatalogoCombos {
		copia[k] = v
	}
	return copia
}

// ObterCombo busca um combo pelo nome no catalogo certificado.
func ObterCombo(nome string) (ComboHonesto, error) {
	c, ok := CatalogoCombos[nome]
	if !ok {
		return ComboHonesto{}, ErrComboInvalido
	}
	return c, nil
}

// ListarCombos retorna os nomes ordenados de todos os combos do catalogo.
func ListarCombos() []string {
	nomes := make([]string, 0, len(CatalogoCombos))
	for k := range CatalogoCombos {
		nomes = append(nomes, k)
	}
	sort.Strings(nomes)
	return nomes
}

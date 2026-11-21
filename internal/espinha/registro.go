package espinha

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// SchemaRegistro define a versao canonica do schema de espinha/registro.json.
const SchemaRegistro = "espinha.v1"

// Registro concentra os verbos, motores, processos e requisitos da espinha dorsal.
type Registro struct {
	Schema              string            `json:"schema"`
	Verbos              []Verbo           `json:"verbos"`
	Motores             []Motor           `json:"motores"`
	Processos           []Processo        `json:"processos"`
	Requisitos          []Requisito       `json:"requisitos"`
	LimitesErroIgnorado map[string]int    `json:"limites_erro_ignorado,omitempty"`
	GlobaisJustificadas []string          `json:"globais_justificadas,omitempty"`
	InitsJustificados   map[string]string `json:"inits_justificados,omitempty"`
}

// Verbo descreve um comando da CLI e sua camada de pertencimento.
type Verbo struct {
	Nome     string `json:"nome"`
	Camada   string `json:"camada"`
	Motor    string `json:"motor"`
	Contrato string `json:"contrato"`
}

// Motor descreve um pacote real em internal/ e se sua execucao e deterministica.
type Motor struct {
	Pacote         string `json:"pacote"`
	Interface      string `json:"interface"`
	Deterministico bool   `json:"deterministico"`
}

// Processo descreve um slug real existente em processos/.
type Processo struct {
	Slug string `json:"slug"`
}

// Requisito vincula uma lei do AGENTS.md ao seu verificador mecanico ou justificativa.
type Requisito struct {
	Lei           string `json:"lei"`
	Verificador   string `json:"verificador,omitempty"`
	Mecanico      bool   `json:"mecanico"`
	Justificativa string `json:"justificativa,omitempty"`
}

// CarregarRegistro le e decodifica o arquivo JSON do registro da espinha dorsal.
func CarregarRegistro(caminho string) (Registro, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return Registro{}, fmt.Errorf("ler %s: %w", caminho, err)
	}
	var reg Registro
	if err := json.Unmarshal(dados, &reg); err != nil {
		return Registro{}, fmt.Errorf("decodificar %s: %w", caminho, err)
	}
	return reg, nil
}

// Validar verifica schema, preenchimento, unicidade, camadas e requisitos.
func (r Registro) Validar() error {
	if r.Schema != SchemaRegistro {
		return fmt.Errorf("schema invalido: %q (esperado %q)", r.Schema, SchemaRegistro)
	}
	if err := validarVerbos(r.Verbos); err != nil {
		return err
	}
	if err := validarMotores(r.Motores); err != nil {
		return err
	}
	if err := validarProcessos(r.Processos); err != nil {
		return err
	}
	return validarRequisitos(r.Requisitos)
}

func validarVerbos(verbos []Verbo) error {
	if len(verbos) == 0 {
		return errors.New("lista verbos nao pode ser vazia")
	}
	vistos := make(map[string]struct{}, len(verbos))
	for _, v := range verbos {
		nome := strings.TrimSpace(v.Nome)
		if nome == "" {
			return errors.New("verbo com nome vazio")
		}
		if _, dup := vistos[nome]; dup {
			return fmt.Errorf("verbo repetido: %q", nome)
		}
		vistos[nome] = struct{}{}
		if v.Camada != "espinha" && v.Camada != "dominio" {
			return fmt.Errorf("camada invalida no verbo %q: %q", nome, v.Camada)
		}
	}
	return nil
}

func validarMotores(motores []Motor) error {
	if len(motores) == 0 {
		return errors.New("lista motores nao pode ser vazia")
	}
	vistos := make(map[string]struct{}, len(motores))
	for _, m := range motores {
		pacote := strings.TrimSpace(m.Pacote)
		if pacote == "" {
			return errors.New("motor com pacote vazio")
		}
		if _, dup := vistos[pacote]; dup {
			return fmt.Errorf("motor repetido: %q", pacote)
		}
		vistos[pacote] = struct{}{}
	}
	return nil
}

func validarProcessos(processos []Processo) error {
	if len(processos) == 0 {
		return errors.New("lista processos nao pode ser vazia")
	}
	vistos := make(map[string]struct{}, len(processos))
	for _, p := range processos {
		slug := strings.TrimSpace(p.Slug)
		if slug == "" {
			return errors.New("processo com slug vazio")
		}
		if _, dup := vistos[slug]; dup {
			return fmt.Errorf("processo repetido: %q", slug)
		}
		vistos[slug] = struct{}{}
	}
	return nil
}

func validarRequisitos(requisitos []Requisito) error {
	if len(requisitos) == 0 {
		return errors.New("lista requisitos nao pode ser vazia")
	}
	vistos := make(map[string]struct{}, len(requisitos))
	for _, req := range requisitos {
		lei := strings.TrimSpace(req.Lei)
		if lei == "" {
			return errors.New("requisito com lei vazia")
		}
		if _, dup := vistos[lei]; dup {
			return fmt.Errorf("requisito repetido: %q", lei)
		}
		vistos[lei] = struct{}{}
		if req.Mecanico && strings.TrimSpace(req.Verificador) == "" {
			return fmt.Errorf("requisito mecanico %q exige verificador", lei)
		}
		if !req.Mecanico && strings.TrimSpace(req.Justificativa) == "" {
			return fmt.Errorf("requisito nao mecanico %q exige justificativa", lei)
		}
	}
	return nil
}

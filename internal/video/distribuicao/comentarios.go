package distribuicao

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Estados de uma resposta a comentario.
const (
	StatusRascunho = "rascunho"
	StatusAprovada = "aprovada"
	StatusIgnorada = "ignorada"
)

// Comentario e um comentario coletado do YouTube ou do Instagram.
type Comentario struct {
	ID         string `json:"id"`
	Rede       string `json:"rede"`
	Autor      string `json:"autor"`
	Texto      string `json:"texto"`
	Respondido bool   `json:"respondido"`
}

// Resposta e o rascunho de resposta; so a aprovada pode sair.
type Resposta struct {
	ComentarioID string `json:"comentario_id"`
	Rede         string `json:"rede"`
	Rascunho     string `json:"rascunho"`
	Status       string `json:"status"`
}

// Respostas e o arquivo acervo/respostas.json.
type Respostas struct {
	SchemaVersion int        `json:"schema_version"`
	Itens         []Resposta `json:"itens"`
}

// CaminhoRespostas devolve acervo/respostas.json.
func CaminhoRespostas(projDir string) string {
	return filepath.Join(projDir, "acervo", "respostas.json")
}

// CarregarRespostas le o arquivo; ausente devolve vazio.
func CarregarRespostas(projDir string) (Respostas, error) {
	r := Respostas{SchemaVersion: 1}
	dados, err := os.ReadFile(CaminhoRespostas(projDir))
	if errors.Is(err, os.ErrNotExist) {
		return r, nil
	}
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(dados, &r); err != nil {
		return Respostas{}, fmt.Errorf("%w: respostas.json: %v", ErrAcervoInvalido, err)
	}
	return r, nil
}

// Gravar escreve o arquivo de forma atomica.
func (r Respostas) Gravar(projDir string) error {
	dados, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return gravarAtomico(CaminhoRespostas(projDir), append(dados, '\n'))
}

// CarregarComentarios le a lista coletada (JSON) de um arquivo.
func CarregarComentarios(caminho string) ([]Comentario, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	var cs []Comentario
	if err := json.Unmarshal(dados, &cs); err != nil {
		return nil, fmt.Errorf("%w: comentarios: %v", ErrPedidoInvalido, err)
	}
	return cs, nil
}

// Rascunhar grava um rascunho para um comentario ainda sem resposta.
func (r *Respostas) Rascunhar(c Comentario, texto string) error {
	switch {
	case strings.TrimSpace(c.ID) == "" || strings.TrimSpace(texto) == "":
		return fmt.Errorf("%w: comentario e rascunho sao obrigatorios", ErrPedidoInvalido)
	case c.Respondido:
		return fmt.Errorf("%w: comentario %q ja foi respondido", ErrPedidoInvalido, c.ID)
	case r.buscar(c.ID) != nil:
		return fmt.Errorf("%w: comentario %q ja tem resposta", ErrPedidoInvalido, c.ID)
	}
	r.Itens = append(r.Itens, Resposta{ComentarioID: c.ID, Rede: c.Rede, Rascunho: strings.TrimSpace(texto), Status: StatusRascunho})
	return nil
}

func (r *Respostas) buscar(id string) *Resposta {
	for i := range r.Itens {
		if r.Itens[i].ComentarioID == id {
			return &r.Itens[i]
		}
	}
	return nil
}

// Decidir aprova (com texto editado opcional) ou ignora um rascunho. Decisao tomada nao muda.
func (r *Respostas) Decidir(id string, aprovar bool, editado string) error {
	x := r.buscar(id)
	if x == nil {
		return fmt.Errorf("%w: comentario %q sem rascunho", ErrPedidoInvalido, id)
	}
	if x.Status != StatusRascunho {
		return fmt.Errorf("%w: comentario %q ja esta %s", ErrPedidoInvalido, id, x.Status)
	}
	if !aprovar {
		x.Status = StatusIgnorada
		return nil
	}
	if e := strings.TrimSpace(editado); e != "" {
		x.Rascunho = e
	}
	x.Status = StatusAprovada
	return nil
}

// Aprovadas devolve so as respostas que o usuario aprovou; nada sai sem isso.
func (r Respostas) Aprovadas() []Resposta {
	var out []Resposta
	for _, x := range r.Itens {
		if x.Status == StatusAprovada {
			out = append(out, x)
		}
	}
	return out
}

// Pendentes devolve os comentarios sem resposta e sem rascunho.
func Pendentes(cs []Comentario, r Respostas) []Comentario {
	var out []Comentario
	for _, c := range cs {
		if !c.Respondido && r.buscar(c.ID) == nil {
			out = append(out, c)
		}
	}
	return out
}

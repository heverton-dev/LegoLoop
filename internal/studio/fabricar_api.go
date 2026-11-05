package studio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/fabricar"
	"legoloop/internal/specgen"
)

// fabricarTimeout é o teto de uma fabricação (varredura, LLM e gravação).
const fabricarTimeout = 10 * time.Minute

// FabricarFunc executa a fábrica de processos. É o ponto de injeção do Studio; nil usa o 9Router do inventário.
type FabricarFunc func(ctx context.Context, req fabricar.Request) (*fabricar.Result, error)

// fabricarRunRequest é tudo que o cliente pode pedir. Caminhos, confiança em mídia removível, link em harness
// e drives são decisões do operador (servidor) e não existem aqui: campo desconhecido vira 400.
type fabricarRunRequest struct {
	Name       string `json:"name"`
	Goal       string `json:"goal"`
	Category   string `json:"category"`
	Model      string `json:"model"`
	Force      bool   `json:"force"`
	Entrevista string `json:"entrevista,omitempty"`
}

// handleFabricarRun materializa um processo novo (skill, loop, dag e documentos). Rota distinta de
// /api/canonize/*, que segue sendo a desduplicação.
func (s *Server) handleFabricarRun(w http.ResponseWriter, r *http.Request) {
	req, rawReq, ok := s.decodeFabricarRequest(w, r)
	if !ok {
		return
	}
	if err := req.Validate(); err != nil {
		writeFabricarError(w, err)
		return
	}
	if s.Discover != nil && !s.inventory(r.Context()).Has(req.ModelChoice) {
		writeError(w, http.StatusBadRequest, "modelo fora do inventário")
		return
	}
	if !s.fabricarMu.TryLock() {
		writeErrorCode(w, http.StatusConflict, "busy", "já existe uma fabricação em andamento")
		return
	}
	defer s.fabricarMu.Unlock()
	// A fabricação termina mesmo se a aba fechar, mas tem teto; e o WriteTimeout de 30 s do servidor
	// cortaria a resposta de uma execução que leva minutos, então o prazo de escrita é estendido.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), fabricarTimeout)
	defer cancel()
	_ = http.NewResponseController(w).SetWriteDeadline(time.Now().Add(fabricarTimeout))
	res, err := s.runFabricar(ctx, req, rawReq.Entrevista)
	if err != nil {
		writeFabricarError(w, err)
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(ctx) // a skill nova aparece no catálogo sem reiniciar o servidor
	}
	writeJSON(w, http.StatusOK, s.maskResult(res))
}

func (s *Server) decodeFabricarRequest(w http.ResponseWriter, r *http.Request) (fabricar.Request, fabricarRunRequest, bool) {
	var in fabricarRunRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&in)
	var tooBig *http.MaxBytesError
	switch {
	case errors.As(err, &tooBig):
		writeError(w, http.StatusRequestEntityTooLarge, "corpo excede 1 MiB")
		return fabricar.Request{}, in, false
	case err != nil, !atEOF(dec):
		writeError(w, http.StatusBadRequest, "pedido inválido: JSON malformado ou com campos desconhecidos")
		return fabricar.Request{}, in, false
	}
	if in.Category == "" {
		in.Category = string(fabricar.CategoryBuild)
	}
	req := fabricar.Request{RawName: in.Name, Goal: in.Goal, Category: fabricar.Category(in.Category), ModelChoice: in.Model, Force: in.Force}
	return req, in, true
}

// runFabricar usa a função injetada ou, na falta dela, a fábrica real com o cérebro escolhido no 9Router.
func (s *Server) runFabricar(ctx context.Context, req fabricar.Request, entrevista string) (*fabricar.Result, error) {
	if s.Fabricar != nil {
		return s.Fabricar(ctx, req)
	}
	return s.defaultFabricar(ctx, req, entrevista)
}

func (s *Server) defaultFabricar(ctx context.Context, req fabricar.Request, entrevista string) (*fabricar.Result, error) {
	if s.Discover == nil || s.ProjectDir == "" {
		return nil, errors.New("studio: fábrica indisponível neste servidor")
	}
	b, _, err := brain.Pick(s.inventory(ctx), req.ModelChoice, s.BrainOpts...)
	if err != nil {
		return nil, err
	}
	vault, err := s.vault()
	if err != nil {
		return nil, err
	}
	deps := fabricar.Deps{LLM: b, Matriz: s.FabricarMatriz, Roots: fabricar.Roots{Vault: vault, Docs: s.ProjectDir}}
	if entrevista != "" && s.Interview != nil {
		if sess, err := s.Interview.GetSession(entrevista); err == nil && sess != nil {
			clar := fabricar.ClarificationFromSession(sess)
			deps.Interview = &clar
		}
	}
	return fabricar.Fabricar(ctx, req, deps)
}

// writeFabricarError traduz o erro em status HTTP com mensagem própria. O texto bruto do erro nunca sobe ao
// navegador: ele pode conter caminhos locais.
func writeFabricarError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, fabricar.ErrInvalidName):
		writeError(w, http.StatusBadRequest, "nome inválido: use letras e números")
	case errors.Is(err, fabricar.ErrInvalidRequest):
		writeError(w, http.StatusBadRequest, "pedido inválido: informe uma categoria válida e um modelo")
	case errors.Is(err, fabricar.ErrExists):
		writeErrorCode(w, http.StatusConflict, "exists", "já existe um processo com este nome")
	case errors.Is(err, brain.ErrModeloDesconhecido):
		writeError(w, http.StatusUnprocessableEntity, "o modelo escolhido não está disponível; escolha outro da lista")
	case errors.Is(err, brain.ErrSemCerebro):
		writeError(w, http.StatusServiceUnavailable, "nenhum modelo está disponível; suba o 9Router, instale o Ollama ou instale um harness com CLI")
	case errors.Is(err, specgen.ErrContract):
		writeError(w, http.StatusBadGateway, "o modelo respondeu fora do contrato após várias tentativas; tente outro modelo")
	default:
		writeError(w, http.StatusInternalServerError, "falha ao fabricar o processo")
	}
}

// writeErrorCode acrescenta um código estável ao erro, para a interface decidir o que oferecer (por exemplo,
// sobrescrever só quando o 409 é de duplicata, e nunca quando é de fabricação em andamento).
func writeErrorCode(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]string{"error": msg, "code": code})
}

// maskResult devolve uma cópia do Result sem caminhos absolutos: projeto vira <projeto> e o HOME vira ~.
func (s *Server) maskResult(res *fabricar.Result) fabricar.Result {
	out := *res
	a := &out.Artifacts
	for _, p := range []*string{&a.IntentPath, &a.SpecPath, &a.PlanPath, &a.RoutePath, &a.TemplatePath, &a.SkillPath, &a.LoopPath, &a.DagPath} {
		*p = s.displayPath(*p)
	}
	a.SkillLinks = nil
	out.CuratedPieces = append([]fabricar.DiscoveredPiece(nil), res.CuratedPieces...)
	for i := range out.CuratedPieces {
		out.CuratedPieces[i].Path = s.displayPath(out.CuratedPieces[i].Path)
	}
	out.Warnings = make([]string, len(res.Warnings))
	for i, w := range res.Warnings {
		out.Warnings[i] = s.maskPaths(w)
	}
	if res.Report != nil {
		repCopy := *res.Report
		var maskedFiles []string
		for _, f := range repCopy.Gravacao.Arquivos {
			maskedFiles = append(maskedFiles, s.displayPath(f))
		}
		repCopy.Gravacao.Arquivos = maskedFiles
		var maskedEfeitos []string
		for _, e := range repCopy.Cartao.Efeitos {
			maskedEfeitos = append(maskedEfeitos, s.maskPaths(e))
		}
		repCopy.Cartao.Efeitos = maskedEfeitos
		out.Report = &repCopy
	}
	return out
}

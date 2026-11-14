package studiodoc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("criar pasta do acervo: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) filePath(id string) (string, error) {
	if err := ValidateID(id); err != nil {
		return "", err
	}
	novo := filepath.Join(s.dir, id+caminhos.SufixoStudio)
	if info, err := os.Stat(novo); err == nil && info.Mode().IsRegular() {
		return novo, nil
	}
	aninhado := filepath.Join(s.dir, id, id+caminhos.SufixoStudio)
	if info, err := os.Stat(aninhado); err == nil && info.Mode().IsRegular() {
		return aninhado, nil
	}
	return filepath.Join(s.dir, id+sufixoAntigo), nil
}

const sufixoAntigo = ".lo" + "op.json"

func validateProcessoDAG(l *Processo) error {
	v := struct {
		ID          string       `json:"id,omitempty"`
		Name        string       `json:"name"`
		Description string       `json:"description,omitempty"`
		Version     string       `json:"version,omitempty"`
		Trigger     *dag.Trigger `json:"trigger,omitempty"`
		Nodes       []dag.Node   `json:"nodes"`
	}{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		Version:     l.Version,
		Trigger:     l.Trigger,
		Nodes:       l.Nodes,
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("serializar Processo: %w", err)
	}
	if _, err := dag.Parse(raw); err != nil {
		return fmt.Errorf("DAG inválido: %w", err)
	}
	return nil
}

func (s *Store) Save(l *Processo) error {
	if err := ValidateID(l.ID); err != nil {
		return err
	}
	p, err := s.filePath(l.ID)
	if err != nil {
		return err
	}
	if err := validateProcessoDAG(l); err != nil {
		return err
	}
	l.UpdatedAt = time.Now().UTC()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = l.UpdatedAt
	}
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(p, data)
}

func writeAtomic(dest string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", dest, time.Now().UnixNano())
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func (s *Store) Get(id string) (*Processo, error) {
	p, err := s.filePath(id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		nested := filepath.Join(s.dir, id, filepath.Base(p))
		data, err = os.ReadFile(nested)
	}
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var lp Processo
	if err := json.Unmarshal(data, &lp); err != nil {
		return nil, fmt.Errorf("decodificar Processo: %w", err)
	}
	return &lp, nil
}

func (s *Store) List() ([]Summary, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var list []Summary
	vistos := map[string]bool{}
	for _, e := range entries {
		id := idDeArquivo(e)
		if id == "" || vistos[id] {
			continue
		}
		vistos[id] = true
		lp, err := s.Get(id)
		if err != nil {
			continue
		}
		list = append(list, lp.Summary())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list, nil
}

func idDeArquivo(e os.DirEntry) string {
	if e.IsDir() {
		if strings.HasPrefix(e.Name(), ".") {
			return ""
		}
		return e.Name()
	}
	if strings.HasSuffix(e.Name(), caminhos.SufixoStudio) {
		return strings.TrimSuffix(e.Name(), caminhos.SufixoStudio)
	}
	if strings.HasSuffix(e.Name(), sufixoAntigo) {
		return strings.TrimSuffix(e.Name(), sufixoAntigo)
	}
	return ""
}

func (l *Processo) Summary() Summary {
	finalArt := ""
	if len(l.Nodes) > 0 {
		finalArt = l.Nodes[len(l.Nodes)-1].Output.Artifact
	}
	trgType := "manual"
	if l.Trigger != nil && l.Trigger.Type != "" {
		trgType = string(l.Trigger.Type)
	}
	return Summary{
		ID:            l.ID,
		Name:          l.Name,
		Description:   l.Description,
		NodesCount:    len(l.Nodes),
		FinalArtifact: finalArt,
		TriggerType:   trgType,
		UpdatedAt:     l.UpdatedAt,
	}
}

func (s *Store) Duplicate(id, newID, newName string) (*Processo, error) {
	orig, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	dup := *orig
	dup.ID = newID
	dup.Name = newName
	dup.CreatedAt = time.Time{}
	dup.UpdatedAt = time.Time{}
	if err := s.Save(&dup); err != nil {
		return nil, err
	}
	return &dup, nil
}

func (s *Store) Remove(id string) error {
	p, err := s.filePath(id)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

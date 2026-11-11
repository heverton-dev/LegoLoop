package brain

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Antibody captura a tripla imunológica de um erro resolvido (Motor III).
type Antibody struct {
	Package        string `json:"package"`
	ErrorSignature string `json:"error_signature"`
	RootCause      string `json:"root_cause"`
	DiffSolution   string `json:"diff_solution"`
	Timestamp      string `json:"timestamp"`
}

// AntibodiesDB armazena e consulta anticorpos imunológicos para injeção preventiva.
type AntibodiesDB struct {
	mu         sync.RWMutex
	antibodies map[string][]Antibody
}

// NewAntibodiesDB instancia um banco de anticorpos em memória.
func NewAntibodiesDB() *AntibodiesDB {
	return &AntibodiesDB{
		antibodies: make(map[string][]Antibody),
	}
}

// Record grava uma nova tripla imunológica associada a um pacote.
func (db *AntibodiesDB) Record(pkg, errSig, rootCause, diffSol string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	ab := Antibody{
		Package:        pkg,
		ErrorSignature: errSig,
		RootCause:      rootCause,
		DiffSolution:   diffSol,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
	}

	db.antibodies[pkg] = append(db.antibodies[pkg], ab)
	return nil
}

// Find localiza um anticorpo pela assinatura do erro dentro de um pacote.
func (db *AntibodiesDB) Find(pkg, errSig string) (*Antibody, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, ab := range db.antibodies[pkg] {
		if strings.EqualFold(ab.ErrorSignature, errSig) {
			return &ab, true
		}
	}
	return nil, false
}

// InjectPreventiveContext gera um texto de advertência preventiva com os anticorpos do pacote.
func (db *AntibodiesDB) InjectPreventiveContext(pkg string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	list := db.antibodies[pkg]
	if len(list) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n🛡️ [ANTICORPO PREVENTIVO // MEMÓRIA IMUNOLÓGICA]:\n")
	for _, ab := range list {
		sb.WriteString(fmt.Sprintf("· Causa Raiz Conhecida: %s\n· Solução Canônica: %s\n", ab.RootCause, ab.DiffSolution))
	}
	return sb.String()
}

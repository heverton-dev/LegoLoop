package routedag

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

// VerificarCadeia aplica a trava da cadeia antes de projetar o DAG: a ROUTE assina a INTENT, a SPEC e o PLAN,
// os hashes conferem com o disco e a INTENT está aprovada. base é a raiz do projeto.
func VerificarCadeia(raw []byte, base string) error {
	var doc struct {
		Upstream plangen.UpstreamSignatures `json:"upstream_signatures"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("ROUTE inválida: %w", err)
	}
	if doc.Upstream.Intent == nil {
		return upstream.ErrSemIntencao
	}
	if err := plangen.VerifySignatures(base, doc.Upstream); err != nil {
		return err
	}
	return upstream.IntencaoAprovada(filepath.Join(base, filepath.FromSlash(doc.Upstream.Intent.Path)))
}

// raizDaRoute devolve a raiz do projeto a partir do caminho da ROUTE (raiz/docs/routes/arquivo).
func raizDaRoute(caminho string) (string, error) {
	abs, err := filepath.Abs(caminho)
	if err != nil {
		return "", err
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(abs))), nil
}

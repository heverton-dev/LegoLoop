package routedag

import "os"

// Projetar converte uma rota em nós de DAG e valida a integridade.
func Projetar(caminhoRoute string) ([]byte, string, error) {
	raw, err := os.ReadFile(caminhoRoute)
	if err != nil {
		return nil, "", err
	}
	return Converter(raw)
}

package orchestrix

import "errors"

// ErrMemoriaIndisponivel indica um SO sem leitura de memória livre.
var ErrMemoriaIndisponivel = errors.New("leitura de memória livre indisponível neste sistema")

// MemoriaLivreSO lê a memória física disponível, em bytes, pela implementação do SO.
func MemoriaLivreSO() (uint64, error) { return memoriaLivreSO() }

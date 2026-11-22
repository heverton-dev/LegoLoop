//go:build !windows && !linux && !darwin

package orchestrix

func memoriaLivreSO() (uint64, error) { return 0, ErrMemoriaIndisponivel }

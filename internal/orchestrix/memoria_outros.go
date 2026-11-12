//go:build !windows && !linux

package orchestrix

func memoriaLivreSO() (uint64, error) { return 0, ErrMemoriaIndisponivel }

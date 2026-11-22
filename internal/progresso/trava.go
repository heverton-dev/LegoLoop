package progresso

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	esperaMaximaTrava  = 10 * time.Second
	idadeTravaObsoleta = 30 * time.Second
	passoEsperaTrava   = 20 * time.Millisecond
)

// Atualizar carrega o progresso, aplica fn e grava o JSON e o espelho MD sob trava exclusiva.
// Vários processos (nós paralelos, swarm) podem chamar Atualizar ao mesmo tempo sem perder escrita.
// Se fn devolver erro, nada é gravado.
func Atualizar(caminhoJSON string, fn func(*Progresso) error) (err error) {
	libera, err := adquirirTrava(caminhoJSON + `.lock`)
	if err != nil {
		return err
	}
	defer func() {
		if errLibera := libera(); err == nil {
			err = errLibera
		}
	}()
	p, err := Carregar(caminhoJSON)
	if err != nil {
		return err
	}
	if err = fn(&p); err != nil {
		return err
	}
	return Salvar(caminhoJSON, CaminhoEspelhoMD(caminhoJSON), p)
}

func adquirirTrava(caminho string) (func() error, error) {
	limite := time.Now().Add(esperaMaximaTrava)
	for {
		f, err := os.OpenFile(caminho, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			return registrarDono(f, caminho)
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf(`trava de progresso %s: %w`, caminho, err)
		}
		if err := removerTravaObsoleta(caminho); err != nil {
			return nil, err
		}
		if time.Now().After(limite) {
			return nil, fmt.Errorf(`trava de progresso ocupada: %s`, caminho)
		}
		time.Sleep(passoEsperaTrava)
	}
}

// registrarDono grava o PID na trava recém-criada e devolve a função que a libera.
func registrarDono(f *os.File, caminho string) (func() error, error) {
	_, errEscrita := fmt.Fprintf(f, `%d`, os.Getpid())
	if err := errors.Join(errEscrita, f.Close()); err != nil {
		return nil, errors.Join(fmt.Errorf(`trava de progresso %s: %w`, caminho, err), liberarTrava(caminho))
	}
	return func() error { return liberarTrava(caminho) }, nil
}

func liberarTrava(caminho string) error {
	if err := os.Remove(caminho); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(`liberar trava %s: %w`, caminho, err)
	}
	return nil
}

// removerTravaObsoleta apaga a trava deixada por um processo que caiu.
func removerTravaObsoleta(caminho string) error {
	info, err := os.Stat(caminho)
	if err != nil || time.Since(info.ModTime()) <= idadeTravaObsoleta {
		return nil
	}
	return liberarTrava(caminho)
}

package interview

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// ErrSessaoInvalida: o identificador não tem o formato intv_<número>; recusar evita sair da pasta de sessões.
	ErrSessaoInvalida = errors.New("interview: identificador de sessão inválido")
	// ErrSessaoNaoEncontrada: não há arquivo para o identificador.
	ErrSessaoNaoEncontrada = errors.New("interview: sessão não encontrada")

	reSessionID = regexp.MustCompile(`^intv_[0-9]{1,20}$`)
)

// Store guarda cada sessão em Dir/<id>.json. É o que deixa uma entrevista continuar em outro processo,
// como no Chat: cada resposta do usuário vira uma chamada nova do comando.
type Store struct{ Dir string }

func (s Store) path(id string) (string, error) {
	if !reSessionID.MatchString(id) {
		return "", fmt.Errorf("%w: %q", ErrSessaoInvalida, id)
	}
	return filepath.Join(s.Dir, id+".json"), nil
}

// Save grava a sessão por inteiro. Escreve num arquivo temporário e renomeia, para nunca deixar JSON pela metade.
func (s Store) Save(sess *Session) error {
	path, err := s.path(sess.ID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.Dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(s.Dir, ".sessao-*.tmp")
	if err != nil {
		return err
	}
	_, werr := tmp.Write(data)
	cerr := tmp.Close()
	if err := errors.Join(werr, cerr); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	return nil
}

// Load lê a sessão pelo identificador.
func (s Store) Load(id string) (*Session, error) {
	path, err := s.path(id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrSessaoNaoEncontrada, id)
	}
	if err != nil {
		return nil, err
	}
	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("interview: sessão %s corrompida: %w", id, err)
	}
	return &sess, nil
}

// Adopt registra no motor uma sessão que veio do disco, para que Step possa continuá-la.
func (e *Engine) Adopt(sess *Session) {
	e.mu.Lock()
	e.sessions[sess.ID] = sess
	e.mu.Unlock()
}

package triagem

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestConferirCoberturaMultiplasEtapas(t *testing.T) {
	dir := t.TempDir()
	spec := GerarTriagemPadrao("demo", []string{"a", "b"})
	dados, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ArquivoTriagem), dados, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a", "b", "c"}); err == nil {
		t.Fatal("esperava erro quando falta nó c")
	}
	specUm := GerarTriagemPadrao("demo", []string{"a"})
	dadosUm, _ := json.Marshal(specUm)
	if err := os.WriteFile(filepath.Join(dir, ArquivoTriagem), dadosUm, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Conferir(dir, []string{"a"}); err != nil {
		t.Fatalf("esperava sucesso com 1 etapa, obteve %v", err)
	}
}

func TestMotorCasaRegraEtapaCamposVazios(t *testing.T) {
	spec := TriagemSpec{
		SchemaVersion: SchemaVersionV1,
		Processo:      "demo",
		RegrasDeterministicas: []RegraDeterministica{
			{PadraoSlug: "demo", Prioridade: 1},
		},
		Etapas: []RegraEtapa{
			{ID: "e1", PalavrasChave: []string{"alvo"}},
		},
		PontosEntrada: []string{"cli"},
	}
	_, err := Decidir(context.Background(), Entrada{Texto: "texto sem match"}, []TriagemSpec{spec}, nil)
	if !errors.Is(err, ErrResolvedorAusente) {
		t.Fatalf("esperava ErrResolvedorAusente, obteve %v", err)
	}
}

func TestTriagemLimparNomeEtapaETemDigito(t *testing.T) {
	spec := GerarTriagemPadrao("p", []string{"01_etapa-composta", "unica", "01", "sem-digito"})
	if !reflect.DeepEqual(spec.Etapas[0].PalavrasChave, []string{"etapa composta", "etapa-composta"}) {
		t.Fatalf("palavras etapa composta inesperadas: %v", spec.Etapas[0].PalavrasChave)
	}
	if !reflect.DeepEqual(spec.Etapas[1].PalavrasChave, []string{"unica"}) {
		t.Fatalf("palavras etapa simples inesperadas: %v", spec.Etapas[1].PalavrasChave)
	}
	if limparNomeEtapa("01") != "01" {
		t.Fatalf("esperava '01', obteve %q", limparNomeEtapa("01"))
	}
	if limparNomeEtapa("sem-digito") != "sem-digito" {
		t.Fatalf("esperava 'sem-digito', obteve %q", limparNomeEtapa("sem-digito"))
	}
	if !temDigito("0") || !temDigito("9") || temDigito("abc") || temDigito("") {
		t.Fatal("temDigito falhou nos limites '0', '9', 'abc' e vazio")
	}
}

func TestPidVivoEPidDoLockLimites(t *testing.T) {
	if pidVivo(-1) || pidVivo(0) {
		t.Fatal("pidVivo deve ser falso para PID <= 0")
	}
	if pidDoLock([]byte("12345 s1")) != 12345 {
		t.Fatalf("esperava 12345, obteve %d", pidDoLock([]byte("12345 s1")))
	}
	if pidDoLock([]byte("")) != 0 || pidDoLock([]byte("abc s1")) != 0 {
		t.Fatal("pidDoLock deve retornar 0 para entrada vazia ou não numérica")
	}
}

func TestMemoriaChaveSessaoEDesserializar(t *testing.T) {
	if _, err := chaveSessao(MemoriaExecucao{SessaoID: "s1", PID: 0}); !errors.Is(err, ErrSessaoInvalida) {
		t.Fatalf("esperava ErrSessaoInvalida para PID 0, obteve %v", err)
	}
	if _, err := chaveSessao(MemoriaExecucao{SessaoID: "", PID: 10}); !errors.Is(err, ErrSessaoInvalida) {
		t.Fatalf("esperava ErrSessaoInvalida para SessaoID vazio, obteve %v", err)
	}
	if err := SalvarMemoria(t.TempDir(), MemoriaExecucao{SessaoID: "s1", PID: 0}); !errors.Is(err, ErrSessaoInvalida) {
		t.Fatalf("SalvarMemoria deveria falhar com ErrSessaoInvalida, obteve %v", err)
	}
	if _, err := desserializarEstado([]byte("{invalido")); err == nil {
		t.Fatal("esperava erro ao desserializar JSON inválido")
	}
	if _, err := desserializarEstado([]byte(`{"sessoes":"invalido"}`)); err == nil {
		t.Fatal("esperava erro ao desserializar envelope com sessoes inválido")
	}
	if _, err := desserializarEstado([]byte(`{"pid":"texto"}`)); err == nil {
		t.Fatal("esperava erro ao desserializar legado com pid inválido")
	}
}

func TestMemoriaErrosDiscoEBackup(t *testing.T) {
	dir := t.TempDir()
	arqMem := filepath.Join(dir, ArquivoMemoria)
	if err := os.WriteFile(arqMem, []byte("{corrompido"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := lerEstadoDisco(dir); err != ErrMemoriaCorrompida {
		t.Fatalf("esperava ErrMemoriaCorrompida exato, obteve %v", err)
	}
	dir2 := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir2, ArquivoMemoria), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := lerEstadoDisco(dir2); err == nil || errors.Is(err, ErrMemoriaCorrompida) {
		t.Fatalf("esperava erro de leitura (não ErrMemoriaCorrompida) ao ler diretório, obteve %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir2, ArquivoMemoria, "x"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := LimparMemoria(dir2); err == nil {
		t.Fatal("esperava erro ao remover diretório não vazio")
	}
}

func TestMemoriaEscritaAtomicaEPersistirErro(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, ArquivoMemoria)
	if err := os.Mkdir(caminho+".tmp", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := escreverAtomico(caminho, []byte("{}")); err == nil {
		t.Fatal("esperava erro quando .tmp é diretório")
	}
	est := EstadoMemoria{
		Sessoes: map[string]MemoriaExecucao{
			"s1": {AtualizadoEm: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	if err := persistirEstado(t.TempDir(), est); err == nil {
		t.Fatal("esperava erro ao serializar ano 10000")
	}
}

func TestLockLiberarNilEErrosAdquirirELimpar(t *testing.T) {
	origTimeout := timeoutLockMemoria
	timeoutLockMemoria = 25 * time.Millisecond
	t.Cleanup(func() { timeoutLockMemoria = origTimeout })
	var nulo *lockMemoria
	nulo.liberar()
	(&lockMemoria{}).liberar()
	arqPai := filepath.Join(t.TempDir(), "arquivo")
	if err := os.WriteFile(arqPai, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirInvalido := filepath.Join(arqPai, "sub")
	_, err := adquirirLockMemoria(dirInvalido, os.Getpid(), "s1")
	var pe *os.PathError
	if err == nil || !errors.As(err, &pe) || (pe.Op != "mkdir" && pe.Op != "open") {
		t.Fatalf("esperava PathError mkdir/open em adquirirLockMemoria, obteve %v", err)
	}
	if _, err := LerMemoria(dirInvalido, "s1"); err == nil || !errors.As(err, &pe) || (pe.Op != "mkdir" && pe.Op != "open") {
		t.Fatalf("esperava PathError mkdir/open em LerMemoria, obteve %v", err)
	}
	if err := ZerarSessao(dirInvalido, "s1"); err == nil || !errors.As(err, &pe) || (pe.Op != "mkdir" && pe.Op != "open") {
		t.Fatalf("esperava PathError mkdir/open em ZerarSessao, obteve %v", err)
	}
	if _, err := LimparOrfaos(dirInvalido); err == nil {
		t.Fatal("esperava erro em LimparOrfaos com dir inválido")
	}
}

func TestZerarELimparOrfaosErrosLockEDisco(t *testing.T) {
	origTimeout := timeoutLockMemoria
	timeoutLockMemoria = 25 * time.Millisecond
	t.Cleanup(func() { timeoutLockMemoria = origTimeout })
	dirLockDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dirLockDir, ArquivoLockMemoria), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := limparLockOrfao(dirLockDir); err == nil {
		t.Fatal("esperava erro em limparLockOrfao quando lock é diretório")
	}
	var pe *os.PathError
	if err := ZerarSessao(dirLockDir, "s1"); err == nil || !errors.As(err, &pe) || (pe.Op != "read" && pe.Op != "open") {
		t.Fatalf("esperava PathError read/open em ZerarSessao, obteve %v", err)
	}
	if _, err := LimparOrfaos(dirLockDir); err == nil || !errors.As(err, &pe) || (pe.Op != "read" && pe.Op != "open") {
		t.Fatalf("esperava PathError read/open em LimparOrfaos, obteve %v", err)
	}
	dirMemDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dirMemDir, ArquivoMemoria), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ZerarSessao(dirMemDir, "s1"); err == nil || strings.Contains(err.Error(), "inexistente") {
		t.Fatalf("esperava erro de leitura em ZerarSessao, obteve %v", err)
	}
	if _, err := LimparOrfaos(dirMemDir); err == nil {
		t.Fatal("esperava erro em LimparOrfaos quando memoria é diretório")
	}
}

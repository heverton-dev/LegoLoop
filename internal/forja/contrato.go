package forja

import (
	"context"
	"errors"
	"io/fs"
	"strings"
)

const (
	// VersaoContrato identifica a versão do contrato de saída da Forja.
	VersaoContrato = "1.0.0"

	EstadoOK       = "ok"
	EstadoFalha    = "falha"
	EstadoIgnorada = "ignorada"
	EstadoPendente = "pendente"

	VereditoOK    = "OK"
	VereditoFalha = "FALHA"
)

var (
	// ErrPerfilDesconhecido sinaliza perfil fora de minimo, padrao e gateway.
	ErrPerfilDesconhecido = errors.New("forja: perfil desconhecido")
	// ErrSemHome sinaliza HOME vazio.
	ErrSemHome = errors.New("forja: home vazio")
	// ErrBackupInexistente sinaliza ausência do diretório de backup solicitado.
	ErrBackupInexistente = errors.New("forja: backup inexistente")
	// ErrManifestoInvalido sinaliza formato corrompido em MANIFESTO.tsv.
	ErrManifestoInvalido = errors.New("forja: manifesto inválido")
)

// AlteracaoPlanejada representa um arquivo que a Forja pretende gravar.
type AlteracaoPlanejada struct {
	Caminho  string
	Conteudo []byte
	Modo     fs.FileMode
}

// GeradorFunc planeja as alterações de um harness sem escrever em disco.
type GeradorFunc func(ctx context.Context, op Opcao, harness string) ([]AlteracaoPlanejada, error)

// Opcao configura uma execução da Forja.
type Opcao struct {
	Home      string                 `json:"home"`
	Perfil    string                 `json:"perfil"`             // minimo | padrao | gateway
	Sem       []string               `json:"sem"`                // harnesses excluídos
	Aplicar   bool                   `json:"aplicar"`            // false = simula
	Desfazer  string                 `json:"desfazer,omitempty"` // data do backup a restaurar
	Data      string                 `json:"-"`                  // data injetada (AAAA-MM-DD)
	SO        string                 `json:"-"`                  // SO injetado (runtime.GOOS)
	Versao    string                 `json:"-"`                  // versão da Forja
	Harnesses []string               `json:"-"`                  // lista injetada de harnesses
	Geradores map[string]GeradorFunc `json:"-"`                  // geradores opcionais por harness
}

// EtapaResultado descreve o que uma Etapa da Forja fez ou faria.
type EtapaResultado struct {
	Nome      string   `json:"nome"`
	Harness   string   `json:"harness,omitempty"`
	Estado    string   `json:"estado"` // ok | falha | ignorada | pendente
	Alterados []string `json:"alterados,omitempty"`
	Backup    string   `json:"backup,omitempty"`
	Mensagem  string   `json:"mensagem,omitempty"`
}

// Relatorio é o artefato de comprovação: os quatro gatilhos geram o mesmo conteúdo.
type Relatorio struct {
	Versao       string           `json:"versao"`
	Data         string           `json:"data"`
	SO           string           `json:"so"`
	Perfil       string           `json:"perfil"`
	Veredito     string           `json:"veredito"` // OK | FALHA
	Etapas       []EtapaResultado `json:"etapas"`
	TokensAntes  map[string]int   `json:"tokens_antes,omitempty"`
	TokensDepois map[string]int   `json:"tokens_depois,omitempty"`
}

// Validate checa a Opcao antes de qualquer E/S.
func (o Opcao) Validate() error {
	if strings.TrimSpace(o.Home) == "" {
		return ErrSemHome
	}
	return ValidarPerfil(o.Perfil)
}

// Package vault gerencia credenciais e segredos de forma segura no Cofre de Segredos com criptografia
// AES-256-GCM, mascaramento estrito em logs e injeção em nós de execução.
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Vault gerencia pares chave-valor criptografados do Cofre de Segredos em memória e persistidos em disco.
type Vault struct {
	mu          sync.RWMutex
	storagePath string
	key         []byte
	secrets     map[string]string
}

// New cria uma nova instância do Cofre de Segredos. Se storagePath for vazio, opera apenas em memória.
func New(storagePath string, masterKey string) (*Vault, error) {
	if masterKey == "" {
		return nil, errors.New("vault: chave mestra não pode ser vazia")
	}

	derivedKey := deriveKey(masterKey)

	return &Vault{
		storagePath: storagePath,
		key:         derivedKey,
		secrets:     make(map[string]string),
	}, nil
}

// deriveKey gera uma chave criptográfica simétrica de 256 bits (32 bytes) a partir da semente.
func deriveKey(masterKey string) []byte {
	hash := sha256.Sum256([]byte(masterKey))
	return hash[:]
}

// Set adiciona ou atualiza um segredo no Cofre de Segredos.
func (v *Vault) Set(key, value string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if key == "" {
		return errors.New("vault: chave não pode ser vazia")
	}
	v.secrets[key] = value
	return nil
}

// Get recupera o valor de um segredo pelo identificador.
func (v *Vault) Get(key string) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	val, ok := v.secrets[key]
	return val, ok
}

// ListKeys lista todos os identificadores de segredos armazenados (sem revelar os valores).
func (v *Vault) ListKeys() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	keys := make([]string, 0, len(v.secrets))
	for k := range v.secrets {
		keys = append(keys, k)
	}
	return keys
}

// Save serializa e criptografa os segredos em disco usando AES-256-GCM.
func (v *Vault) Save() error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.storagePath == "" {
		return nil
	}

	payload, err := json.Marshal(v.secrets)
	if err != nil {
		return fmt.Errorf("vault: falha na serialização json: %w", err)
	}

	encrypted, err := encrypt(payload, v.key)
	if err != nil {
		return fmt.Errorf("vault: falha na criptografia: %w", err)
	}

	dir := filepath.Dir(v.storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("vault: falha ao criar diretório do storage: %w", err)
	}

	if err := os.WriteFile(v.storagePath, encrypted, 0600); err != nil {
		return fmt.Errorf("vault: falha ao gravar arquivo seguro: %w", err)
	}

	return nil
}

// Load lê e decriptografa os segredos salvos em disco.
func (v *Vault) Load() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.storagePath == "" {
		return errors.New("vault: storagePath não configurado")
	}

	raw, err := os.ReadFile(v.storagePath)
	if err != nil {
		return fmt.Errorf("vault: falha ao ler arquivo: %w", err)
	}

	decrypted, err := decrypt(raw, v.key)
	if err != nil {
		return fmt.Errorf("vault: falha na decriptografia (chave inválida ou arquivo corrompido): %w", err)
	}

	var data map[string]string
	if err := json.Unmarshal(decrypted, &data); err != nil {
		return fmt.Errorf("vault: dados decriptografados inválidos: %w", err)
	}

	v.secrets = data
	return nil
}

// MaskSecrets substitui valores sensíveis presentes no texto fornecido por ***MASKED***.
func (v *Vault) MaskSecrets(raw string) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	result := raw
	for _, val := range v.secrets {
		if len(val) >= 4 { // Apenas valores com comprimento mínimo razoável para evitar falsos positivos
			result = strings.ReplaceAll(result, val, "***MASKED***")
		}
	}
	return result
}

var envVarRegex = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)

// ResolveEnv substitui variáveis no formato ${VAR_NAME} pelos valores do Cofre de Segredos ou variáveis do SO.
func (v *Vault) ResolveEnv(tmpl string) string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return envVarRegex.ReplaceAllStringFunc(tmpl, func(match string) string {
		varName := match[2 : len(match)-1]
		if val, ok := v.secrets[varName]; ok {
			return val
		}
		if sysVal := os.Getenv(varName); sysVal != "" {
			return sysVal
		}
		return match
	})
}

// encrypt cifra os bytes usando AES-256-GCM com nonce aleatório prependado.
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt autentica e decifra os bytes usando AES-256-GCM.
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext muito curto")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

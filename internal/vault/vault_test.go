package vault_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/vault"
)

func TestVault_EncryptDecryptAndPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_test_*")
	if err != nil {
		t.Fatalf("falha ao criar tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vaultFile := filepath.Join(tmpDir, "vault.enc")
	masterKey := "super-secret-master-key-12345"

	// 1. Cria vault novo e insere segredos
	v, err := vault.New(vaultFile, masterKey)
	if err != nil {
		t.Fatalf("falha ao instanciar vault: %v", err)
	}

	if err := v.Set("OPENAI_API_KEY", "sk-proj-1234567890abcdef"); err != nil {
		t.Fatalf("erro ao setar segredo: %v", err)
	}
	if err := v.Set("DB_PASSWORD", "mypassword!@#"); err != nil {
		t.Fatalf("erro ao setar segredo: %v", err)
	}

	// Persiste em disco
	if err := v.Save(); err != nil {
		t.Fatalf("falha ao salvar vault: %v", err)
	}

	// 2. Garante que o arquivo salvo em disco NÃO contém os segredos em texto puro
	rawDisk, err := os.ReadFile(vaultFile)
	if err != nil {
		t.Fatalf("falha ao ler vaultFile do disco: %v", err)
	}
	if strings.Contains(string(rawDisk), "sk-proj-1234567890abcdef") {
		t.Fatalf("vazamento de segredo! texto puro encontrado no arquivo gravado")
	}
	if strings.Contains(string(rawDisk), "mypassword!@#") {
		t.Fatalf("vazamento de senha! texto puro encontrado no arquivo gravado")
	}

	// 3. Carrega o vault em uma nova instância usando a mesma chave mestra
	v2, err := vault.New(vaultFile, masterKey)
	if err != nil {
		t.Fatalf("falha ao instanciar vault 2: %v", err)
	}
	if err := v2.Load(); err != nil {
		t.Fatalf("falha ao carregar vault existente: %v", err)
	}

	val, ok := v2.Get("OPENAI_API_KEY")
	if !ok || val != "sk-proj-1234567890abcdef" {
		t.Fatalf("valor recuperado incorreto: got %q, ok=%v", val, ok)
	}

	valPass, ok := v2.Get("DB_PASSWORD")
	if !ok || valPass != "mypassword!@#" {
		t.Fatalf("valor de senha recuperado incorreto: got %q, ok=%v", valPass, ok)
	}
}

func TestVault_WrongMasterKeyFails(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_test_wrong_*")
	if err != nil {
		t.Fatalf("falha ao criar tmpDir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vaultFile := filepath.Join(tmpDir, "vault.enc")
	v, err := vault.New(vaultFile, "correct-key")
	if err != nil {
		t.Fatalf("falha ao criar vault: %v", err)
	}
	_ = v.Set("SECRET", "valor-ultra-confidencial")
	_ = v.Save()

	// Tenta carregar com chave incorreta
	vWrong, err := vault.New(vaultFile, "wrong-key")
	if err != nil {
		t.Fatalf("falha ao criar vault com wrong-key: %v", err)
	}
	err = vWrong.Load()
	if err == nil {
		t.Fatalf("esperava erro de autenticação/decriptografia com chave errada, mas retornou nil")
	}
}

func TestVault_MaskSecretsAndResolveEnv(t *testing.T) {
	v, err := vault.New("", "test-key")
	if err != nil {
		t.Fatalf("falha ao criar vault in-memory: %v", err)
	}

	_ = v.Set("API_TOKEN", "super-token-xyz-999")
	_ = v.Set("PORT", "8080")

	// Teste de Mascaramento
	rawLog := "Conectando ao provedor com token super-token-xyz-999 na porta 8080"
	masked := v.MaskSecrets(rawLog)
	if strings.Contains(masked, "super-token-xyz-999") {
		t.Fatalf("falha no mascaramento, segredo exposto no log: %s", masked)
	}
	if !strings.Contains(masked, "***MASKED***") {
		t.Fatalf("esperava máscara no texto sanitizado: %s", masked)
	}

	// Teste de Resolução de Template
	tmpl := "curl -H 'Authorization: Bearer ${API_TOKEN}' http://localhost:${PORT}"
	resolved := v.ResolveEnv(tmpl)
	expected := "curl -H 'Authorization: Bearer super-token-xyz-999' http://localhost:8080"
	if resolved != expected {
		t.Fatalf("esperava %q, obteve %q", expected, resolved)
	}
}

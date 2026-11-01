package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MCPCatalogEntry descreve um MCP do catálogo embutido no LegoLoop (R-MCP-2).
type MCPCatalogEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	Command     string `json:"command"`
	Args        []string `json:"args,omitempty"`
	Installed   bool   `json:"installed"`
}

var builtinMCPCatalog = []MCPCatalogEntry{
	{
		Name:        "playwright",
		Description: "Automação e navegação web headless com Playwright",
		Origin:      "npm",
		Command:     "npx",
		Args:        []string{"-y", "@playwright/mcp@latest"},
	},
	{
		Name:        "codebase-memory-mcp",
		Description: "Grafo semântico e memória persistente do repositório",
		Origin:      "npm",
		Command:     "npx",
		Args:        []string{"-y", "codebase-memory-mcp"},
	},
	{
		Name:        "rea",
		Description: "Engenharia reversa universal e dissecação profunda de alvos",
		Origin:      "local/npm",
		Command:     "npx",
		Args:        []string{"-y", "@rea/engine"},
	},
	{
		Name:        "chrome-devtools-mcp",
		Description: "Inspeção e controle via Chrome DevTools Protocol",
		Origin:      "npm",
		Command:     "npx",
		Args:        []string{"-y", "chrome-devtools-mcp"},
	},
}

func (s *Server) handleMCPCatalog(w http.ResponseWriter, _ *http.Request) {
	catalog := make([]MCPCatalogEntry, len(builtinMCPCatalog))
	copy(catalog, builtinMCPCatalog)

	// Checa se já está instalado no inventário de assets
	if s.Assets != nil {
		snap := s.Assets.Snapshot()
		installedNames := map[string]bool{}
		for _, item := range snap.Items {
			if item.Kind == "mcp" {
				installedNames[strings.ToLower(item.Name)] = true
			}
		}
		for i := range catalog {
			if installedNames[strings.ToLower(catalog[i].Name)] {
				catalog[i].Installed = true
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"catalog": catalog,
	})
}

type mcpInstallRequest struct {
	Name    string `json:"name"`
	Consent bool   `json:"consent"`
}

func (s *Server) handleMCPInstall(w http.ResponseWriter, r *http.Request) {
	var req mcpInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "corpo inválido: " + err.Error()})
		return
	}

	if !req.Consent {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "instalação de MCP requer consentimento explícito do usuário (R-MCP-2)",
		})
		return
	}

	var found *MCPCatalogEntry
	for _, entry := range builtinMCPCatalog {
		if strings.EqualFold(entry.Name, req.Name) {
			found = &entry
			break
		}
	}

	if found == nil {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": fmt.Sprintf("MCP %q não encontrado no catálogo", req.Name)})
		return
	}

	// Registra o servidor MCP no arquivo .mcp.json do projeto
	mcpFilePath := filepath.Join(s.ProjectDir, ".mcp.json")
	type projectMCPFile struct {
		MCPServers map[string]map[string]any `json:"mcpServers"`
	}

	var mcpConfig projectMCPFile
	if data, err := os.ReadFile(mcpFilePath); err == nil {
		_ = json.Unmarshal(data, &mcpConfig)
	}
	if mcpConfig.MCPServers == nil {
		mcpConfig.MCPServers = make(map[string]map[string]any)
	}

	mcpConfig.MCPServers[found.Name] = map[string]any{
		"command": found.Command,
		"args":    found.Args,
	}

	data, err := json.MarshalIndent(mcpConfig, "", "  ")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "falha ao serializar .mcp.json: " + err.Error()})
		return
	}

	if err := os.WriteFile(mcpFilePath, data, 0o644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "falha ao gravar .mcp.json: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"name":    found.Name,
		"origin":  found.Origin,
		"command": found.Command,
		"file":    mcpFilePath,
	})
}

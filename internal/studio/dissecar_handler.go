package studio

import (
	"encoding/json"
	"net/http"
)

type DissecarRequest struct {
	Target string `json:"target"`
}

type DissecarActionRequest struct {
	Action string `json:"action"`
}

func handleDissecar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"targets": []string{"pipeline", "tool", "component"},
		})
		return
	}

	if r.Method == http.MethodPost {
		var req DissecarRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "dissected",
			"target":  req.Target,
			"canvas":  "3d",
			"actions": []string{"ingest", "forensics", "curate", "govern"},
		})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleDissecarAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DissecarActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validAction := false
	for _, a := range []string{"ingest", "forensics", "curate", "govern"} {
		if req.Action == a {
			validAction = true
			break
		}
	}

	if !validAction {
		http.Error(w, "Invalid cycle action", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"action":  req.Action,
		"message": "Ciclo executado com sucesso",
	})
}

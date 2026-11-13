package studio

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDissecarHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/dissecar", nil)
	rec := httptest.NewRecorder()
	handleDissecar(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperado 200, obteve %d", rec.Code)
	}

	bodyReq := DissecarRequest{Target: "pipeline"}
	jsonBody, _ := json.Marshal(bodyReq)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/dissecar", bytes.NewBuffer(jsonBody))
	recPost := httptest.NewRecorder()
	handleDissecar(recPost, reqPost)

	if recPost.Code != http.StatusOK {
		t.Errorf("esperado 200 no post, obteve %d", recPost.Code)
	}
}

func TestDissecarActionHandler(t *testing.T) {
	bodyReq := DissecarActionRequest{Action: "ingest"}
	jsonBody, _ := json.Marshal(bodyReq)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/dissecar/action", bytes.NewBuffer(jsonBody))
	recPost := httptest.NewRecorder()
	handleDissecarAction(recPost, reqPost)

	if recPost.Code != http.StatusOK {
		t.Errorf("esperado 200 na acao, obteve %d", recPost.Code)
	}
}

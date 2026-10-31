package studio

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"sync"
	"time"
)

// maxBodyBytes limita o corpo de qualquer POST.
const maxBodyBytes = 1 << 20

const csp = "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; " +
	"connect-src 'self'; font-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'; object-src 'none'"

// secure aplica, na ordem: cabeçalhos, Host allowlist, rate limit e regras de POST.
func (s *Server) secure(next http.Handler) http.Handler {
	lim := newLimiter(s.RatePerSecond, s.Burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", csp)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("Cache-Control", "no-store")
		if !s.hostAllowed(r.Host) {
			writeError(w, http.StatusMisdirectedRequest, "host não permitido")
			return
		}
		if !lim.allow() {
			writeError(w, http.StatusTooManyRequests, "muitas requisições")
			return
		}
		if r.Method == http.MethodPost && !s.postAllowed(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// hostAllowed bloqueia DNS rebinding: só o endereço de escuta (ou localhost na mesma porta).
func (s *Server) hostAllowed(host string) bool {
	if host == s.Host {
		return true
	}
	_, port, err := net.SplitHostPort(s.Host)
	return err == nil && host == net.JoinHostPort("localhost", port)
}

// postAllowed exige mesma origem e application/json (impede CSRF por formulário simples).
func (s *Server) postAllowed(w http.ResponseWriter, r *http.Request) bool {
	if o := r.Header.Get("Origin"); o != "" && o != "http://"+r.Host {
		writeError(w, http.StatusForbidden, "origem cruzada bloqueada")
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" && site != "none" {
		writeError(w, http.StatusForbidden, "requisição cross-site bloqueada")
		return false
	}
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type deve ser application/json")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	return true
}

// decodeComposition decodifica estritamente (campos desconhecidos ⇒ 400).
func decodeComposition(w http.ResponseWriter, r *http.Request) (Composition, bool) {
	var c Composition
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&c)
	var tooBig *http.MaxBytesError
	switch {
	case errors.As(err, &tooBig):
		writeError(w, http.StatusRequestEntityTooLarge, "corpo excede 1 MiB")
		return c, false
	case err != nil:
		writeError(w, http.StatusBadRequest, "composição inválida: "+err.Error())
		return c, false
	case !atEOF(dec):
		writeError(w, http.StatusBadRequest, "composição inválida: dados extras após o objeto JSON")
		return c, false
	}
	return c, true
}

// limiter é um token bucket simples (servidor local de usuário único).
type limiter struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	tokens float64
	last   time.Time
}

func newLimiter(rate float64, burst int) *limiter {
	if rate <= 0 {
		rate = 50
	}
	if burst <= 0 {
		burst = 100
	}
	return &limiter{rate: rate, burst: float64(burst), tokens: float64(burst), last: time.Now()}
}

func (l *limiter) allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.tokens = min(l.burst, l.tokens+now.Sub(l.last).Seconds()*l.rate)
	l.last = now
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// atEOF confirma que não há nada além do objeto JSON (inclusive '}' ou ']' soltos).
func atEOF(dec *json.Decoder) bool {
	_, err := dec.Token()
	return errors.Is(err, io.EOF)
}

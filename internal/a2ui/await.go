package a2ui

import (
	"bufio"
	"context"
	"io"
	"strings"
)

type stdinResult struct {
	texto string
	err   error
}

// AwaitAction aguarda em corrida a primeira resposta entre o canal da UI e o stdin.
func AwaitAction(ctx context.Context, sessionID, surfaceID string, uiCh <-chan Action, stdin io.Reader) (Action, error) {
	if err := ctx.Err(); err != nil {
		return Action{}, err
	}
	stdinCh := iniciarLeituraStdin(stdin)
	for uiCh != nil || stdinCh != nil {
		select {
		case <-ctx.Done():
			return Action{}, ctx.Err()
		case act, ok := <-uiCh:
			if !ok {
				uiCh = nil
				continue
			}
			return normalizarAcaoUI(act, sessionID, surfaceID), nil
		case res := <-stdinCh:
			if res.err != nil {
				if uiCh != nil {
					stdinCh = nil
					continue
				}
				return Action{}, res.err
			}
			return acaoStdin(sessionID, surfaceID, res.texto), nil
		}
	}
	return Action{}, io.EOF
}

func iniciarLeituraStdin(r io.Reader) <-chan stdinResult {
	if r == nil {
		return nil
	}
	ch := make(chan stdinResult, 1)
	go func() {
		linha, err := bufio.NewReader(r).ReadString('\n')
		texto := strings.TrimSpace(linha)
		if err != nil && (err != io.EOF || texto == "") {
			ch <- stdinResult{err: err}
			return
		}
		ch <- stdinResult{texto: texto}
	}()
	return ch
}

func normalizarAcaoUI(act Action, sessionID, surfaceID string) Action {
	act.Origin = "ui"
	if act.SessionID == "" {
		act.SessionID = sessionID
	}
	if act.SurfaceID == "" {
		act.SurfaceID = surfaceID
	}
	return act
}

func acaoStdin(sessionID, surfaceID, texto string) Action {
	return Action{
		SessionID: sessionID,
		SurfaceID: surfaceID,
		Origin:    "stdin",
		Name:      "stdin_reply",
		Payload:   map[string]any{"text": texto},
	}
}

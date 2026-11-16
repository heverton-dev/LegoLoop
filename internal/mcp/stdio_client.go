package mcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const callTimeout = 90 * time.Second

// stdioClient fala JSON-RPC (uma mensagem por linha) com um servidor MCP filho.
type stdioClient struct {
	cmd    *exec.Cmd
	in     io.WriteCloser
	out    *bufio.Reader
	mu     sync.Mutex
	nextID int
}

func startStdio(spec ServerSpec) (*stdioClient, error) {
	if spec.Command == "" {
		return nil, errors.New("servidor sem command (remoto/URL nao suportado no proxy stdio)")
	}
	cmd := exec.Command(spec.Command, spec.Args...)
	cmd.Dir = spec.Cwd
	cmd.Env = os.Environ()
	for k, v := range spec.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stderr = io.Discard
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	c := &stdioClient{cmd: cmd, in: in, out: bufio.NewReaderSize(out, 1<<20)}
	if err := c.handshake(); err != nil {
		c.close()
		return nil, err
	}
	return c, nil
}

func (c *stdioClient) handshake() error {
	params := map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "legoloop-mcp-proxy", "version": "1.0.0"},
	}
	if _, err := c.call("initialize", params); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}
	return c.notify("notifications/initialized")
}

func (c *stdioClient) notify(method string) error {
	msg, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "method": method})
	_, err := c.in.Write(append(msg, '\n'))
	return err
}

// call envia uma requisicao e espera a resposta de mesmo id, com timeout.
func (c *stdioClient) call(method string, params any) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextID++
	id := c.nextID
	msg, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params})
	if _, err := c.in.Write(append(msg, '\n')); err != nil {
		return nil, err
	}
	type reply struct {
		raw json.RawMessage
		err error
	}
	ch := make(chan reply, 1)
	go func() {
		raw, err := c.readResponse(id)
		ch <- reply{raw, err}
	}()
	select {
	case r := <-ch:
		return r.raw, r.err
	case <-time.After(callTimeout):
		c.close()
		return nil, fmt.Errorf("timeout de %s em %s", callTimeout, method)
	}
}

func (c *stdioClient) readResponse(id int) (json.RawMessage, error) {
	for {
		line, err := c.out.ReadBytes('\n')
		if len(line) > 0 {
			var resp struct {
				ID     any             `json:"id"`
				Result json.RawMessage `json:"result"`
				Error  *struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if json.Unmarshal(line, &resp) == nil && sameID(resp.ID, id) {
				if resp.Error != nil {
					return nil, errors.New(resp.Error.Message)
				}
				return resp.Result, nil
			}
		}
		if err != nil {
			return nil, err
		}
	}
}

func sameID(got any, want int) bool {
	f, ok := got.(float64)
	return ok && int(f) == want
}

// close encerra o servidor e toda a arvore de processos (npx -> node).
func (c *stdioClient) close() {
	_ = c.in.Close()
	if c.cmd.Process != nil {
		killTree(c.cmd)
	}
	_ = c.cmd.Wait()
}

func killTree(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
		return
	}
	_ = cmd.Process.Kill()
}

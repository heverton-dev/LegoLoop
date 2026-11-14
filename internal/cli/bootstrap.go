package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type bootstrapResult struct {
	Vault           string   `json:"vault"`
	Skills          []string `json:"skills"`
	Harnesses       []string `json:"harnesses"`
	HarnessesSynced bool     `json:"harnesses_synced"`
	EnvReady        bool     `json:"env_ready"`
	DryRun          bool     `json:"dry_run"`
	Status          string   `json:"status"`
}

func cmdBootstrap(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("bootstrap", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	dry := fs.Bool("dry-run", false, "apenas planeja sem gravar")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	v, err := resolveVault(*vault)
	if err != nil {
		return fail(d, err)
	}
	res, err := runBootstrap(ctx, d, v, *dry)
	if err != nil {
		return fail(d, err)
	}
	return reportBootstrap(d, res, *asJSON)
}

func runBootstrap(ctx context.Context, d Deps, vault string, dry bool) (*bootstrapResult, error) {
	skRes, err := runInstalarSkills(d, vault, dry)
	if err != nil {
		return nil, fmt.Errorf("skills: %w", err)
	}
	syncFlag := "--apply"
	if dry {
		syncFlag = "--dry-run"
	}
	syncDeps := Deps{Stdout: io.Discard, Stderr: d.Stderr, Home: d.Home}
	if code := cmdSyncHarnesses(ctx, []string{syncFlag}, syncDeps); code != 0 {
		return nil, fmt.Errorf("sync-harnesses retornou código %d", code)
	}
	envReady := checkBootstrapEnv(d, vault, dry)
	return &bootstrapResult{
		Vault:           vault,
		Skills:          skRes.Skills,
		Harnesses:       skRes.Harnesses,
		HarnessesSynced: true,
		EnvReady:        envReady,
		DryRun:          dry,
		Status:          "ok",
	}, nil
}

func checkBootstrapEnv(d Deps, vault string, dry bool) bool {
	if dry {
		return true
	}
	if _, err := os.Stat(vault); err != nil {
		return false
	}
	home := resolveHome(d, "")
	return home != ""
}

func reportBootstrap(d Deps, res *bootstrapResult, asJSON bool) int {
	if asJSON {
		enc := json.NewEncoder(d.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(res)
		return 0
	}
	fmt.Fprintf(d.Stdout, "Bootstrap %s: %d skills, %d harnesses, mcp-router=%v, env=%v\n",
		res.Status, len(res.Skills), len(res.Harnesses), res.HarnessesSynced, res.EnvReady)
	return 0
}

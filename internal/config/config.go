// Package config resolves jjay's own configuration — distinct from
// openspec/config.yaml. Today it carries per-agent command profiles
// {Launch, Resume}, resolved per field from an ordered list of sources
// (project .jjay/config.yaml → global ~/.config/jjay/config.yaml → Go built-in).
// See ADR-014. The built-in is the single source of truth: `jjay init`
// materializes it into the project file, and runtime falls back to it, so the
// seeded file and the fallback cannot drift.
package config

import (
	"os"
	"path/filepath"

	yaml "go.yaml.in/yaml/v3"
)

// AgentProfile is the pair of command templates for one agent. Both are
// templates over {change}/{prompt}/{wsdir}. Launch starts the work; Resume
// reopens an existing workspace without re-running the work.
type AgentProfile struct {
	Launch string `yaml:"launch"`
	Resume string `yaml:"resume"`
}

// Config is the on-disk jjay config schema (one file layer).
type Config struct {
	Agents map[string]AgentProfile `yaml:"agents"`
}

// builtin is the Go-native default config — the single source of truth for both
// the runtime fallback and the file `jjay init` seeds. The launch command is the
// historical spawn.DefaultAgentCommand (kept identical so spawn behavior is
// unchanged when no config exists).
var builtin = Config{
	Agents: map[string]AgentProfile{
		"claude": {
			Launch: `claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}`,
			Resume: `claude --resume --add-dir {wsdir}`,
		},
	},
}

// Builtin returns a copy of the built-in default config. Used by `jjay init` to
// seed <repo>/.jjay/config.yaml from the same const the runtime falls back to.
func Builtin() Config {
	out := Config{Agents: map[string]AgentProfile{}}
	for name, p := range builtin.Agents {
		out.Agents[name] = p
	}
	return out
}

// Load reads a single config file. A missing file yields an empty Config and no
// error (an absent layer simply contributes nothing to resolution).
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// ProjectConfigPath returns the project config path for a repo root.
func ProjectConfigPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".jjay", "config.yaml")
}

// GlobalConfigPath returns ~/.config/jjay/config.yaml (empty if HOME is unset).
func GlobalConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "jjay", "config.yaml")
}

// resolveField returns the first non-empty value of agents[agent].<field> across
// the ordered sources, falling through to the built-in. An empty string counts
// as unset, so a present-but-blank field still falls back. field selects Launch
// or Resume via the accessor.
func resolveField(agent string, accessor func(AgentProfile) string, sources []Config) string {
	for _, src := range sources {
		if p, ok := src.Agents[agent]; ok {
			if v := accessor(p); v != "" {
				return v
			}
		}
	}
	if p, ok := builtin.Agents[agent]; ok {
		return accessor(p)
	}
	return ""
}

// ResolveProfile resolves an agent's full profile per field over the ordered
// sources, with the built-in as the final fallback. Each field is resolved
// independently: a source that sets only Launch keeps a later source's Resume.
func ResolveProfile(agent string, sources []Config) AgentProfile {
	return AgentProfile{
		Launch: resolveField(agent, func(p AgentProfile) string { return p.Launch }, sources),
		Resume: resolveField(agent, func(p AgentProfile) string { return p.Resume }, sources),
	}
}

// Resolve loads the standard source order (project → global) and resolves the
// agent's profile, with the built-in as the final fallback. repoRoot locates the
// project file; pass the main repo root, not a spawned workspace. Either file may
// be absent. The agent name selects which profile; today only "claude" is
// populated in the built-in.
func Resolve(agent, repoRoot string) (AgentProfile, error) {
	var sources []Config

	if repoRoot != "" {
		proj, err := Load(ProjectConfigPath(repoRoot))
		if err != nil {
			return AgentProfile{}, err
		}
		sources = append(sources, proj)
	}
	if gp := GlobalConfigPath(); gp != "" {
		glob, err := Load(gp)
		if err != nil {
			return AgentProfile{}, err
		}
		sources = append(sources, glob)
	}
	return ResolveProfile(agent, sources), nil
}

// DefaultAgent is the agent jjay resolves when none is otherwise specified.
const DefaultAgent = "claude"

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProfile_NoSourcesUsesBuiltin(t *testing.T) {
	p := ResolveProfile("claude", nil)
	if p.Launch != builtin.Agents["claude"].Launch {
		t.Errorf("launch should be built-in, got %q", p.Launch)
	}
	if p.Resume != builtin.Agents["claude"].Resume {
		t.Errorf("resume should be built-in, got %q", p.Resume)
	}
}

func TestResolveProfile_ProjectOverridesOneFieldInheritsRest(t *testing.T) {
	project := Config{Agents: map[string]AgentProfile{
		"claude": {Launch: "custom-launch {change}"}, // resume left blank
	}}
	p := ResolveProfile("claude", []Config{project})
	if p.Launch != "custom-launch {change}" {
		t.Errorf("launch should be the project override, got %q", p.Launch)
	}
	if p.Resume != builtin.Agents["claude"].Resume {
		t.Errorf("resume should fall back to built-in, got %q", p.Resume)
	}
}

func TestResolveProfile_OrderedSourcePrecedence(t *testing.T) {
	project := Config{Agents: map[string]AgentProfile{"claude": {Resume: "project-resume"}}}
	global := Config{Agents: map[string]AgentProfile{"claude": {Launch: "global-launch", Resume: "global-resume"}}}
	// project first, then global: project wins where set, global fills the gap.
	p := ResolveProfile("claude", []Config{project, global})
	if p.Resume != "project-resume" {
		t.Errorf("project resume should win, got %q", p.Resume)
	}
	if p.Launch != "global-launch" {
		t.Errorf("launch should come from global (project left it blank), got %q", p.Launch)
	}
}

func TestResolveProfile_PresentButBlankFallsBack(t *testing.T) {
	project := Config{Agents: map[string]AgentProfile{"claude": {Launch: "x", Resume: ""}}}
	p := ResolveProfile("claude", []Config{project})
	if p.Resume != builtin.Agents["claude"].Resume {
		t.Errorf("blank resume should fall back to built-in, got %q", p.Resume)
	}
}

func TestLoad_MissingFileIsEmptyNoError(t *testing.T) {
	c, err := Load(filepath.Join(t.TempDir(), "nope.yaml"))
	if err != nil {
		t.Fatalf("missing file should not error, got %v", err)
	}
	if len(c.Agents) != 0 {
		t.Errorf("missing file should yield empty config, got %v", c.Agents)
	}
}

func TestLoad_ParsesYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	yaml := "agents:\n  claude:\n    launch: 'L {change}'\n    resume: 'R {wsdir}'\n"
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.Agents["claude"].Launch != "L {change}" || c.Agents["claude"].Resume != "R {wsdir}" {
		t.Errorf("parsed wrong: %+v", c.Agents["claude"])
	}
}

func TestBuiltin_IsACopy(t *testing.T) {
	b := Builtin()
	b.Agents["claude"] = AgentProfile{Launch: "mutated"}
	if builtin.Agents["claude"].Launch == "mutated" {
		t.Error("Builtin() must return a copy; mutation leaked into the package var")
	}
}

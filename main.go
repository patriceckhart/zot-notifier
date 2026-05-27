package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/patriceckhart/zot/packages/agent/ext"
)

type config struct {
	Enabled bool `json:"enabled"`
}

var (
	cfg   config
	cfgMu sync.Mutex
)

func configPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "config.json")
}

func soundPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "ding.mp3")
}

func loadConfig() {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	data, err := os.ReadFile(configPath())
	if err != nil {
		cfg = config{Enabled: true}
		return
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		cfg = config{Enabled: true}
	}
}

func saveConfig() {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(configPath(), data, 0o644)
}

func isEnabled() bool {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	return cfg.Enabled
}

func setEnabled(v bool) {
	cfgMu.Lock()
	cfg.Enabled = v
	cfgMu.Unlock()
	saveConfig()
}

func playSound() {
	if !isEnabled() {
		return
	}
	_ = exec.Command("afplay", soundPath()).Start()
}

const panelID = "notifier-settings"

func renderPanel(e *ext.Extension) {
	var toggle string
	if isEnabled() {
		toggle = "  [enter] Sound: on"
	} else {
		toggle = "  [enter] Sound: off"
	}
	e.RenderPanel(panelID, "Notifier", []string{toggle}, "enter: toggle · esc: close")
}

func main() {
	loadConfig()

	e := ext.New("notifier", "1.0.0")

	// Only play when the model actually finished a reply (final summary).
	// Skip tool_use turns (intermediate), errors, and user aborts.
	e.On("turn_end", func(ev ext.Event) {
		if ev.Stop != "end" {
			return
		}
		playSound()
	})

	e.OnPanelKey(panelID, func(key, text string) {
		if key == "enter" {
			setEnabled(!isEnabled())
			renderPanel(e)
		}
	}, func() {
		// panel closed — nothing to clean up
	})

	e.Command("notifier", "configure notification sound", func(args string) ext.Response {
		return ext.OpenPanel(panelID, "Notifier", panelLines(), "enter: toggle · esc: close")
	})

	if err := e.Run(); err != nil {
		e.Logf("fatal: %v", err)
	}
}

func panelLines() []string {
	if isEnabled() {
		return []string{"  [enter] Sound: on"}
	}
	return []string{"  [enter] Sound: off"}
}

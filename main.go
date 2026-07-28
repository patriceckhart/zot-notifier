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

// extDir returns the directory the extension was launched from.
// zot sets cmd.Dir = ext.Dir before exec, so os.Getwd() gives us the
// extension folder regardless of whether we run from a committed
// binary, `go run .`, or any other launcher.
//
// We deliberately do NOT use os.Executable(): with `go run .` that
// points at a temp build dir under $GOCACHE, which breaks both
// config persistence and the path to ding.mp3.
func extDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func configPath() string { return filepath.Join(extDir(), "config.json") }
func soundPath() string  { return filepath.Join(extDir(), "ding.mp3") }

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

func shouldNotify(ev ext.Event) bool {
	return ev.Name == "tool_confirmation_requested" || (ev.Name == "turn_end" && ev.Stop == "end")
}

func playSound(e *ext.Extension) {
	if !isEnabled() {
		return
	}
	sp := soundPath()
	// Log if the sound file is missing or afplay fails so we get a
	// hint in `zot ext logs notifier` instead of silent nothing.
	if _, err := os.Stat(sp); err != nil {
		e.Logf("sound file missing: %s (%v)", sp, err)
		return
	}
	cmd := exec.Command("afplay", sp)
	if err := cmd.Start(); err != nil {
		e.Logf("afplay failed to start: %v", err)
		return
	}
	// Reap the child in the background so we don't leak zombies.
	go func() { _ = cmd.Wait() }()
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

	e := ext.New("notifier", "1.1.0")

	notify := func(ev ext.Event) {
		if shouldNotify(ev) {
			playSound(e)
		}
	}
	// Notify when zot needs input and when the model finishes its final reply.
	// Intermediate tool-use turns, errors, and user aborts remain silent.
	e.On("tool_confirmation_requested", notify)
	e.On("turn_end", notify)

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

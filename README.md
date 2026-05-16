# zot-notifier

A small [zot](https://www.zot.sh) extension that plays a notification sound when the assistant finishes a turn. Handy for long-running tasks so you don't have to babysit the terminal.

## Features

- Plays a short "ding" via macOS `afplay` whenever a turn ends with a final reply.
- Skips intermediate tool-use turns, errors, and aborts — only the real "done" moment triggers a sound.
- Toggle on/off from inside zot via the `/notifier` slash command, which opens a small settings panel.
- State persists to `config.json` next to the binary.

## Requirements

- macOS (uses the built-in `afplay` command).
- Go 1.22+ if you want to rebuild from source.
- A working [zot](https://github.com/patriceckhart/zot) installation.

## Install

Clone (or copy) this directory into your zot extensions folder. zot will pick it up automatically based on `extension.json` (it's enabled by default).

A prebuilt `notifier` binary is included, so **no build step is required** to use the extension. If you'd rather run it straight from source — or you've changed `main.go` — you can either rebuild:

```sh
go build -o notifier .
```

…or skip the binary entirely and have zot run the Go source directly by changing the `exec` field in `extension.json` to:

```json
"exec": "go run ."
```

## Usage

- Just run zot as usual — when a turn ends, you'll hear the ding.
- Type `/notifier` to open the settings panel.
  - `enter` toggles the sound on/off.
  - `esc` closes the panel.

## Files

| File | Purpose |
|------|---------|
| `extension.json` | zot extension manifest. |
| `main.go` | Extension source. |
| `notifier` | Prebuilt binary (the `exec` target). |
| `ding.mp3` | The notification sound. |
| `config.json` | Persisted on/off state. |

## Credits

Sound Effect by <a href="https://pixabay.com/users/lesiakower-25701529/?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=246413">Lesiakower</a> from <a href="https://pixabay.com//?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=246413">Pixabay</a>

## License

MIT

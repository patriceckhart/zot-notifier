# zot-notifier

A small [zot](https://www.zot.sh) extension that plays a notification sound when the assistant finishes a turn or needs tool approval. Handy for long-running tasks so you don't have to babysit the terminal.

## Features

- Plays a short "ding" via macOS `afplay` whenever a turn ends with a final reply.
- Also plays when zot is waiting for you to approve a tool call.
- Skips intermediate tool-use turns, automatically allowed calls, errors, and aborts.
- Toggle on/off from inside zot via the `/notifier` slash command, which opens a small settings panel.
- State persists to `config.json` next to the extension.

## Requirements

- macOS (uses the built-in `afplay` command).
- Go 1.22+ on `$PATH` (the extension runs from source via `go run .`, no
  prebuilt binary is shipped). Install Go from <https://go.dev/dl>.
- A working [zot](https://github.com/patriceckhart/zot) installation. Tool
  approval notifications require a zot version that supports the
  `tool_confirmation_requested` lifecycle event. Older versions still notify
  when turns finish.

## Install

```bash
zot ext install https://github.com/patriceckhart/zot-notifier
```

Or clone / copy this directory into your zot extensions folder. zot will pick it up automatically based on `extension.json` (it's enabled by default).

No prebuilt binary is shipped — `extension.json` runs `go run .`, so the
first launch compiles `main.go` on the fly. Make sure `go` is on your
`$PATH`.

## Usage

- Just run zot as usual. You will hear the ding when a turn ends or zot needs tool approval.
- Type `/notifier` to open the settings panel.
  - `enter` toggles the sound on/off.
  - `esc` closes the panel.

## Files

| File | Purpose |
|------|---------|
| `extension.json` | zot extension manifest. |
| `main.go` | Extension source. |
| `ding.mp3` | The notification sound. |
| `config.json` | Persisted on/off state. |

## Credits

Sound Effect by <a href="https://pixabay.com/users/lesiakower-25701529/?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=246413">Lesiakower</a> from <a href="https://pixabay.com//?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=246413">Pixabay</a>

## License

MIT

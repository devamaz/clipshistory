# Clipshistory [WIP]

A cli clipboard history manager built in Go with a terminal UI (TUI). Automatically monitors your clipboard, stores entries in a SQLite database, and provides an intuitive interface to search, view, and manage your clipboard history.

## Screenshots

```
┌─ SEARCH ───────────────────────────────────────────────────────┐
│                                                                                                              │
├─────────────────────────────────────────────────────────── ─┤
│  ┌─ PREVIEW ──────── 1 ─┐ ┌─ CONTENT ─────────── 3 ─┐ ┌─ PIN ──────── 2 ─┐   │
│  │ > hello world                │ │ Full Content:                      │ │ Pinned                   │   │
│  │   (11 chars)                 │ │                                    │ │ clips                    │   │
│  │   another text               │ │ hello world                        │ │                          │   │
│  │   (12 chars)                 │ │                                    │ │                          │   │
│  │   ...                        │ │                                    │ │                          │   │
│  │                              │ │                                    │ │                          │   │
│  │                              │ │                                    │ │                          │   │
│  │                              │ │                                    │ │                          │   │
│  │                              │ │                                    │ │                          │   │
│  │                              │ └────────────────────┘ │                          │   │
│  │                              │ ┌─ INFO ───────────── 4 ─┐ │                          │   │
│  │                              │ │ Clip Details:                      │ │                          │   │
│  │                              │ │                                    │ │                          │   │
│  │                              │ │ Char Count:               11       │ │                          │   │
│  │                              │ │ Times:                    1        │ │                          │   │
│  │                              │ │ Words:                    3        │ │                          │   │
│  │                              │ │ CopiedAt:                 Today    │ │                          │   │
│  │                              │ │ Last Copied:              Today    │ │                          │   │
│  └────────────── ──┘ └────────────────────┘└───────────── ─┘   │
├────────────────────────────────────────────────────────── ──┤
│ Press ↑/↓ to navigate • 'q' to quit                                                                          │
└───────────────────────────────────────────────────────── ───┘
```

## Installation

### Prerequisites

- Go 1.24 or higher
- SQLite (included as a dependency)
- fswatch (optional, for development hot-reload)

### From Source

```bash
git clone https://github.com/yourusername/clipshistory.git
cd clipshistory
go build -o aclips
```

### Using the watch script (Development)

```bash
./watch.sh
```

This will automatically rebuild and restart the application whenever you make changes.

## Usage

### Starting the Application

```bash
./aclips
```

### TUI Navigation

| Key            | Action                          |
| -------------- | ------------------------------- |
| `↑` / `k`      | Navigate up in clip list        |
| `↓` / `j`      | Navigate down in clip list      |
| `q` / `Ctrl+C` | Quit the application            |
| `y`            | Copy selected clip to clipboard |

## Roadmap

- [ ] Full clipboard monitoring integration
- [ ] Search functionality (FTS5)
- [ ] Pin/unpin clips
- [ ] Delete clips
- [ ] Export/Import clipboard history
- [ ] Keyboard shortcuts for quick actions
- [ ] Configurable content filtering

## License

MIT License.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Author

- [Ahmad Abdul-Aziz](https://x.com/devamaz)

Created with ❤️ using Go, Bubble Tea, and SQLite.

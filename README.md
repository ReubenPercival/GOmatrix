# GOmatrix

A 1:1 clone of [cmatrix](https://github.com/abishekvashok/cmatrix) in Go — the classic "Matrix" falling-code screensaver for your terminal.

## Build

```sh
go build -o cmatrix .
# or
make
```

## Usage

```sh
./cmatrix
```

Press `q` to quit.

### Command-line flags

| Flag | Description |
|------|-------------|
| `-a` | Asynchronous scroll |
| `-b` | Bold characters on |
| `-B` | All bold characters (overrides `-b`) |
| `-c` | Use Japanese characters |
| `-n` | No bold (overrides `-b` and `-B`) |
| `-o` | Old-style scrolling |
| `-s` | Screensaver mode (exit on any key) |
| `-r` | Rainbow mode |
| `-m` | Lambda mode |
| `-k` | Characters change while scrolling |
| `-L` | Lock mode |
| `-M <msg>` | Display message in center |
| `-u <0-10>` | Screen update delay (default: 4) |
| `-C <color>` | Matrix color (default: green) |
| `-h` | Print help |
| `-V` | Print version |

### Runtime controls

| Key | Action |
|-----|--------|
| `q` | Quit |
| `p` | Toggle pause |
| `a` | Toggle async scroll |
| `b` | Bold on |
| `B` | All bold |
| `n` | No bold |
| `r` | Toggle rainbow |
| `m` | Toggle lambda |
| `0-9` | Set speed |
| <code>!@#$%^&</code> | Set color (red, green, yellow, blue, magenta, cyan, white) |

## Dependencies

- [termbox-go](https://github.com/nsf/termbox-go) — terminal UI library

## License

GPL-3.0-or-later. See [LICENSE](LICENSE).

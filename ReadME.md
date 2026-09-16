# ascii-art-web

A small Go web server that turns plain text into ASCII-art banners,
using the same character-font logic as the CLI `ascii-art` project,
served through a simple browser UI.

Type a line (or several lines) of text, pick a font, and the server
renders it into a fixed-width ASCII banner.

## Features

- Renders text into ASCII art using one of three fonts:
  `standard`, `shadow`, `thinkertoy`
- Supports multi-line input (each line entered in the textarea is
  rendered on its own)
- Returns proper HTTP status codes on error (`400`, `404`, `405`, `500`)
  instead of failing silently
- Zero third-party dependencies — standard library only

## Project structure

```
.
├── banners/
│   ├── standard.txt      # glyph data for the "standard" font
│   ├── shadow.txt        # glyph data for the "shadow" font
│   └── thinkertoy.txt    # glyph data for the "thinkertoy" font
├── frontend/
│   └── index.html        # the web UI (form + fetch calls)
├── handlers/
│   └── handlers.go        # HTTP handlers
├── pkg/
│   ├── banner/            # loads a .txt font file into a rune → glyph map
│   ├── parser/             # turns raw input text into a rune matrix
│   └── renderer/           # turns the rune matrix + glyph map into ASCII art
├── main.go                 # wires up routes and starts the server
├── main_test.go             # test suite for the whole project
└── go.mod
```

## Requirements

- Go 1.21 or newer
- No external packages — everything used is part of the Go standard
  library (`net/http`, `os`, `strings`, `sync`, `errors`, `io`, `fmt`)

## How to run

From the project root:

```bash
go run .
```

You should see:

```
Server started at localhost:9091
```

Then open **http://localhost:9091** in a browser.

> The port (`9091`) is hardcoded in `main.go`. If it's already in use,
> change the string passed to `http.ListenAndServe` and restart.

## How to use it

1. Type your text into the textarea. Each line becomes its own row
   of ASCII art.
2. Pick a font with the **Standard**, **Shadow**, or **ThinkerToy**
   buttons (this just tells the server which `banners/*.txt` file to
   use for the next render — it doesn't render anything by itself).
3. Click **Render**. The result appears below the form.

### API endpoints

| Method | Path          | Description                                   |
|--------|---------------|------------------------------------------------|
| GET    | `/`           | Serves the web UI (`frontend/index.html`)      |
| POST   | `/input`      | Body = raw text to render; returns ASCII art   |
| GET/POST | `/standard` | Switches the active font to `standard`       |
| GET/POST | `/shadow`   | Switches the active font to `shadow`         |
| GET/POST | `/thinkertoy` | Switches the active font to `thinkertoy`   |

### Status codes you should expect

| Code | When |
|------|------|
| 200  | Successful render / successful font switch / successful page load |
| 400  | The submitted text contains a character outside the printable ASCII range (32–126) |
| 404  | Any path other than `/` |
| 405  | Wrong HTTP method for the endpoint |
| 500  | The active font file (`banners/*.txt`) couldn't be loaded from disk |

## How to test it

All tests live in a single file, `main_test.go`, at the project root,
and cover every package: `parser`, `banner`, `renderer`, and
`handlers`.

Run the full suite from the project root:

```bash
go test -v ./...
```

Run it with the race detector as well (recommended — the font-switch
handlers share state behind a mutex):

```bash
go test -v -race ./...
```

**Important:** run tests from the project root, not from a
subdirectory. The handler tests use the real `banners/standard.txt`
and `frontend/index.html` files with relative paths, exactly like the
running server does. If those files aren't found relative to the
current directory, the affected tests are skipped rather than failed.

### What's covered

- **`parser.Parse`** — plain text, multi-line input, empty string,
  trailing newline, boundary characters, control characters, DEL,
  non-ASCII input
- **`banner.Load`** — missing file, a valid fixture file, a truncated
  fixture file
- **`renderer.Render`** — a known glyph map, blank-line handling, and
  a regression test around rendering a character missing from the
  font map
- **`handlers`** — every status code above, for every handler,
  plus a concurrency test on the shared font-selection state

## Known limitations

- `banner.Load` expects every font file to cover the full printable
  ASCII range (32–126) in a fixed 9-line-per-character block. A font
  file that's missing or misaligned for a character the user typed
  will not be caught gracefully — check `pkg/renderer/renderer.go`
  before assuming any input is safe.
- The font is stored as shared, mutex-guarded package state rather
  than being passed per-request, so switching fonts affects every
  concurrent user of the server.
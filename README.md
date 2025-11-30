# checklist

Interactive terminal checklist that reads items from a text file and lets you move through them with arrow keys.

## Installation

Install the CLI into your `$GOBIN` using Go modules:

```bash
go install ./cmd/checklist
```

Or use the included Makefile targets:

```bash
make build    # builds bin/checklist
make install  # installs the checklist binary via go install
```

## Usage

```bash
go run ./cmd/checklist --input checklist.txt --output done.txt
```

- Provide an input file with one item per line.
- Move the active row with the ↑ and ↓ arrows (or `j`/`k`).
- The active line is marked with a leading `>` before the checkbox.
- Toggle selection with the space bar.
- Press Enter to write selected items to the output file. If nothing is selected, the app will keep running and remind you to pick something.
- Press `q`, `Esc`, or `Ctrl+C` to quit without saving.

## Releases

Tags that start with `v` trigger the GitHub Actions workflow in `.github/workflows/release.yml`, which builds Linux and macOS binaries and attaches them to a GitHub release.

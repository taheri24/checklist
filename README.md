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

### Placeholders

Lines in your checklist can contain placeholders wrapped in curly braces, such as `{ticket}` or `{env}`. After you press Enter to save, the CLI will prompt you to fill in each unique placeholder and substitute the values in the output file.

If a placeholder contains a pipe (`|`), the segments are treated as choices. For example, `{staging|production}` will present a numbered list and you can pick the option by entering the corresponding number.

#### Example

Input checklist (`checklist.txt`):

```
Deploy to {env}
Verify {service} logs
Share updates in {channel}
```

During save, you will be asked for `env`, `service`, and `channel`. If you enter `staging`, `payments`, and `#deployments`, the written output becomes:

```
- Deploy to staging
- Verify payments logs
- Share updates in #deployments
```

## Releases

Tags that start with `v` trigger the GitHub Actions workflow in `.github/workflows/release.yml`, which builds Linux and macOS binaries and attaches them to a GitHub release.

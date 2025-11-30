package checklist

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"checklist/internal/terminal"
)

type item struct {
	Text     string
	Selected bool
}

// Run executes the interactive checklist workflow.
func Run(checklistPath, outputPath string) error {
	items, err := loadItems(checklistPath)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return errors.New("checklist is empty")
	}

	state, err := terminal.EnableRawMode(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to enable raw mode: %w", err)
	}
	defer func() {
		terminal.Restore(int(os.Stdin.Fd()), state)
	}()

	reader := bufio.NewReader(os.Stdin)
	active := 0

	for {
		render(items, active, checklistPath, outputPath)
		action, idx, err := terminal.ReadKey(reader)
		if err != nil {
			return fmt.Errorf("read input: %w", err)
		}

		switch action {
		case terminal.ActionUp:
			active--
			if active < 0 {
				active = len(items) - 1
			}
		case terminal.ActionDown:
			active++
			if active >= len(items) {
				active = 0
			}
		case terminal.ActionToggle:
			items[active].Selected = !items[active].Selected
		case terminal.ActionToggleAt:
			if idx >= 0 && idx < len(items) {
				items[idx].Selected = !items[idx].Selected
			}
		case terminal.ActionEnter:
			state, err = promptAndWriteSelected(items, outputPath, state)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			fmt.Printf("\r %q\n\r saved", outputPath)
			return nil
		case terminal.ActionQuit:
			fmt.Println("\rExiting without saving")
			return nil
		case terminal.ActionNone:
			// ignore
		}
	}
}

func loadItems(path string) ([]item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open checklist: %w", err)
	}
	defer f.Close()

	var items []item
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		items = append(items, item{Text: text})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read checklist: %w", err)
	}
	return items, nil
}

func writeSelected(lines []string, path string) error {
	if len(lines) == 0 {
		return errors.New("no items selected")
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}
	return nil
}

var placeholderRegexp = regexp.MustCompile(`\{([^{}]+)\}`)

func collectPlaceholderKeys(lines []string) []string {
	seen := make(map[string]bool)
	var keys []string

	for _, line := range lines {
		matches := placeholderRegexp.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			key := match[1]
			if !seen[key] {
				seen[key] = true
				keys = append(keys, key)
			}
		}
	}

	return keys
}

func promptPlaceholderValues(keys []string, buffered map[string]string) (map[string]string, error) {
	reader := bufio.NewReader(os.Stdin)
	values := buffered

	for _, key := range keys {
		if _, ok := values[key]; ok {
			continue
		}

		fmt.Printf("\rEnter value for %s: ", key)
		text, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read placeholder %q: %w", key, err)
		}
		values[key] = strings.TrimSpace(text)
	}

	return values, nil
}

func replacePlaceholders(lines []string, values map[string]string) []string {
	if len(values) == 0 {
		return lines
	}

	resolved := make([]string, 0, len(lines))

	replacer := func(match string) string {
		parts := placeholderRegexp.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}
		if val, ok := values[parts[1]]; ok {
			return val
		}

		return match
	}

	for _, line := range lines {
		resolved = append(resolved, fmt.Sprintf("- %s", placeholderRegexp.ReplaceAllStringFunc(line, replacer)))
	}

	return resolved
}

func promptAndWriteSelected(items []item, outputPath string, state *terminal.State) (*terminal.State, error) {
	var selectedLines []string
	for _, it := range items {
		if it.Selected {
			selectedLines = append(selectedLines, it.Text)
		}
	}

	keys := collectPlaceholderKeys(selectedLines)
	placeholderValues := make(map[string]string)
	if len(keys) > 0 {
		if err := terminal.Restore(int(os.Stdin.Fd()), state); err != nil {
			return state, fmt.Errorf("restore terminal: %w", err)
		}

		values, err := promptPlaceholderValues(keys, placeholderValues)
		newState, rawErr := terminal.EnableRawMode(int(os.Stdin.Fd()))
		if rawErr != nil {
			return state, fmt.Errorf("failed to re-enable raw mode: %w", rawErr)
		}
		state = newState
		if err != nil {
			return state, err
		}
		placeholderValues = values
	}

	resolved := replacePlaceholders(selectedLines, placeholderValues)
	if err := writeSelected(resolved, outputPath); err != nil {
		return state, err
	}

	return state, nil
}

func numToChar(n int) string {
	if n < 0 || n > 26 {
		return ""
	}
	return string(rune('A') + rune(n))
}

func render(items []item, active int, checklistPath, outputPath string) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Interactive checklist")
	fmt.Printf("\rSource: %s\n", checklistPath)
	fmt.Printf("Output: %s\n\r", outputPath)
	fmt.Println(numToChar(2))
	fmt.Println("Use ↑/↓ to move, space to toggle, digits/letters to toggle an item directly,\n\rEnter to save, q or Esc to quit.")
	fmt.Println()
	for idx, it := range items {
		orderCh := strconv.FormatInt(int64(idx)+1, 10)
		if idx >= 9 {
			orderCh = numToChar(idx - 9)
		}
		pointer := " "
		if idx == active {
			pointer = ">"
		}
		check := " "
		if it.Selected {
			check = "x"
		}
		fmt.Printf("\r%s-[%s] %s.%s\n", pointer, check, orderCh, it.Text)
	}
}

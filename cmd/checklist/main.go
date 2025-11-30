package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"checklist/internal/terminal"
)

type item struct {
	Text     string
	Selected bool
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

func writeSelected(items []item, path string) error {
	var lines []string
	for _, it := range items {
		if it.Selected {
			lines = append(lines, it.Text)
		}
	}
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

func main() {
	checklistPath := flag.String("input", "checklist.txt", "path to checklist file with one item per line")
	outputPath := flag.String("output", "selected.txt", "path where selected items will be written")
	flag.Parse()

	items, err := loadItems(*checklistPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(items) == 0 {
		fmt.Fprintln(os.Stderr, "checklist is empty")
		os.Exit(1)
	}

	state, err := terminal.EnableRawMode(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to enable raw mode:", err)
		os.Exit(1)
	}
	defer terminal.Restore(int(os.Stdin.Fd()), state)

	reader := bufio.NewReader(os.Stdin)
	active := 0

	for {
		render(items, active, *checklistPath, *outputPath)
		action, idx, err := terminal.ReadKey(reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read input:", err)
			return
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
			if err := writeSelected(items, *outputPath); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			fmt.Printf("\rSaved in %q\n\r", *outputPath)
			return
		case terminal.ActionQuit:
			fmt.Println("\rExiting without saving")
			return
		case terminal.ActionNone:
			// ignore
		}
	}
	fmt.Println()
}

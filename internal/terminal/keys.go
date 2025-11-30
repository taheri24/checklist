package terminal

import (
	"bufio"
)

type Action int

const (
	ActionNone Action = iota
	ActionUp
	ActionDown
	ActionToggle
	ActionToggleAt
	ActionEnter
	ActionQuit
)

// ReadKey interprets key presses from a raw-mode reader.
func ReadKey(reader *bufio.Reader) (Action, int, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return ActionNone, -1, err
	}

	switch b {
	case 'q', 'Q':
		return ActionQuit, -1, nil
	case ' ', 'x', 'X':
		return ActionToggle, -1, nil
	case '\r', '\n':
		return ActionEnter, -1, nil
	case 3: // Ctrl+C
		return ActionQuit, -1, nil
	case 0x1b:
		// escape sequence, attempt to parse arrows
		if reader.Buffered() >= 2 {
			seq, err := reader.Peek(2)
			if err == nil && len(seq) == 2 && seq[0] == '[' {
				reader.Discard(2)
				switch seq[1] {
				case 'A':
					return ActionUp, -1, nil
				case 'B':
					return ActionDown, -1, nil
				case 'C':
					return ActionDown, -1, nil
				case 'D':
					return ActionUp, -1, nil
				}
			}
		}
		return ActionQuit, -1, nil
	case 'k':
		return ActionUp, -1, nil
	case 'j':
		return ActionDown, -1, nil
	}

	switch {
	case b >= '1' && b <= '9':
		return ActionToggleAt, int(b - '1'), nil
	case b >= 'a' && b <= 'z':
		return ActionToggleAt, 9 + int(b-'a'), nil
	case b >= 'A' && b <= 'Z':
		return ActionToggleAt, 9 + int(b-'A'), nil
	}

	return ActionNone, -1, nil
}

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
	ActionEnter
	ActionQuit
)

// ReadKey interprets key presses from a raw-mode reader.
func ReadKey(reader *bufio.Reader) (Action, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return ActionNone, err
	}

	switch b {
	case 'q', 'Q':
		return ActionQuit, nil
	case ' ', 'x', 'X':
		return ActionToggle, nil
	case '\r', '\n':
		return ActionEnter, nil
	case 3: // Ctrl+C
		return ActionQuit, nil
	case 0x1b:
		// escape sequence, attempt to parse arrows
		if reader.Buffered() >= 2 {
			seq, err := reader.Peek(2)
			if err == nil && len(seq) == 2 && seq[0] == '[' {
				reader.Discard(2)
				switch seq[1] {
				case 'A':
					return ActionUp, nil
				case 'B':
					return ActionDown, nil
				case 'C':
					return ActionDown, nil
				case 'D':
					return ActionUp, nil
				}
			}
		}
		return ActionQuit, nil
	case 'k':
		return ActionUp, nil
	case 'j':
		return ActionDown, nil
	}

	return ActionNone, nil
}

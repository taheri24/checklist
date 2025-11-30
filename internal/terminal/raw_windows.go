//go:build windows

package terminal

import "syscall"

type State struct {
	mode uint32
}

// EnableRawMode switches the terminal to raw mode and returns the previous state.
func EnableRawMode(fd int) (*State, error) {
	handle := syscall.Handle(fd)

	var originalMode uint32
	if err := syscall.GetConsoleMode(handle, &originalMode); err != nil {
		return nil, err
	}

	rawMode := originalMode
	rawMode &^= syscall.ENABLE_ECHO_INPUT | syscall.ENABLE_LINE_INPUT | syscall.ENABLE_PROCESSED_INPUT
	rawMode |= syscall.ENABLE_VIRTUAL_TERMINAL_INPUT

	if err := syscall.SetConsoleMode(handle, rawMode); err != nil {
		return nil, err
	}

	return &State{mode: originalMode}, nil
}

// Restore resets the terminal to a previous state.
func Restore(fd int, state *State) error {
	if state == nil {
		return nil
	}

	handle := syscall.Handle(fd)
	if err := syscall.SetConsoleMode(handle, state.mode); err != nil {
		return err
	}

	return nil
}

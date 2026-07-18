package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// PromptSecret asks the user for a secret on the controlling terminal,
// without echoing it back. Nothing is ever written to disk. Falls back to
// a plain (echoed) stdin read when stdin isn't a terminal (e.g. piped in
// from a script), so scripted use is still possible but not the default.
func PromptSecret(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)

	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

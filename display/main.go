package display

import (
	"fmt"
	"strings"
	"os"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
	log "github.com/sirupsen/logrus"
)

const (
	WINDOW_WIDTH  = 40
	PADDING_WIDTH = 2
)

type Terminal struct {
	Width  int
	Height int
}

func InitTerminal() (*Terminal, error) {

	retv := &Terminal{}
	var err error

	if term.IsTerminal(0) {
		retv.Width, retv.Height, err = term.GetSize(0)
		if err != nil {
			return retv, err
		}
	} else {
		return retv, fmt.Errorf("not a terminal")
	}
	return retv, nil
}

// GetWindowSizeAndPadding gets the window size and padding
//
//	width: the width of the terminal
//	window: the size of the window
//
// returns:
//
//	width: the width of the window
//	padding: the padding of the window
//	error: any error
//
// e.g.
//
//	GetWindowSizeAndPadding(80, 10)
//	returns
//	width: 70
//
// padding: 5
//
//	+--------------------------------------+
//	|              width                   |
//	| padding | window           | padding |
//	+--------------------------------------+
func GetWindowSizeAndPadding(width, window int) (int, error) {
	if width < 2 {
		return 0, fmt.Errorf("width too small")
	}

	if width < window {
		return 0, fmt.Errorf("width too small for window")
	}

	paddingSize := width - window // e.g. 80 - 10 = 70
	paddingSize -= paddingSize % 2
	retv := paddingSize / 2
	return retv, nil
}

func WrapLines(msg string, window, padding int) []string {
	retv := []string{}
	foldwidth := uint(window - (PADDING_WIDTH * 2))
	paddingStr := strings.Repeat(" ", padding)

	for _, line := range strings.Split(wordwrap.WrapString(msg, foldwidth), "\n") {
		retv = append(retv, fmt.Sprintf("%s%s%s", paddingStr, line, paddingStr))
	}
	return retv
}

func MidLine(window, padding int) (string) {
	paddingStr := strings.Repeat(" ", padding)
	line := strings.Repeat("-", window)
	return fmt.Sprintf("%s%s%s", paddingStr, line, paddingStr)
}

func FatalIfFailed(err error, msg string) {
	errstr := fmt.Sprintf("Error: %s", err)

	// We have no error
	if err == nil {
		log.Debugf("No error")
		return
	}

	// Get terminal size
	terminal, err := InitTerminal()
	if err != nil {
		return
	}

	padding, err := GetWindowSizeAndPadding(terminal.Width, WINDOW_WIDTH)
	if err != nil {
		return
	}

	fmt.Printf("\n\n")
	for _, line := range WrapLines(errstr, WINDOW_WIDTH, padding) {
		fmt.Printf("%s\n", line)
	}
	fmt.Println(MidLine( WINDOW_WIDTH, padding))
	fmt.Printf("\n\n")
	for _, line := range WrapLines(msg, WINDOW_WIDTH, padding) {
		fmt.Printf("%s\n", line)
	}
	fmt.Printf("\n\n")
	os.Exit(1)
}

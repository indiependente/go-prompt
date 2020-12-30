// +build !windows

package term

import (
	"sync"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var (
	saveTermios     unix.Termios
	saveTermiosFD   int
	saveTermiosOnce sync.Once
)

func getOriginalTermios(fd int) (unix.Termios, error) {
	var (
		origTermiosPtr *unix.Termios
		err            error
	)
	saveTermiosOnce.Do(func() {
		saveTermiosFD = fd
		origTermiosPtr, err = termios.Tcgetattr(uintptr(fd))
		saveTermios = *origTermiosPtr
	})
	return saveTermios, err
}

// Restore terminal's mode.
func Restore() error {
	o, err := getOriginalTermios(saveTermiosFD)
	if err != nil {
		return err
	}
	return termios.Tcsetattr(uintptr(saveTermiosFD), termios.TCSANOW, &o)
}

// +build !windows

package os

import (
	"os"
	"os/signal"
	"syscall"
)

// IgnoreUserEnteredSignals ignores user signals
func IgnoreUserEnteredSignals() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTSTP, syscall.SIGQUIT}
	for _, s := range signals {
		signal.Ignore(s)
	}
}

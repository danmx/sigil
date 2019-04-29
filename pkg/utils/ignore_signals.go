// +build !windows

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// IgnoreUserEnteredSignals ignores user signals
func IgnoreUserEnteredSignals() {
	var signals []os.Signal
	signals = []os.Signal{syscall.SIGINT, syscall.SIGSTOP, syscall.SIGQUIT}
	for _, s := range signals {
		signal.Ignore(s)
	}
}

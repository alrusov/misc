package misc

import (
	"os"
	"os/signal"
	"syscall"
)

//----------------------------------------------------------------------------------------------------------------------------//

func signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGCHLD)

	for {
		signal := <-c

		switch signal {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			Logger("IN", "Signal \"%s\" received", signal.String())
			StopApp(0)
		case syscall.SIGCHLD:
			fallthrough
		default:
			Logger("DE", "Signal \"%s\" received", signal.String())
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

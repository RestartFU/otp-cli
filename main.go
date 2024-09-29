package main

import (
	"golang.org/x/term"
	"os"
	"os/signal"
	"otp/cmd/otp"
	"otp/internal/command"
	"syscall"
)

func main() {
	fd := int(os.Stdin.Fd())
	state, _ := term.GetState(fd)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	capture[os.Signal](c, func() {
		println()
		_ = term.Restore(fd, state)
		os.Exit(0)
	})

	command.Register("list", "list all registered services", "l", command.List{})
	command.Register("get", "get the latest passcode", "g", command.Get{})
	command.Register("add", "add a service with a secret", "a", command.Add{})
	command.Register("remove", "remove a specific service", "r", command.Remove{})
	command.Register("purge", "purge all registered services", "p", command.Purge{})

	otp.Run()
}

func capture[T comparable](c chan T, fn func()) {
	go func() {
		<-c
		fn()
	}()
}

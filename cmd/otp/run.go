package otp

import (
	"fmt"
	"os"
	"otp/internal/command"
	"strings"
)

func Run() {
	args := os.Args[1:]
	if len(args) == 0 {
		command.Usage()
	}

	arg := args[0]
	if strings.HasPrefix(arg, "-") {
		args[0] = arg[1:]
		cmd(args, command.ByAlias)
	} else {
		cmd(args, command.ByName)
	}
}

func cmd(args []string, fn func(string) (command.Runnable, bool)) {
	runnable, ok := fn(args[0])
	if !ok {
		fmt.Printf("unknown command: %s\n", args[0])
		return
	}
	runnable = command.ParseArgs(args[1:], runnable)

	err := runnable.Run()
	if err != nil {
		fmt.Println(err)
	}
}

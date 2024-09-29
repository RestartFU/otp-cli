package term

import (
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
)

func StdinPassword(prefix string) (string, error) {
	if psw, ok := os.LookupEnv("OTP_PASSWORD"); ok {
		return psw, nil
	}

	fmt.Print(prefix)
	fd := int(os.Stdin.Fd())
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}

	var psw string
	for {
		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil || b[0] == 'q' || b[0] == '\x03' {
			defer_(fd, oldState)
			os.Exit(0)
		}
		if b[0] == '\n' || b[0] == '\r' {
			break
		}
		if b[0] == '\u007F' {
			if len(psw) == 0 {
				continue
			}
			psw = psw[:len(psw)-1]
			_, _ = os.Stdin.Write([]byte("\b \b"))
			continue
		}
		psw += string(b)
		fmt.Print("*")
	}
	defer_(fd, oldState)
	if len(psw) <= 0 {
		return "", errors.New("no password was provided")
	}
	return psw, nil
}

func defer_(fd int, oldState *term.State) {
	_ = term.Restore(fd, oldState)
	fmt.Println("\r")
}

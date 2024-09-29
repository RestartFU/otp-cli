package command

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"otp/internal/crypto"
	"otp/internal/otp"
	"otp/internal/term"
)

type Get struct {
	Service string
}

func (g Get) Run() error {
	psw, _ := term.StdinPassword("Password: ")

	buf, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}

	content, err := crypto.Decrypt(psw, string(buf))
	if err != nil {
		return ErrWrongPassword
	}
	ko, err := godotenv.Unmarshal(content)
	if err != nil {
		return ErrWrongPassword
	}

	secret, ok := ko[g.Service]
	if !ok {
		return ErrServiceNotFound
	}

	passcode, _ := otp.Generate(secret)
	fmt.Println(passcode)
	return err
}

package command

import (
	"github.com/joho/godotenv"
	"os"
	"otp/internal/crypto"
	"otp/internal/term"
)

type Purge struct{}

func (Purge) Run() error {
	psw, _ := term.StdinPassword("Password: ")

	buf, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}

	content, err := crypto.Decrypt(psw, string(buf))
	if err != nil {
		return ErrWrongPassword
	}
	_, err = godotenv.Unmarshal(content)
	if err != nil {
		return ErrWrongPassword
	}

	return os.Remove(dataPath)
}

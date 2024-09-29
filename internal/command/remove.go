package command

import (
	"github.com/joho/godotenv"
	"os"
	"otp/internal/crypto"
	"otp/internal/term"
)

type Remove struct {
	Service string
}

func (r Remove) Run() error {
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

	_, ok := ko[r.Service]
	if !ok {
		return ErrServiceNotFound
	}

	delete(ko, r.Service)
	newStr, err := godotenv.Marshal(ko)
	if err != nil {
		return err
	}

	newBuf, err := crypto.Encrypt(psw, newStr)
	if err != nil {
		return err
	}
	return os.WriteFile(dataPath, []byte(newBuf), os.ModePerm)
}

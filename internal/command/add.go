package command

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"otp/internal/crypto"
	"otp/internal/otp"
	"otp/internal/term"
)

type Add struct {
	Service string
}

func (a Add) Run() error {
	buf, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}

	prefix := "Password: "
	if len(buf) <= 0 {
		prefix = "Initialize " + prefix
	}

	ko := make(map[string]string)
	psw, _ := term.StdinPassword(prefix)

	if len(buf) > 0 {
		content, err := crypto.Decrypt(psw, string(buf))
		if err != nil {
			return ErrWrongPassword
		}
		ko, err = godotenv.Unmarshal(content)
		if err != nil {
			return ErrWrongPassword
		}

		_, ok := ko[a.Service]
		if ok {
			return ErrServiceExists
		}
	} else {
		confirmation, _ := term.StdinPassword("Confirm Password: ")
		if confirmation != psw {
			return ErrPasswordNotMatch
		}
	}

	secret, _ := term.StdinPassword("Secret: ")
	if len(secret) < 16 {
		return ErrSecretTooShort
	}
	ko[a.Service] = secret

	passcode, _ := otp.Generate(secret)
	fmt.Println(passcode)

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

package command

import (
	"fmt"
	"github.com/joho/godotenv"
	"maps"
	"os"
	"otp/internal/crypto"
	"otp/internal/term"
	"slices"
	"strings"
)

type List struct{}

func (List) Run() error {
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

	keys := slices.Collect(maps.Keys(ko))
	fmt.Print("(" + dataPath + ") OTP LIST:\n - ")
	fmt.Println(strings.Join(keys, "\n - "))
	return err
}

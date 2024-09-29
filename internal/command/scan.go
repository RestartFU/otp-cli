package command

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"otp/internal/crypto"
	"otp/internal/otp"
	"otp/internal/term"
)

type Scan struct {
	Service string
	Path    string
}

func (s Scan) Run() error {
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

		_, ok := ko[s.Service]
		if ok {
			return ErrServiceExists
		}
	} else {
		confirmation, _ := term.StdinPassword("Confirm Password: ")
		if confirmation != psw {
			return ErrPasswordNotMatch
		}
	}

	secret := scan(s.Path)
	if len(secret) < 16 {
		return ErrSecretTooShort
	}
	ko[s.Service] = secret

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

func scan(path string) string {
	// open and decode image file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		panic(err)
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		panic(err)
	}

	return result.GetText()
}

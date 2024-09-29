package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	dataPath = filepath.Join(dataDir(), "data.ko")

	ErrWrongPassword    = errors.New("wrong password")
	ErrPasswordNotMatch = errors.New("passwords don't match")
	ErrSecretTooShort   = errors.New("secret must be at least 16 characters")
	ErrServiceNotFound  = errors.New("service not found")
	ErrServiceExists    = errors.New("service already exists")
)

func init() {
	_ = os.MkdirAll(dataDir(), os.ModePerm)
	_, err := os.Stat(dataPath)
	if errors.Is(err, os.ErrNotExist) {
		f, _ := os.Create(dataPath)
		_ = f.Close()
	}
}

func dataDir() string {
	dir, ok := os.LookupEnv("OTP_DATA_DIR")
	if ok {
		return filepath.Clean(dir)
	}

	dir, err := os.UserCacheDir()
	if err != nil {
		dir, err = os.UserConfigDir()
		if err != nil {
			dir, err = os.Getwd()
			if err != nil {
				dir = os.TempDir()
			}
		}
	}

	return filepath.Clean(filepath.Join(dir, "otp_"))
}

type Runnable interface {
	Run() error
}

var (
	commands     = map[string]Runnable{}
	descriptions = map[string]string{}
	aliases      = map[string]string{}
)

func ParseArgs(args []string, runnable Runnable) Runnable {
	typ := reflect.ValueOf(runnable)
	newTyp := reflect.New(typ.Type()).Elem()
	numField := newTyp.NumField()

	if len(args) < numField {
		Usage()
	}

	for i := range numField {
		field := newTyp.Field(i)
		if field.Type() != reflect.TypeOf("") {
			panic("field type must be string")
		}
		field.SetString(args[i])
	}
	return newTyp.Interface().(Runnable)
}

func fetchLowerCasedFields(runnable Runnable) []string {
	var fields []string
	typ := reflect.TypeOf(runnable)

	for i := range typ.NumField() {
		field := typ.Field(i)
		if field.Type != reflect.TypeOf("") {
			panic("field type must be string")
		}
		fields = append(fields, strings.ToLower(field.Name))
	}
	return fields
}

func ByName(name string) (Runnable, bool) {
	r, ok := commands[name]
	return r, ok
}

func ByAlias(alias string) (Runnable, bool) {
	n, ok := aliases[alias]
	if !ok {
		return nil, false
	}
	return ByName(n)
}

func Register(name string, description string, alias string, runnable Runnable) {
	commands[name] = runnable
	descriptions[alias] = description
	aliases[alias] = name
}

func Usage() {
	fmt.Printf("Usage: %s <command> <args>\n", filepath.Base(os.Args[0]))

	for name, runnable := range commands {
		var alias string
		for a, c := range aliases {
			if c == name {
				alias = a
			}
		}

		fmt.Printf("  - (-%s) %s", alias, name)
		for _, field := range fetchLowerCasedFields(runnable) {
			fmt.Printf(" <%s>", field)
		}

		println()
	}
	os.Exit(0)
}

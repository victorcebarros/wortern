package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	defaultPrefix = "$HOME/.local/share/wortern"
	prefix        = pflag.StringP("path", "p", defaultPrefix, "Sets storage path for words")
	lang          = pflag.StringP("lang", "l", "en", "Sets word language")
)

const (
	wordDB string = "words"
)

func OpenWordDB(prefix string) (io.WriteCloser, error) {
	prefix = os.ExpandEnv(prefix)
	err := os.MkdirAll(prefix, 0755)

	if err != nil {
		return nil, err
	}

	prefix = path.Join(prefix, wordDB)
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND

	file, err := os.OpenFile(prefix, flags, 0644)

	return file, err
}

func StoreWord(w io.Writer, lang, word, sentence string) error {
	_, err := fmt.Fprintf(w, "%s::%s::%s::%s\n", lang, time.Now(), word, sentence)
	return err
}

func main() {
	pflag.Parse()

	f, err := OpenWordDB(*prefix)
	defer f.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error: Could not open Word DB: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	if len(pflag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s word sentence?\n", os.Args[0])
		os.Exit(1)
	}

	word := pflag.Args()[0]
	sentence := strings.Join(pflag.Args()[1:], " ")

	err = StoreWord(f, *lang, word, sentence)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error: Could not store word: %s\n", os.Args[0], err)
	}
}

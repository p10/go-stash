package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	limit := flag.Int("limit", -1, "limit to show")
	flag.IntVar(limit, "l", -1, "alias for limit")

	take := flag.Int("take", -1, "take a stash")
	flag.IntVar(take, "t", -1, "alias for take")

	flag.Parse()

	dir := stashesDir()

	if *limit != -1 {
		err := list(dir, limit)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if *take != -1 {
		content, err := takeStash(dir, *take)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("%s:\n", string(content))
		return
	}

	b, _ := io.ReadAll(os.Stdin)
	if len(b) == 0 {
		fmt.Fprintln(os.Stderr, "no content to stash")
		os.Exit(1)
	}
	create(dir, b)
	fmt.Println("Stash created")
	fmt.Println(string(b))
}

func create(dir string, content []byte) {
	today := time.Now().Format(time.DateTime)
	path := filepath.Join(dir, fmt.Sprintf("%s.txt", today))
	if err := os.WriteFile(path, content, 0644); err != nil {
		panic(err)
	}
}

func list(dir string, limit *int) error {
	files := stashesFiles()

	if *limit > len(files) || *limit < 1 {
		return fmt.Errorf("limit %d is out of range: from 1 to %d", *limit, len(files))
	}

	files = files[len(files)-*limit:]

	for i, fileName := range files {
		body := readFile(filepath.Join(dir, fileName))
		// TODO: trim body (last 10 lines)

		color.Set(color.FgYellow)
		fmt.Printf("%d) %s\n", len(files)-i, strings.ReplaceAll(fileName, ".txt", ""))
		color.Unset()

		fmt.Printf("%s\n", string(body))
		if i != len(files)-1 {
			fmt.Printf("\n")
		}
	}

	return nil
}

func takeStash(dir string, reversedNumber int) ([]byte, error) {
	files := stashesFiles()

	if reversedNumber > len(files) || reversedNumber < 1 {
		return nil, fmt.Errorf(
			"stash number %d is out of range: from 1 to %d",
			reversedNumber,
			len(files),
		)
	}

	fileName := files[len(files)-reversedNumber]
	body := readFile(filepath.Join(dir, fileName))

	return body, nil
}

func readFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func stashesFiles() []string {
	files, err := os.ReadDir(stashesDir())
	if err != nil {
		panic(err)
	}

	fileNames := make([]string, len(files))

	for i, file := range files {
		fileNames[i] = file.Name()
	}
	return fileNames
}

func stashesDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".stashes")
}

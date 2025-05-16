package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	blockedWords []*regexp.Regexp
	redaction    = "redacted"
)

func checkBlockedWordsFilePerm(path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to stat %s: %v", path, err)
	}

	mode := info.Mode().Perm()
	if mode != 0o700 {
		log.Fatalf("Permission on %s must be 700, found %o", path, mode)
	}
}

func loadBlockedWords(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "redaction:") {
			if redaction != "redacted" {
				log.Fatalf("Multiple redaction definitions found in %s", path)
			}
			redaction = strings.TrimSpace(strings.TrimPrefix(line, "redaction:"))
			continue
		}
		blockedWords = append(blockedWords, regexp.MustCompile(regexp.QuoteMeta(line)))
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading %s: %v", path, err)
	}
}

func censorLine(line string) string {
	for _, re := range blockedWords {
		line = re.ReplaceAllString(line, redaction)
	}
	return line
}

func process(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			fmt.Fprintln(writer, line)
			continue
		}
		censored := censorLine(line)
		fmt.Fprintln(writer, censored)
	}
	writer.Flush()
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

func main() {
	if len(os.Args) == 1 && isInputFromTerminal() {
		fmt.Fprintln(os.Stderr, "No input provided. Exiting.")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to determine current user: %v", err)
	}
	confPath := filepath.Join(usr.HomeDir, ".blocked_words.txt")
	checkBlockedWordsFilePerm(confPath)
	loadBlockedWords(confPath)

	switch len(os.Args) {
	case 1:
		process(os.Stdin, os.Stdout)
	case 2:
		inFile, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Cannot open input file: %v", err)
		}
		defer inFile.Close()
		process(inFile, os.Stdout)
	case 3:
		inFile, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Cannot open input file: %v", err)
		}
		defer inFile.Close()
		outFile, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatalf("Cannot create output file: %v", err)
		}
		defer outFile.Close()
		process(inFile, outFile)
	default:
		log.Fatalf("Usage: %s [input_file] [output_file]", os.Args[0])
	}
}

func isInputFromTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return true
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}


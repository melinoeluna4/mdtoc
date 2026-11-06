package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const usage = `mdtoc generates a table of contents from Markdown headings.

Usage:
  mdtoc [flags] [file]

Flags:
  -max-depth int   maximum heading depth (default 3)
  -bullet string   bullet marker (default "-")

Examples:
  mdtoc README.md
  cat README.md | mdtoc -max-depth 2
  mdtoc -bullet "*" docs/guide.md

Exit codes:
  0  the document was read
  2  usage or I/O error
`

type heading struct {
	level int
	text  string
	slug  string
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("mdtoc", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprint(stderr, usage)
	}

	maxDepth := flags.Int("max-depth", 3, "maximum heading depth")
	bullet := flags.String("bullet", "-", "bullet marker")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		return 2
	}

	if *maxDepth < 1 || *maxDepth > 6 {
		fmt.Fprintln(stderr, "mdtoc: max-depth must be between 1 and 6")
		return 2
	}

	if *bullet == "" {
		fmt.Fprintln(stderr, "mdtoc: bullet cannot be empty")
		return 2
	}

	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "mdtoc: at most one file is allowed")
		return 2
	}

	reader := stdin
	if flags.NArg() == 1 {
		file, err := os.Open(flags.Arg(0))
		if err != nil {
			fmt.Fprintf(stderr, "mdtoc: %v\n", err)
			return 2
		}
		defer file.Close()

		reader = file
	}

	headings, err := parseHeadings(reader)
	if err != nil {
		fmt.Fprintf(stderr, "mdtoc: %v\n", err)
		return 2
	}

	printTOC(stdout, headings, *maxDepth, *bullet)

	return 0
}

func parseHeadings(reader io.Reader) ([]heading, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	headings := make([]heading, 0)
	seen := make(map[string]int)
	inFence := false

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(text, "```") || strings.HasPrefix(text, "~~~") {
			inFence = !inFence
			continue
		}

		if inFence {
			continue
		}

		level, title, ok := parseHeading(text)
		if !ok {
			continue
		}

		slug := slugify(title)
		if count, exists := seen[slug]; exists {
			seen[slug] = count + 1
			slug = fmt.Sprintf("%s-%d", slug, count)
		} else {
			seen[slug] = 0
		}

		headings = append(headings, heading{level: level, text: title, slug: slug})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return headings, nil
}

func parseHeading(line string) (int, string, bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", false
	}

	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}

	if level < 1 || level > 6 {
		return 0, "", false
	}

	if level == len(line) || line[level] != ' ' {
		return 0, "", false
	}

	title := strings.TrimSpace(line[level:])
	if strings.HasSuffix(title, "#") {
		title = strings.TrimSpace(strings.TrimRight(title, "#"))
	}

	title = strings.ReplaceAll(title, "`", "")
	title = strings.ReplaceAll(title, "**", "")
	title = strings.ReplaceAll(title, "__", "")

	if title == "" {
		return 0, "", false
	}

	return level, title, true
}

func slugify(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	var builder strings.Builder
	lastDash := false

	for _, character := range text {
		switch {
		case unicode.IsLetter(character), unicode.IsDigit(character):
			builder.WriteRune(character)
			lastDash = false
		case character == ' ', character == '-', character == '_', character == '/', character == '.':
			if !lastDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}

	slug := strings.Trim(builder.String(), "-")
	if slug == "" {
		return "section"
	}

	return slug
}

func printTOC(writer io.Writer, headings []heading, maxDepth int, bullet string) {
	for _, item := range headings {
		if item.level > maxDepth {
			continue
		}

		indent := strings.Repeat("  ", item.level-1)
		fmt.Fprintf(writer, "%s%s [%s](#%s)\n", indent, bullet, item.text, item.slug)
	}
}
